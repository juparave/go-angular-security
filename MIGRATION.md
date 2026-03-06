# Migration Guide

This guide helps you upgrade from previous versions of go-angular-security to the new multitenant architecture.

## Overview of Changes

### Database Architecture

**Before:** Single database with all users and data mixed together.

**After:**
- Master database: Stores accounts and user credentials
- Tenant databases: Each account has its own isolated database at `data/{accountID}/data.db`

### User Model Changes

The `User` model has been updated:

```go
// New fields
AccountID               string     // Foreign key to tenant
RoleID                  uint       // Single role (replaces Roles array)
PasswordChangeRequired  bool       // Security flag
PasswordDigest          string     // For reset token validation
Enabled                 bool       // Account status
LastLoginAt             *time.Time // Audit trail

// Removed
StripeCustomerID        // Moved to Account model
```

## Migration Steps

### Step 1: Backup Your Data

```bash
# Backup your existing database
cp gorm.db gorm.db.backup
```

### Step 2: Update Dependencies

```bash
cd server
go mod tidy
```

### Step 3: Create Migration Script

Create a file `scripts/migrate_to_multitenant.go`:

```go
package main

import (
    "log"
    "server/internal/database"
    "server/internal/models"
)

func main() {
    // Initialize database manager
    if err := database.InitManager("gorm.db", "data"); err != nil {
        log.Fatal(err)
    }
    defer database.Manager.CloseAll()

    masterDB := database.Manager.GetMasterDB()

    // Get all existing users
    var users []models.User
    if err := masterDB.Find(&users).Error; err != nil {
        log.Fatal(err)
    }

    // Group users by potential account (customize this logic)
    // For now, each user gets their own account
    for _, user := range users {
        // Create account for user
        account := models.Account{
            Name:         user.FirstName + "'s Account",
            ContactEmail: user.Email,
            Status:       models.AccountStatusActive,
            PlanTier:     models.TierFree,
            MaxUsers:     1,
        }

        if err := masterDB.Create(&account).Error; err != nil {
            log.Printf("Failed to create account for %s: %v", user.Email, err)
            continue
        }

        // Update user with account ID
        user.AccountID = account.ID
        user.RoleID = 1 // Admin
        user.Enabled = true

        if err := masterDB.Save(&user).Error; err != nil {
            log.Printf("Failed to update user %s: %v", user.Email, err)
            continue
        }

        // Sync user to tenant database
        tenantDB, err := database.Manager.GetConnection(account.ID)
        if err != nil {
            log.Printf("Failed to get tenant DB: %v", err)
            continue
        }

        if err := tenantDB.Create(&user).Error; err != nil {
            log.Printf("Failed to sync user to tenant: %v", err)
        }
    }

    log.Println("Migration completed!")
}
```

### Step 4: Run Migration

```bash
go run scripts/migrate_to_multitenant.go
```

### Step 5: Update Frontend

1. Update `package.json` dependencies
2. Update auth service to use new endpoints

```typescript
// Old endpoint
this.http.post('/api/v1/request_reset_password', data)

// New endpoint
this.http.post('/api/v1/request-password-reset', data)
```

### Step 6: Update Environment Variables

Add new required variables to your `.env`:

```env
# New variables
JWT_SECRET=your-jwt-secret
JWT_REFRESH_SECRET=your-refresh-secret
JWT_RESET_SECRET=your-reset-secret
TENANT_DATA_PATH=./data
```

## Breaking Changes

### API Endpoints

| Old | New |
|-----|-----|
| `POST /register` | `POST /register-account` |
| `POST /request_reset_password` | `POST /request-password-reset` |
| `POST /refresh_token` | `POST /refresh-token` |

### User Response

The user object now includes `account` and `role` instead of `roles`:

```json
// Before
{
  "id": "abc123",
  "roles": [{"name": "admin"}]
}

// After
{
  "id": "abc123",
  "accountId": "acc_xyz",
  "role": {"name": "Admin"},
  "account": {"name": "My Account"}
}
```

### Login Response

The login response now includes `passwordChangeRequired`:

```json
{
  "message": "success",
  "user": {...},
  "passwordChangeRequired": false
}
```

## Rollback

If you need to rollback:

1. Stop the server
2. Restore the backup: `cp gorm.db.backup gorm.db`
3. Remove tenant databases: `rm -rf data/`
4. Revert to previous version
5. Restart the server

## Need Help?

If you encounter issues during migration:

1. Check the server logs for errors
2. Verify database connections
3. Ensure all environment variables are set
4. Open an issue on GitHub with details
