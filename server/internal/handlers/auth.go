package handlers

import (
	"log"
	"server/internal/database"
	"server/internal/models"
	"server/internal/services"
	"server/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	keyPassword        = "password"
	keyConfirmPassword = "confirmPassword"
	keyFirstName       = "firstName"
	keyLastName        = "lastName"
	keyEmail           = "email"
	keyRefreshToken    = "refreshToken"
	keyAccountName     = "accountName"
)

var userSyncService = services.NewUserSyncService()

// Register handles new user registration (legacy, for backward compatibility).
// For self-service account registration, use RegisterAccount instead.
func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	if data[keyPassword] != data[keyConfirmPassword] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName: data[keyFirstName],
		LastName:  data[keyLastName],
		Email:     data[keyEmail],
		RoleID:    1, // Admin role for standalone users
		Enabled:   true,
	}

	user.SetPassword(data[keyPassword])

	res := database.DB.Create(&user)
	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	err := utils.GenerateUserTokens(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate tokens",
		})
	}

	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(user)
}

// RegisterAccount handles self-service account registration.
// Creates a new account (tenant) with the first user as admin.
func RegisterAccount(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	// Validate required fields
	if data[keyEmail] == "" || data[keyPassword] == "" || data[keyAccountName] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email, password, and account name are required",
		})
	}

	if data[keyPassword] != data[keyConfirmPassword] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	// Check if email already exists
	var existingUser models.User
	if masterDB.Where("email = ?", data[keyEmail]).First(&existingUser).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "an account with this email already exists",
		})
	}

	// Create account with free tier active subscription
	account := models.Account{
		Name:               data[keyAccountName],
		ContactEmail:       data[keyEmail],
		Status:             models.AccountStatusActive,
		PlanTier:           models.TierFree,
		SubscriptionStatus: models.SubscriptionStatusActive,
		MaxUsers:           1, // Free tier: 1 user
	}

	if err := masterDB.Create(&account).Error; err != nil {
		log.Printf("Failed to create account: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create account",
		})
	}

	// Create user in master database
	user := models.User{
		AccountID:              account.ID,
		FirstName:              data[keyFirstName],
		LastName:               data[keyLastName],
		Email:                  data[keyEmail],
		RoleID:                 1, // Admin
		Enabled:                true,
		PasswordChangeRequired: false,
	}
	user.SetPassword(data[keyPassword])

	if err := masterDB.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		// Clean up account
		masterDB.Delete(&account)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
		})
	}

	// Sync user to tenant database
	if err := userSyncService.SyncUserToTenant(&user); err != nil {
		log.Printf("Failed to sync user to tenant: %v", err)
	}

	// Generate tokens
	if err := utils.GenerateUserTokens(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate tokens",
		})
	}

	// Load role for response
	masterDB.Preload("Account").First(&user, "id = ?", user.ID)

	setCookies(c, user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "account created successfully",
		"user":    user,
	})
}

// Login handles user authentication.
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	var user models.User
	if err := masterDB.Preload("Account").Where("email = ?", data[keyEmail]).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{
				"email": []string{"not found"},
			},
		})
	}

	// Check if user is enabled
	if !user.Enabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "account is disabled",
		})
	}

	// Verify password
	if err := user.ComparePassword(data[keyPassword]); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fiber.Map{
				"password": []string{"incorrect"},
			},
		})
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	masterDB.Model(&user).Update("last_login_at", now)

	// Generate tokens
	if err := utils.GenerateUserTokens(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate tokens",
		})
	}

	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":                "success",
		"user":                   user,
		"passwordChangeRequired": user.PasswordChangeRequired,
	})
}

// RefreshToken handles JWT token refresh.
func RefreshToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	refreshToken := data[keyRefreshToken]
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "refresh token is required",
		})
	}

	issuer, err := utils.ParseJwt(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "token invalid or expired",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	var user models.User
	if err := masterDB.Preload("Account").First(&user, "id = ?", issuer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{
				"user": []string{"not found"},
			},
		})
	}

	// Check if user is still enabled
	if !user.Enabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "account is disabled",
		})
	}

	if err := utils.GenerateUserTokens(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate tokens",
		})
	}

	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

// User retrieves the current authenticated user.
func User(c *fiber.Ctx) error {
	jwt := utils.GetJWT(c)
	userID, err := utils.ParseJwt(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	var user models.User
	if err := masterDB.Preload("Account").Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{
				"user": []string{"not found"},
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Logout handles user logout.
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     utils.CookieJWT,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:     utils.CookieRefreshJWT,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// UpdateInfo handles updating user profile information.
func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	user := models.User{
		ID:        userID,
		FirstName: data[keyFirstName],
		LastName:  data[keyLastName],
		Email:     data[keyEmail],
	}

	if err := masterDB.Model(&user).Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update user",
		})
	}

	// Sync to tenant database
	masterDB.First(&user, "id = ?", userID)
	if user.AccountID != "" {
		userSyncService.UpdateUserInTenant(&user)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdatePassword handles password change (legacy, use ChangePassword for authenticated changes).
func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if data[keyPassword] != data[keyConfirmPassword] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	var user models.User
	if err := masterDB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	user.SetPassword(data[keyPassword])
	user.PasswordChangeRequired = false

	if err := masterDB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update password",
		})
	}

	// Sync to tenant database
	if user.AccountID != "" {
		userSyncService.UpdateUserPasswordInTenant(&user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password updated successfully",
	})
}

// setCookies sets JWT cookies for authentication.
func setCookies(c *fiber.Ctx, user models.User) {
	secure := app != nil && app.Mode == "production"
	cookie := fiber.Cookie{
		Name:     utils.CookieJWT,
		Value:    user.AccessToken,
		Expires:  time.Now().Add(utils.AccessTokenDuration()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "lax",
	}
	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:     utils.CookieRefreshJWT,
		Value:    user.RefreshToken,
		Expires:  time.Now().Add(utils.RefreshTokenDuration()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "lax",
	}
	c.Cookie(&refreshCookie)
}
