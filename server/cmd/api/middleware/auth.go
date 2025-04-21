package middleware

import (
	"server/internal/util"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	jwt := util.GetJWT(c)

	if _, err := util.ParseJwt(jwt); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// we should test if authenticated user exists on db

	return c.Next()
}
