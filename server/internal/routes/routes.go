package routes

import (
	"server/internal/handlers"
	"server/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup sets up the API routes.
// It registers public and authenticated routes using the Fiber framework.
func Setup(app fiber.Router) {
	// Public endpoints: accessible without authentication.
	app.Post("/register", handlers.Register) // Registers a new user.
	app.Post("/login", handlers.Login)       // Logs in an existing user.

	// Protected endpoints: require authentication.
	app.Use(middleware.IsAuthenticated) // Applies authentication middleware.

	app.Put("/users/info", handlers.UpdateInfo)         // Updates user information.
	app.Put("/users/password", handlers.UpdatePassword) // Updates user password.

	app.Post("/logout", handlers.Logout) // Logs out the current user.
	app.Get("/user", handlers.User)      // Retrieves information about the current user.

	// Subscription endpoints:
	app.Get("/subscriptions/current", handlers.GetCurrentSubscription)
	app.Patch("/subscriptions/:id", handlers.PatchSubscription)
	app.Post("/subscriptions/:id/cancel", handlers.PostCancelSubscription)
	app.Post("/subscriptions/:id/reactivate", handlers.PostReactivateSubscription)
	app.Post("/subscriptions/:id/change-plan", handlers.PostChangeSubscription)

	// Serve static files from the uploads directory.
	app.Static("/uploads", "./uploads")
}
