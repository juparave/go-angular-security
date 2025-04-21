package middleware

import (
	"errors"
	"server/internal/database"
	"server/internal/models"
	"server/internal/util"

	"github.com/gofiber/fiber/v2"
)

// IsAuthorized is a middleware function that checks if the authenticated user has the required permission
// to access a specific resource or perform an action.
// It retrieves the user ID from the JWT cookie, fetches the user's roles and associated permissions from the database.
//
// Parameters:
//   - c: The Fiber context, providing access to the request and response.
//   - reqPerm: A string representing the base name of the required permission (e.g., "users", "products").
//     The function checks for permissions like "view_users" or "edit_users" based on the request method.
//
// Behavior:
//  1. Parses the JWT cookie to get the user ID. Returns an error if parsing fails.
//  2. If `reqPerm` is empty, it allows the request to proceed (no specific permission required).
//  3. Fetches the user's roles from the database based on the user ID.
//  4. Fetches all permissions associated with those roles.
//  5. Checks permissions based on the HTTP method:
//     - For GET requests: Allows access if the user has either "view_<reqPerm>" or "edit_<reqPerm>".
//     - For other methods (POST, PUT, DELETE, etc.): Allows access only if the user has "edit_<reqPerm>".
//  6. If the required permission is found, it returns `nil`, allowing the request to proceed to the next handler.
//  7. If the required permission is not found, it sets the HTTP status to 401 Unauthorized and returns an "unauthorized" error.
func IsAuthorized(c *fiber.Ctx, reqPerm string) error {
	// Retrieve JWT from cookie
	cookie := c.Cookies("jwt")

	// Parse JWT to get user ID
	userID, err := util.ParseJwt(cookie)
	if err != nil {
		// If JWT parsing fails, consider the user unauthorized
		c.Status(fiber.StatusUnauthorized)
		return errors.New("unauthorized: invalid token")
	}

	// If no specific permission is required for this route, allow access
	if reqPerm == "" {
		return nil // Proceed to the next handler
	}

	// Prepare user model to query associations
	user := models.User{
		ID: userID,
	}

	// Find all roles associated with the user
	var roles []models.Role
	database.DB.Model(&user).Association("Roles").Find(&roles)

	// Find all permissions associated with the user's roles
	var permissions []models.Permission
	// Note: This might fetch duplicate permissions if a user has multiple roles with the same permission.
	// Consider optimizing this if performance becomes an issue or using a Set-like structure.
	database.DB.Model(&roles).Association("Permissions").Find(&permissions)

	// Check permissions based on the request method
	requiredViewPerm := "view_" + reqPerm
	requiredEditPerm := "edit_" + reqPerm

	if c.Method() == fiber.MethodGet {
		// For GET requests, view or edit permissions are sufficient
		for _, permission := range permissions {
			if permission.Name == requiredViewPerm || permission.Name == requiredEditPerm {
				return nil // User has required permission, proceed
			}
		}
	} else {
		// For non-GET requests (POST, PUT, DELETE, etc.), edit permission is required
		for _, permission := range permissions {
			if permission.Name == requiredEditPerm {
				return nil // User has required permission, proceed
			}
		}
	}

	// If no matching permission was found, deny access
	c.Status(fiber.StatusForbidden)
	return errors.New("forbidden: insufficient permissions")
}
