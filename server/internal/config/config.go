package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	lconfig "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/joho/godotenv"
	"github.com/juparave/mylogger"
)

const (
	// configFileName is the expected name for the JSON configuration file.
	configFileName = "config.json"
)

// AppConfig holds the application's static configuration values loaded from
// environment variables and/or a JSON file.
type AppConfig struct {
	Name        string `env:"APP_NAME" json:"name"`
	Mode        string `env:"APP_MODE" json:"mode"` // e.g., "development", "production"
	Port        int    `env:"APP_PORT" json:"port"`
	UploadsPath string `env:"UPLOADS_PATH" json:"uploads_path"`
	Log         *mylogger.MyLogger
	MailChan    chan any
	Database    DatabaseConfig `json:"database"`
	Email       EmailConfig
	Stripe      struct {
		SecretKey      string `env:"STRIPE_SECRET_KEY"`
		PublishableKey string `env:"STRIPE_PUBLISHABLE_KEY"`
		WebhookSecret  string `env:"STRIPE_WEBHOOK_SECRET"`
	}
}

type DatabaseConfig struct {
	Driver string `env:"DATABASE_DRIVER" json:"driver"`
	DSN    string `env:"DATABASE_DSN" json:"dsn"`
	Path   string `env:"DATABASE_PATH" json:"path"`
}

// EmailConfig holds email-related configuration.
type EmailConfig struct {
	Host        string `env:"EMAIL_HOST"`
	Port        string `env:"EMAIL_PORT"`
	Account     string `env:"EMAIL_ACCOUNT"`
	Password    string `env:"EMAIL_PASSWORD"`
	ContactTo   string `env:"EMAIL_CONTACT_TO"`
	ContactFrom string `env:"EMAIL_CONTACT_FROM"`
}

// GetAppConfig loads configuration from environment variables and optionally
// from a 'config.json' file. Environment variables take precedence.
// It returns a populated AppConfig struct.
func GetAppConfig() *AppConfig {
	cfg := AppConfig{}

	// 1. Load .env file if it exists (non-critical)
	// Useful for local development. In production, rely on actual env vars.
	err := godotenv.Load()
	if err != nil {
		// This is not necessarily an error, could be running in an environment
		// where env vars are set directly (like Docker, K8s).
		log.Println("Info: No local .env file found or loaded.")
	} else {
		log.Println("Info: Loaded configuration from .env file.")
	}

	// 1.1 Load default values, these values will be replaced, if needed with the file
	cfg.Name = getEnv("APP_NAME", "GO-ANGULAR-SECURITY").(string)
	cfg.Mode = getEnv("APP_MODE", "development").(string)
	cfg.Port = getEnv("APP_PORT", 5000).(int)

	// 2. Prepare configuration feeders
	var feeders []lconfig.Feeder

	// JSON Feeder (Optional)
	if _, err := os.Stat(configFileName); err == nil {
		// config.json exists, add it as a feeder
		log.Printf("Info: Found %s, adding as configuration source.\n", configFileName)
		feeders = append(feeders, &feeder.Json{Path: configFileName})
	} else if !os.IsNotExist(err) {
		// Another error occurred trying to stat the file (e.g., permissions)
		log.Printf("Warning: Could not stat %s: %v\n", configFileName, err)
	} else {
		log.Printf("Info: No %s file found.\n", configFileName)
	}

	// Environment Variable Feeder (Primary)
	// Env vars will override values from config.json if keys match.
	feeders = append(feeders, feeder.Env{})

	// 3. Create config loader and feed the struct
	loader := lconfig.New()
	loader.AddFeeder(feeders...) // Add all prepared feeders
	loader.AddStruct(&cfg)       // Register the struct to be populated

	err = loader.Feed()
	if err != nil {
		// Log the error, but consider if panicking might be more appropriate
		// if the application cannot run without valid configuration.
		log.Printf("Error: Failed to feed configuration: %v\n", err)
		// Depending on requirements, you might os.Exit(1) here.
	}

	// 4. Log the loaded configuration (optional, for debugging)
	// Be cautious about logging sensitive data like passwords in production.
	log.Println("Info: Configuration loaded successfully.")
	configJSON, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println("--- Loaded Configuration ---")
	fmt.Println(string(configJSON))
	fmt.Println("--------------------------")

	return &cfg
}

// Helper function to get env var or return default, now it can be string, int or float
func getEnv(key string, fallback interface{}) interface{} {
	if value, exists := os.LookupEnv(key); exists {
		switch fallback.(type) {
		case string:
			return value
		case int:
			var intValue int
			fmt.Sscan(value, &intValue)
			return intValue
		case float64:
			var floatValue float64
			fmt.Sscan(value, &floatValue)
			return floatValue
		default:
			return value
		}
	}
	return fallback
}
