package main

import (
	"fmt"
	"server/internal/config"
	"server/internal/database"
	"server/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/juparave/mylogger"
)

var app *config.AppConfig

func main() {
	// initialize the configuration
	app = config.GetAppConfig()
	app.Log = mylogger.NewLogger()

	app.Log.Info("Starting application setup...") // Optional: Log startup

	database.Connect()

	server := fiber.New(fiber.Config{
		AppName: app.Name,
	})

	// let frontend get and store the auth cookie from diferent ports
	server.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200", // let angular frontend
		AllowCredentials: true,
	}))

	// group handlers, set api prefix path to /api/v1
	api := server.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	routes.Setup(v1)

	app.Log.Info("Starting server on port 5000")
	server.Listen(fmt.Sprintf(":%d", app.Port))
}
