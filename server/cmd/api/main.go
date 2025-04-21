package main

import (
	"server/cmd/api/routes"
	"server/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connect()

	app := fiber.New()

	// let frontend get and store the auth cookie from diferent ports
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200", // let angular frontend
		AllowCredentials: true,
	}))

	// group handlers, set api prefix path to /api/v1
	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	routes.Setup(v1)

	app.Listen(":5000")
}
