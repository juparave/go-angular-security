package middleware

import (
	"errors"
	"server/internal/database"
	"server/internal/models"
	"server/internal/util"

	"github.com/gofiber/fiber/v2"
)

// IsAuthorized verifies user have the required permission
func IsAuthorized(c *fiber.Ctx, reqPerm string) error {
	cookie := c.Cookies("jwt")

	userID, err := util.ParseJwt(cookie)
	if err != nil {
		return err
	}

	// if no permission required, allow
	if reqPerm == "" {
		return nil
	}

	user := models.User{
		ID: userID,
	}

	var roles []models.Role
	// find all user's roles
	database.DB.Model(&user).Association("Roles").Find(&roles)

	var permissions []models.Permission
	// find all permissions for user's roles
	database.DB.Model(&roles).Association("Permissions").Find(&permissions)

	if c.Method() == "GET" {
		for _, permission := range permissions {
			if permission.Name == "view_"+reqPerm || permission.Name == "edit_"+reqPerm {
				return nil
			}
		}
	} else {
		for _, permission := range permissions {
			if permission.Name == "edit_"+reqPerm {
				return nil
			}
		}
	}

	c.Status(fiber.StatusUnauthorized)
	return errors.New("unauthorized")
}
