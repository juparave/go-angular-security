package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"server/internal/config"
	"server/internal/database"
	"server/internal/handlers"
	"server/internal/routes"
	"server/internal/utils"
	"server/pkg/mail"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/juparave/mylogger"
)

var app *config.AppConfig

func main() {
	// Initialize configuration
	app = config.GetAppConfig()
	app.Log = mylogger.NewLogger()
	app.Log.Info("Starting application setup...")

	// Initialize JWT with config-driven secrets
	utils.InitJWT(
		app.JWT.Secret,
		app.JWT.RefreshSecret,
		app.JWT.ResetSecret,
		app.JWT.AccessTokenTTL,
		app.JWT.RefreshTokenTTL,
	)

	// Initialize database manager for multitenancy
	masterDSN := "gorm.db" // Default SQLite file
	if app.Database.Path != "" {
		masterDSN = app.Database.Path
	}

	if err := database.InitManager(masterDSN, "data"); err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer database.Manager.CloseAll()

	// Initialize email service if configured
	var emailService *mail.Service
	if app.Email.Host != "" {
		emailCfg := mail.NewConfig(
			app.Email.Host,
			app.Email.Port,
			app.Email.Account,
			app.Email.Password,
			app.Email.ContactFrom,
			app.Name,
		)
		emailService = mail.NewService(emailCfg, 100)
		defer emailService.Close()
		handlers.SetEmailService(emailService)
	}

	// Set app config for handlers
	handlers.SetRepo(app)

	// Create Fiber server
	server := fiber.New(fiber.Config{
		AppName: app.Name,
	})

	// Build CORS origins from config
	corsOrigins := app.Domain
	if app.Mode != "production" {
		corsOrigins = "http://localhost:4200,http://localhost:5000"
		if app.Domain != "" {
			corsOrigins = app.Domain + "," + corsOrigins
		}
	}

	// Middleware
	server.Use(recover.New())
	server.Use(logger.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Token",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Stripe webhook handler (needs raw body, so it's at root level)
	server.Post("/stripe/webhook", handlers.PostStripeWebhook)

	// API routes
	api := server.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	routes.Setup(v1)

	// Health check
	server.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": app.Version,
		})
	})

	// Start server in background goroutine
	app.Log.Info(fmt.Sprintf("Starting server on port %d", app.Port))
	go func() {
		if err := server.Listen(fmt.Sprintf(":%d", app.Port)); err != nil {
			app.Log.Error("Server stopped", "error", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	app.Log.Info("Shutting down server...")
	if err := server.Shutdown(); err != nil {
		app.Log.Error("Server shutdown error", "error", err)
	}
}
