package routes

import (
	"server/cmd/api/handlers"
	"server/cmd/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	// public endpoints
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// auth endpoints
	app.Use(middleware.IsAuthenticated)

	app.Put("/users/info", handlers.UpdateInfo)
	app.Put("/users/password", handlers.UpdatePassword)

	app.Post("/logout", handlers.Logout)
	app.Get("/user", handlers.User)

	// define static content
	app.Static("/uploads", "./uploads")

}
