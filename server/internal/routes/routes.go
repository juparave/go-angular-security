package routes

import (
	"server/internal/handlers"
	"server/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup sets up the API routes.
// It registers public and authenticated routes using the Fiber framework.
func Setup(app fiber.Router) {
	// ===================
	// Public endpoints
	// ===================

	// Authentication (public)
	app.Post("/register", handlers.Register)                      // Legacy registration
	app.Post("/register-account", handlers.RegisterAccount)       // Self-service account registration
	app.Post("/login", handlers.Login)
	app.Post("/refresh-token", handlers.RefreshToken)

	// Password reset (public)
	app.Post("/request-password-reset", handlers.RequestPasswordReset)
	app.Post("/reset-password", handlers.ResetPassword)

	// Webhooks (public, but verified by Stripe signature)
	app.Post("/webhooks/stripe", handlers.HandleStripeWebhook)

	// ===================
	// Protected endpoints
	// ===================
	app.Use(middleware.IsAuthenticated)

	// User profile
	app.Get("/user", handlers.User)
	app.Post("/logout", handlers.Logout)
	app.Put("/users/info", handlers.UpdateInfo)
	app.Put("/users/password", handlers.UpdatePassword)

	// Change password (authenticated, requires current password)
	app.Post("/change-password", handlers.ChangePassword)

	// ===================
	// Team management (requires account)
	// ===================
	app.Get("/team", middleware.RequireAccount, handlers.GetTeamMembers)
	app.Post("/team", middleware.RequireEditor(), handlers.InviteTeamMember)
	app.Put("/team/:id", middleware.RequireEditor(), handlers.UpdateTeamMember)
	app.Delete("/team/:id", middleware.RequireAdmin(), handlers.DeleteTeamMember)
	app.Post("/team/:id/resend", middleware.RequireEditor(), handlers.ResendInvitation)

	// ===================
	// Subscriptions (requires account)
	// ===================
	app.Get("/subscriptions/current", middleware.RequireAccount, handlers.GetCurrentSubscription)
	app.Post("/subscriptions/create-checkout-session", middleware.RequireAccount, handlers.CreateCheckoutSession)
	app.Patch("/subscriptions/:id", middleware.RequireAccount, handlers.PatchSubscription)
	app.Post("/subscriptions/:id/cancel", middleware.RequireAccount, handlers.PostCancelSubscription)
	app.Post("/subscriptions/:id/reactivate", middleware.RequireAccount, handlers.PostReactivateSubscription)
	app.Post("/subscriptions/:id/change-plan", middleware.RequireAccount, handlers.PostChangeSubscription)

	// ===================
	// Static files
	// ===================
	app.Static("/uploads", "./uploads")
}
