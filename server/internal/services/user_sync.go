package services

import (
	"fmt"
	"server/internal/database"
	"server/internal/models"
)

// UserSyncService handles syncing user data between master and tenant databases
type UserSyncService struct{}

// NewUserSyncService creates a new UserSyncService
func NewUserSyncService() *UserSyncService {
	return &UserSyncService{}
}

// SyncUserToTenant syncs a user from the master database to their tenant database.
// This ensures the user exists in both places with consistent data.
func (s *UserSyncService) SyncUserToTenant(user *models.User) error {
	if user.AccountID == "" {
		return fmt.Errorf("user has no account ID")
	}

	tenantDB, err := database.Manager.GetConnection(user.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get tenant connection: %w", err)
	}

	// Create or update user in tenant database
	// We use a map to update only specific fields to avoid overwriting
	result := tenantDB.Where("id = ?", user.ID).
		Assign(models.User{
			FirstName:               user.FirstName,
			LastName:                user.LastName,
			Email:                   user.Email,
			Password:                user.Password,
			RoleID:                  user.RoleID,
			PasswordChangeRequired:  user.PasswordChangeRequired,
			PasswordDigest:          user.PasswordDigest,
			Enabled:                 user.Enabled,
			LastLoginAt:             user.LastLoginAt,
		}).
		FirstOrCreate(&models.User{ID: user.ID})

	return result.Error
}

// DeleteUserFromTenant removes a user from their tenant database.
// The user record in the master database is preserved for audit purposes.
func (s *UserSyncService) DeleteUserFromTenant(userID, accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID is required")
	}

	tenantDB, err := database.Manager.GetConnection(accountID)
	if err != nil {
		return fmt.Errorf("failed to get tenant connection: %w", err)
	}

	result := tenantDB.Where("id = ?", userID).Delete(&models.User{})
	return result.Error
}

// GetUserFromTenant retrieves a user from their tenant database
func (s *UserSyncService) GetUserFromTenant(userID, accountID string) (*models.User, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	tenantDB, err := database.Manager.GetConnection(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant connection: %w", err)
	}

	var user models.User
	if err := tenantDB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserInTenant updates a user's information in their tenant database
func (s *UserSyncService) UpdateUserInTenant(user *models.User) error {
	if user.AccountID == "" {
		return fmt.Errorf("user has no account ID")
	}

	tenantDB, err := database.Manager.GetConnection(user.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get tenant connection: %w", err)
	}

	return tenantDB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"first_name":              user.FirstName,
			"last_name":               user.LastName,
			"email":                   user.Email,
			"role_id":                 user.RoleID,
			"enabled":                 user.Enabled,
			"password_change_required": user.PasswordChangeRequired,
		}).Error
}

// UpdateUserPasswordInTenant updates a user's password in both master and tenant
func (s *UserSyncService) UpdateUserPasswordInTenant(user *models.User) error {
	if user.AccountID == "" {
		return fmt.Errorf("user has no account ID")
	}

	// Update in master
	masterDB := database.Manager.GetMasterDB()
	if err := masterDB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password":                 user.Password,
			"password_digest":          user.PasswordDigest,
			"password_change_required": user.PasswordChangeRequired,
		}).Error; err != nil {
		return fmt.Errorf("failed to update password in master: %w", err)
	}

	// Update in tenant
	tenantDB, err := database.Manager.GetConnection(user.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get tenant connection: %w", err)
	}

	return tenantDB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password":                 user.Password,
			"password_digest":          user.PasswordDigest,
			"password_change_required": user.PasswordChangeRequired,
		}).Error
}

// CountUsersInTenant returns the number of users in a tenant
func (s *UserSyncService) CountUsersInTenant(accountID string) (int64, error) {
	if accountID == "" {
		return 0, fmt.Errorf("account ID is required")
	}

	tenantDB, err := database.Manager.GetConnection(accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get tenant connection: %w", err)
	}

	var count int64
	if err := tenantDB.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
