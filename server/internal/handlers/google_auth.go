package handlers

import (
	"log"
	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// GoogleLogin handles Google OAuth sign-in from the frontend.
// The frontend completes the Google sign-in flow and sends the verified user info.
func GoogleLogin(c *fiber.Ctx) error {
	var data struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if data.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email is required",
		})
	}

	masterDB := database.Manager.GetMasterDB()

	var user models.User
	result := masterDB.Preload("Account").Where("email = ?", data.Email).First(&user)

	if result.Error != nil {
		// Create a new user without a password (Google auth only)
		user = models.User{
			Email:                  data.Email,
			FirstName:              data.FirstName,
			LastName:               data.LastName,
			RoleID:                 1, // Admin for self-registered Google users
			Enabled:                true,
			PasswordChangeRequired: false,
		}

		if err := masterDB.Create(&user).Error; err != nil {
			log.Printf("Failed to create Google user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to create user",
			})
		}

		// Sync to tenant if account is set
		if user.AccountID != "" {
			if err := userSyncService.SyncUserToTenant(&user); err != nil {
				log.Printf("Failed to sync Google user to tenant: %v", err)
			}
		}

		masterDB.Preload("Account").First(&user, "id = ?", user.ID)
	}

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
