package routes

import (
	"go-app/controllers"
	"go-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// public endpoints
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

	// auth endpoints
	app.Use(middleware.IsAuthenticated)

	app.Put("/api/users/info", controllers.UpdateInfo)
	app.Put("/api/users/password", controllers.UpdatePassword)

	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/user", controllers.User)

	// define static content
	app.Static("/api/uploads", "./uploads")

}
