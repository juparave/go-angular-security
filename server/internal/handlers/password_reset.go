package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	// ResetTokenDuration is the time a password reset token is valid
	ResetTokenDuration = time.Hour * 1 // 60 minutes
	// ResetTokenSecretEnv is the environment variable for the reset token secret
	ResetTokenSecretEnv = "RESET_TOKEN_SECRET"
)

// RequestResetPassword is an alias for RequestPasswordReset for backwards compatibility with tests.
var RequestResetPassword = RequestPasswordReset

// PasswordResetClaims represents the JWT claims for password reset tokens
type PasswordResetClaims struct {
	Email          string `json:"email"`
	PasswordDigest string `json:"digest"` // Hash of current password to detect changes
	jwt.RegisteredClaims
}

// RequestPasswordReset initiates the password reset process.
// It sends a reset link to the user's email if the account exists.
// For security, it always returns success even if the email doesn't exist.
func RequestPasswordReset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	email := data[keyEmail]
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email is required",
		})
	}

	// Find user by email
	var user models.User
	result := database.Manager.GetMasterDB().Where("email = ?", email).First(&user)

	// Always return success to prevent email enumeration
	if result.Error != nil {
		log.Printf("Password reset requested for non-existent email: %s", email)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If an account with that email exists, a reset link has been sent.",
		})
	}

	// Check if user is enabled
	if !user.Enabled {
		log.Printf("Password reset requested for disabled account: %s", email)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If an account with that email exists, a reset link has been sent.",
		})
	}

	// Generate password digest (hash of current password)
	passwordDigest := generatePasswordDigest(user.Password)

	// Generate reset token
	token, err := generateResetToken(user.Email, passwordDigest)
	if err != nil {
		log.Printf("Failed to generate reset token: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If an account with that email exists, a reset link has been sent.",
		})
	}

	// Store the password digest in the user record for validation
	user.PasswordDigest = passwordDigest
	database.Manager.GetMasterDB().Model(&user).Update("password_digest", passwordDigest)

	// Build reset URL from handler-level app config
	resetURL := app.Domain + "/reset-password?token=" + token

	if app.Mode == "development" {
		log.Printf("Password reset URL for %s: %s", email, resetURL)
	}

	// Send email via email service, fall back to stub
	if emailService != nil {
		if err := emailService.Send(user.Email, "Reset your password", "password_reset", map[string]interface{}{
			"Name":     user.FirstName,
			"ResetURL": resetURL,
		}); err != nil {
			log.Printf("Failed to queue password reset email: %v", err)
		}
	} else {
		log.Printf("Email service not configured; reset URL for %s: %s", email, resetURL)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If an account with that email exists, a reset link has been sent.",
	})
}

// ResetPassword handles the actual password reset using a valid token.
// It validates the token, checks for token reuse via password digest, and updates the password.
func ResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	token := data["token"]
	password := data[keyPassword]
	confirmPassword := data[keyConfirmPassword]

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "token is required",
		})
	}

	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "password is required",
		})
	}

	if password != confirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// Validate password strength
	if err := validatePasswordStrength(password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Parse and validate token
	claims, err := parseResetToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	// Find user by email
	var user models.User
	result := database.Manager.GetMasterDB().Where("email = ?", claims.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	// Check if password has been changed since token was issued (token reuse prevention)
	if user.PasswordDigest != "" && user.PasswordDigest != claims.PasswordDigest {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "this reset link is no longer valid",
		})
	}

	// Update password
	user.SetPassword(password)
	user.PasswordChangeRequired = false
	user.PasswordDigest = generatePasswordDigest(user.Password) // Update digest

	if err := database.Manager.GetMasterDB().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update password",
		})
	}

	// Also update in tenant database if user has an account
	if user.AccountID != "" {
		tenantDB, err := database.Manager.GetConnection(user.AccountID)
		if err == nil {
			tenantDB.Model(&models.User{}).
				Where("id = ?", user.ID).
				Updates(map[string]interface{}{
					"password":                 user.Password,
					"password_digest":          user.PasswordDigest,
					"password_change_required": false,
				})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password has been reset successfully",
	})
}

// ChangePassword handles password change for authenticated users.
// It requires the current password to be provided for verification.
func ChangePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	currentPassword := data["currentPassword"]
	newPassword := data["newPassword"]
	confirmPassword := data["confirmPassword"]

	if currentPassword == "" || newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "current password and new password are required",
		})
	}

	if newPassword != confirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// Validate password strength
	if err := validatePasswordStrength(newPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get current user
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	var user models.User
	if err := database.Manager.GetMasterDB().First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// Verify current password
	if err := user.ComparePassword(currentPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fiber.Map{
				"currentPassword": []string{"incorrect"},
			},
		})
	}

	// Check new password is different from current
	if err := user.ComparePassword(newPassword); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "new password must be different from current password",
		})
	}

	// Update password
	user.SetPassword(newPassword)
	user.PasswordChangeRequired = false
	user.PasswordDigest = generatePasswordDigest(user.Password)

	if err := database.Manager.GetMasterDB().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update password",
		})
	}

	// Also update in tenant database
	if user.AccountID != "" {
		tenantDB, err := database.Manager.GetConnection(user.AccountID)
		if err == nil {
			tenantDB.Model(&models.User{}).
				Where("id = ?", user.ID).
				Updates(map[string]interface{}{
					"password":                 user.Password,
					"password_digest":          user.PasswordDigest,
					"password_change_required": false,
				})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password changed successfully",
	})
}

// generateResetToken creates a JWT token for password reset
func generateResetToken(email, passwordDigest string) (string, error) {
	claims := PasswordResetClaims{
		Email:          email,
		PasswordDigest: passwordDigest,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ResetTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetResetSecret()))
}

// parseResetToken validates and parses a password reset token
func parseResetToken(tokenString string) (*PasswordResetClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PasswordResetClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(utils.GetResetSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*PasswordResetClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generatePasswordDigest creates a hash of the password for token reuse detection
func generatePasswordDigest(password []byte) string {
	hash := sha256.Sum256(password)
	return hex.EncodeToString(hash[:])
}

// validatePasswordStrength checks if the password meets minimum requirements
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
