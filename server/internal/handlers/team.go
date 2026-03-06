package handlers

import (
	"crypto/rand"
	"log"
	"math/big"
	"server/internal/database"
	"server/internal/middleware"
	"server/internal/models"
	"server/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	userSync = services.NewUserSyncService()
)

// TeamMember represents a team member response
type TeamMemberResponse struct {
	ID                     string     `json:"id"`
	FirstName              string     `json:"firstName"`
	LastName               string     `json:"lastName"`
	Email                  string     `json:"email"`
	Role                   string     `json:"role"`
	RoleID                 uint       `json:"roleId"`
	Enabled                bool       `json:"enabled"`
	PasswordChangeRequired bool       `json:"passwordChangeRequired"`
	LastLoginAt            *time.Time `json:"lastLoginAt"`
}

// GetTeamMembers returns all team members for the current account
func GetTeamMembers(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "account ID required",
		})
	}

	tenantDB, err := database.Manager.GetConnection(accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to connect to tenant database",
		})
	}

	var users []models.User
	if err := tenantDB.Preload("Role").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch team members",
		})
	}

	// Convert to response format
	members := make([]TeamMemberResponse, len(users))
	for i, user := range users {
		members[i] = TeamMemberResponse{
			ID:                     user.ID,
			FirstName:              user.FirstName,
			LastName:               user.LastName,
			Email:                  user.Email,
			Role:                   user.Role.Name,
			RoleID:                 user.RoleID,
			Enabled:                user.Enabled,
			PasswordChangeRequired: user.PasswordChangeRequired,
			LastLoginAt:            user.LastLoginAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(members)
}

// InviteTeamMember invites a new member to the team
func InviteTeamMember(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	userID := middleware.GetUserID(c)
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "account ID required",
		})
	}

	// Check account limits
	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "account not found",
		})
	}

	// Count current users
	count, err := userSync.CountUsersInTenant(accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to count users",
		})
	}

	if !account.CanAddUser(int(count)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "user limit reached for your plan",
		})
	}

	// Parse request
	var data struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		RoleID    uint   `json:"roleId"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	// Validate
	if data.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email is required",
		})
	}

	if data.RoleID == 0 {
		data.RoleID = 3 // Default to Viewer
	}

	// Check if email already exists
	var existingUser models.User
	if masterDB.Where("email = ?", data.Email).First(&existingUser).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "a user with this email already exists",
		})
	}

	// Get inviter's name
	var inviter models.User
	masterDB.First(&inviter, "id = ?", userID)
	inviterName := inviter.FirstName
	if inviterName == "" {
		inviterName = inviter.Email
	}

	// Generate temporary password
	tempPassword := generateTempPassword(12)

	// Create user in master database
	newUser := models.User{
		AccountID:              accountID,
		FirstName:              data.FirstName,
		LastName:               data.LastName,
		Email:                  data.Email,
		RoleID:                 data.RoleID,
		Enabled:                true,
		PasswordChangeRequired: true,
	}
	newUser.SetPassword(tempPassword)

	if err := masterDB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
		})
	}

	// Sync to tenant database
	if err := userSync.SyncUserToTenant(&newUser); err != nil {
		log.Printf("Failed to sync user to tenant: %v", err)
	}

	// Send invitation email
	loginURL := app.Domain + "/login"
	if emailService != nil {
		if err := emailService.Send(data.Email, "You've been invited to join "+account.Name, "team_invitation", map[string]interface{}{
			"InviterName":  inviterName,
			"AccountName":  account.Name,
			"LoginURL":     loginURL,
			"TempPassword": tempPassword,
		}); err != nil {
			log.Printf("Failed to queue invitation email for %s: %v", data.Email, err)
		}
	} else {
		log.Printf("Team invitation for %s - Login: %s, Temp Password: %s", data.Email, loginURL, tempPassword)
	}

	return c.Status(fiber.StatusCreated).JSON(TeamMemberResponse{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		RoleID:    newUser.RoleID,
		Enabled:   newUser.Enabled,
	})
}

// UpdateTeamMember updates a team member's role or status
func UpdateTeamMember(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	currentUserID := middleware.GetUserID(c)

	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "account ID required",
		})
	}

	memberID := c.Params("id")
	if memberID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "member ID required",
		})
	}

	// Prevent self-modification of role
	if memberID == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot modify your own role",
		})
	}

	var data struct {
		RoleID  uint `json:"roleId"`
		Enabled bool `json:"enabled"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	// Update in master
	if err := masterDB.Model(&models.User{}).
		Where("id = ? AND account_id = ?", memberID, accountID).
		Updates(map[string]interface{}{
			"role_id": data.RoleID,
			"enabled": data.Enabled,
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update member",
		})
	}

	// Update in tenant
	tenantDB, err := database.Manager.GetConnection(accountID)
	if err == nil {
		tenantDB.Model(&models.User{}).
			Where("id = ?", memberID).
			Updates(map[string]interface{}{
				"role_id": data.RoleID,
				"enabled": data.Enabled,
			})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member updated successfully",
	})
}

// DeleteTeamMember removes a member from the team
func DeleteTeamMember(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	currentUserID := middleware.GetUserID(c)

	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "account ID required",
		})
	}

	memberID := c.Params("id")
	if memberID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "member ID required",
		})
	}

	// Prevent self-deletion
	if memberID == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot delete yourself",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	// Soft delete in master
	if err := masterDB.Where("id = ? AND account_id = ?", memberID, accountID).
		Delete(&models.User{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete member",
		})
	}

	// Delete from tenant
	if err := userSync.DeleteUserFromTenant(memberID, accountID); err != nil {
		log.Printf("Warning: failed to delete user from tenant: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member removed successfully",
	})
}

// ResendInvitation resends the invitation email to a team member
func ResendInvitation(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "account ID required",
		})
	}

	memberID := c.Params("id")
	if memberID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "member ID required",
		})
	}

	masterDB := database.Manager.GetMasterDB()
	var user models.User
	if err := masterDB.First(&user, "id = ? AND account_id = ?", memberID, accountID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "member not found",
		})
	}

	// Generate new temporary password
	tempPassword := generateTempPassword(12)
	user.SetPassword(tempPassword)
	user.PasswordChangeRequired = true
	masterDB.Save(&user)

	loginURL := app.Domain + "/login"
	if emailService != nil {
		if err := emailService.Send(user.Email, "Your invitation has been resent", "team_invitation", map[string]interface{}{
			"InviterName":  "your admin",
			"AccountName":  accountID,
			"LoginURL":     loginURL,
			"TempPassword": tempPassword,
		}); err != nil {
			log.Printf("Failed to queue resend invitation email for %s: %v", user.Email, err)
		}
	} else {
		log.Printf("Resend invitation for %s - Login: %s, Temp Password: %s", user.Email, loginURL, tempPassword)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "invitation resent successfully",
	})
}

// generateTempPassword generates a cryptographically random temporary password
func generateTempPassword(length int) string {
	b := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// SetPassword hashes a password using bcrypt
func SetPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return hash
}
