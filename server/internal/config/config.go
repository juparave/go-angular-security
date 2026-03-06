package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/juparave/mylogger"
)

const configFileName = "config.json"

// AppConfig holds the application's static configuration values loaded from
// environment variables and/or a JSON file.
type AppConfig struct {
	Name        string `env:"APP_NAME" json:"name"`
	Version     string `env:"APP_VERSION" json:"version"`
	Mode        string `env:"APP_MODE" json:"mode"`
	Port        int    `env:"APP_PORT" json:"port"`
	Domain      string `env:"APP_DOMAIN" json:"domain"`
	UploadsPath string `env:"UPLOADS_PATH" json:"uploads_path"`
	Log         *mylogger.MyLogger
	MailChan    chan any       `json:"-"`
	Database    DatabaseConfig `json:"database"`
	Email       EmailConfig
	JWT         JWTConfig
	Stripe      struct {
		SecretKey      string `env:"STRIPE_SECRET_KEY"`
		PublishableKey string `env:"STRIPE_PUBLISHABLE_KEY"`
		WebhookSecret  string `env:"STRIPE_WEBHOOK_SECRET"`
	}
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret          string `env:"JWT_SECRET"`
	RefreshSecret   string `env:"JWT_REFRESH_SECRET"`
	ResetSecret     string `env:"JWT_RESET_SECRET"`
	AccessTokenTTL  int    `env:"JWT_ACCESS_TOKEN_TTL"`
	RefreshTokenTTL int    `env:"JWT_REFRESH_TOKEN_TTL"`
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

// GetAppConfig loads configuration from environment variables. A .env file and
// config.json are optional — errors in either are logged as warnings and
// env vars always take precedence.
func GetAppConfig() *AppConfig {
	// 1. Load .env file (optional — errors are warnings, not fatal)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file loaded:", err)
	} else {
		log.Println("Info: Loaded configuration from .env file.")
	}

	// 2. Warn if config.json exists but is invalid — don't fail
	if _, err := os.Stat(configFileName); err == nil {
		log.Printf("Info: Found %s (note: env vars take precedence).\n", configFileName)
	}

	// 3. Build config purely from env vars with sensible defaults
	cfg := &AppConfig{}
	cfg.Name = envStr("APP_NAME", "GO-ANGULAR-SECURITY")
	cfg.Mode = envStr("APP_MODE", "development")
	cfg.Port = envInt("APP_PORT", 5000)
	cfg.Domain = envStr("APP_DOMAIN", "")
	cfg.UploadsPath = envStr("UPLOADS_PATH", "./uploads")

	cfg.Database.Driver = envStr("DATABASE_DRIVER", "sqlite")
	cfg.Database.DSN = envStr("DATABASE_DSN", "")
	cfg.Database.Path = envStr("DATABASE_PATH", "gorm.db")

	cfg.JWT.Secret = envStr("JWT_SECRET", "")
	cfg.JWT.RefreshSecret = envStr("JWT_REFRESH_SECRET", "")
	cfg.JWT.ResetSecret = envStr("JWT_RESET_SECRET", "")
	cfg.JWT.AccessTokenTTL = envInt("JWT_ACCESS_TOKEN_TTL", 24)
	cfg.JWT.RefreshTokenTTL = envInt("JWT_REFRESH_TOKEN_TTL", 720)

	cfg.Email.Host = envStr("EMAIL_HOST", "")
	cfg.Email.Port = envStr("EMAIL_PORT", "587")
	cfg.Email.Account = envStr("EMAIL_ACCOUNT", "")
	cfg.Email.Password = envStr("EMAIL_PASSWORD", "")
	cfg.Email.ContactTo = envStr("EMAIL_CONTACT_TO", "")
	cfg.Email.ContactFrom = envStr("EMAIL_CONTACT_FROM", "")

	cfg.Stripe.SecretKey = envStr("STRIPE_SECRET_KEY", "")
	cfg.Stripe.PublishableKey = envStr("STRIPE_PUBLISHABLE_KEY", "")
	cfg.Stripe.WebhookSecret = envStr("STRIPE_WEBHOOK_SECRET", "")

	// 4. Validate required fields
	missing := []string{}
	if cfg.JWT.Secret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		log.Fatalf("Missing required configuration: %v", missing)
	}

	// 5. Set default version
	if cfg.Version == "" {
		cfg.Version = "0.1.0"
	}

	log.Println("Info: Configuration loaded successfully.")
	fmt.Printf("--- App: %s | Mode: %s | Port: %d ---\n", cfg.Name, cfg.Mode, cfg.Port)

	return cfg
}

func envStr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		log.Printf("Warning: %s=%q is not a valid integer, using default %d\n", key, v, fallback)
	}
	return fallback
}
