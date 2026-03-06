package middleware

import (
	"errors"
	"net/http"
	"server/internal/database"
	"server/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Role constants for simplified RBAC
const (
	RoleAdmin  = "Admin"
	RoleEditor = "Editor"
	RoleViewer = "Viewer"
)

// Role levels for hierarchical permission checking
var roleLevels = map[string]int{
	RoleAdmin:  3,
	RoleEditor: 2,
	RoleViewer: 1,
}

// IsAuthorized is a middleware that checks if the user has the required permission.
// It uses a simplified single-role RBAC system with hierarchical permissions:
// - Admin: Full access
// - Editor: Can create, edit, and view
// - Viewer: Can only view
//
// Parameters:
//   - c: The Fiber context
//   - reqPerm: The resource name (e.g., "users", "products")
//
// For GET requests: Viewer, Editor, and Admin can access
// For other methods (POST, PUT, DELETE): Editor and Admin can access
func IsAuthorized(c *fiber.Ctx, reqPerm string) error {
	userID := GetUserID(c)
	if userID == "" {
		c.Status(http.StatusUnauthorized)
		return errors.New("unauthorized: not authenticated")
	}

	// If no specific permission required, allow access
	if reqPerm == "" {
		return nil
	}

	// Get user's role from the appropriate database
	accountID := GetAccountID(c)
	var user models.User
	var role models.Role

	if accountID != "" {
		// Get role from tenant database
		tenantDB, err := database.Manager.GetConnection(accountID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return errors.New("internal server error")
		}

		if err := tenantDB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
			c.Status(http.StatusForbidden)
			return errors.New("forbidden: user not found in tenant")
		}
		role = user.Role
	} else {
		// Legacy: Get role from master database
		masterDB := database.Manager.GetMasterDB()
		if err := masterDB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
			c.Status(http.StatusForbidden)
			return errors.New("forbidden: user not found")
		}
		role = user.Role
	}

	// Check permissions based on role level
	minRequiredLevel := getMinRequiredLevel(c.Method())

	userLevel, exists := roleLevels[role.Name]
	if !exists {
		c.Status(http.StatusForbidden)
		return errors.New("forbidden: unknown role")
	}

	if userLevel >= minRequiredLevel {
		return nil
	}

	c.Status(http.StatusForbidden)
	return errors.New("forbidden: insufficient permissions")
}

// RequireRole is a middleware that checks if the user has at least the specified role.
func RequireRole(minRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		if userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "not authenticated",
			})
		}

		accountID := GetAccountID(c)
		var user models.User

		if accountID != "" {
			tenantDB, err := database.Manager.GetConnection(accountID)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
			if err := tenantDB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{
					"message": "forbidden",
				})
			}
		} else {
			masterDB := database.Manager.GetMasterDB()
			if err := masterDB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{
					"message": "forbidden",
				})
			}
		}

		userLevel, exists := roleLevels[user.Role.Name]
		if !exists {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

		minLevel, exists := roleLevels[minRole]
		if !exists {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "invalid role configuration",
			})
		}

		if userLevel < minLevel {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() fiber.Handler {
	return RequireRole(RoleAdmin)
}

// RequireEditor is a convenience middleware for editor+ routes
func RequireEditor() fiber.Handler {
	return RequireRole(RoleEditor)
}

// getMinRequiredLevel returns the minimum role level required for the HTTP method
func getMinRequiredLevel(method string) int {
	switch method {
	case http.MethodGet:
		return roleLevels[RoleViewer] // Viewers can read
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return roleLevels[RoleEditor] // Editors can modify
	default:
		return roleLevels[RoleAdmin] // Admin for anything else
	}
}

// IsAdmin checks if the current user has admin role
func IsAdmin(c *fiber.Ctx) bool {
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	return user.Role.Name == RoleAdmin
}

// IsEditorOrAdmin checks if the current user has editor or admin role
func IsEditorOrAdmin(c *fiber.Ctx) bool {
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	level, exists := roleLevels[user.Role.Name]
	if !exists {
		return false
	}
	return level >= roleLevels[RoleEditor]
}
