package routes

import (
	"go-app/controllers"
	"go-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	// public endpoints
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	// auth endpoints
	app.Use(middleware.IsAuthenticated)

	app.Put("/users/info", controllers.UpdateInfo)
	app.Put("/users/password", controllers.UpdatePassword)

	app.Post("/logout", controllers.Logout)
	app.Get("/user", controllers.User)

	// define static content
	app.Static("/uploads", "./uploads")

}
