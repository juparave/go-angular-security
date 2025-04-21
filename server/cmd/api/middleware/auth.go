package middleware

import (
	"net/http"
	"server/internal/util"

	"github.com/gofiber/fiber/v2"
)

// IsAuthenticated is a middleware function that checks if the user is authenticated.
// It retrieves the JWT from the request (either from the "Token" header or the "jwt" cookie).
// It then attempts to parse the JWT. If parsing fails or the token is invalid,
// it returns a 401 Unauthorized response with a JSON message "unauthenticated".
// If the JWT is valid, it calls the next handler in the chain using c.Next().
// Note: Currently, it does not verify if the user ID from the token exists in the database.
func IsAuthenticated(c *fiber.Ctx) error {
	jwt := util.GetJWT(c)

	if _, err := util.ParseJwt(jwt); err != nil {
		c.Status(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "not authenticated",
		})
	}

	// TODO: we should test if authenticated user exists on db

	return c.Next()
}
