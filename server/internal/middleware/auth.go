package middleware

import (
	"net/http"
	"server/internal/database"
	"server/internal/models"
	util "server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Context keys for storing values in Fiber context
const (
	ContextKeyUserID    = "userID"
	ContextKeyAccountID = "accountID"
	ContextKeyUser      = "user"
)

// IsAuthenticated is a middleware that checks if the user is authenticated.
// It parses the JWT from the request and validates it.
// If valid, it stores the user ID and account ID in the context locals.
func IsAuthenticated(c *fiber.Ctx) error {
	jwt := util.GetJWT(c)

	userID, err := util.ParseJwt(jwt)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	// Store user ID in context
	c.Locals(ContextKeyUserID, userID)

	// Fetch user to get account ID
	masterDB := database.Manager.GetMasterDB()
	var user models.User
	if err := masterDB.First(&user, "id = ?", userID).Error; err == nil {
		// Store account ID in context
		c.Locals(ContextKeyAccountID, user.AccountID)
		c.Locals(ContextKeyUser, &user)
	}

	return c.Next()
}

// GetUserID retrieves the user ID from the Fiber context.
// Returns an empty string if not found.
func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetAccountID retrieves the account ID from the Fiber context.
// Returns an empty string if not found or user has no account.
func GetAccountID(c *fiber.Ctx) string {
	if accountID, ok := c.Locals(ContextKeyAccountID).(string); ok {
		return accountID
	}
	return ""
}

// GetCurrentUser retrieves the current user from the Fiber context.
// Returns nil if not found.
func GetCurrentUser(c *fiber.Ctx) *models.User {
	if user, ok := c.Locals(ContextKeyUser).(*models.User); ok {
		return user
	}
	return nil
}

// RequireAccount is a middleware that ensures the user belongs to an account.
// This is useful for routes that require multitenancy context.
func RequireAccount(c *fiber.Ctx) error {
	accountID := GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "account required",
		})
	}

	// Verify account exists and is active
	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "account not found",
		})
	}

	if !account.IsActive() {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "account is not active",
		})
	}

	return c.Next()
}

// RequireActiveSubscription is a middleware that ensures the user's account
// has an active paid subscription.
func RequireActiveSubscription(c *fiber.Ctx) error {
	accountID := GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "account required",
		})
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "account not found",
		})
	}

	if !account.HasActiveSubscription() {
		return c.Status(http.StatusPaymentRequired).JSON(fiber.Map{
			"message": "active subscription required",
		})
	}

	return c.Next()
}
