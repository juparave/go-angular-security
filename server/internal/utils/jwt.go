package utils

import (
	"fmt"
	"log"
	"server/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtConfig struct {
	secret          string
	refreshSecret   string
	resetSecret     string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	initialized     bool
}

// InitJWT initializes JWT configuration from app config values.
// Must be called before any JWT operations.
func InitJWT(secret, refreshSecret, resetSecret string, accessTTLHours, refreshTTLHours int) {
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if refreshSecret == "" {
		refreshSecret = secret + "-refresh"
	}
	if resetSecret == "" {
		resetSecret = secret + "-reset"
	}
	if accessTTLHours <= 0 {
		accessTTLHours = 24 // 1 day default
	}
	if refreshTTLHours <= 0 {
		refreshTTLHours = 720 // 30 days default
	}

	jwtConfig.secret = secret
	jwtConfig.refreshSecret = refreshSecret
	jwtConfig.resetSecret = resetSecret
	jwtConfig.accessTokenTTL = time.Hour * time.Duration(accessTTLHours)
	jwtConfig.refreshTokenTTL = time.Hour * time.Duration(refreshTTLHours)
	jwtConfig.initialized = true
}

// AccessTokenDuration returns the configured access token TTL
func AccessTokenDuration() time.Duration {
	if !jwtConfig.initialized {
		return time.Hour * 24
	}
	return jwtConfig.accessTokenTTL
}

// RefreshTokenDuration returns the configured refresh token TTL
func RefreshTokenDuration() time.Duration {
	if !jwtConfig.initialized {
		return time.Hour * 720
	}
	return jwtConfig.refreshTokenTTL
}

// GetResetSecret returns the configured reset token secret
func GetResetSecret() string {
	if !jwtConfig.initialized {
		log.Fatal("JWT not initialized: call InitJWT first")
	}
	return jwtConfig.resetSecret
}

func getSecret() string {
	if !jwtConfig.initialized {
		log.Fatal("JWT not initialized: call InitJWT first")
	}
	return jwtConfig.secret
}

func getRefreshSecret() string {
	if !jwtConfig.initialized {
		log.Fatal("JWT not initialized: call InitJWT first")
	}
	return jwtConfig.refreshSecret
}

func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.accessTokenTTL)),
	})

	return claims.SignedString([]byte(getSecret()))
}

func GenerateRefreshJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.refreshTokenTTL)),
	})

	return claims.SignedString([]byte(getRefreshSecret()))
}

// GetJWT returns jwt from `Authorization` header or from `jwt` cookie
func GetJWT(c *fiber.Ctx) string {
	var jwtToken string

	token := c.Get("Token")    // get value from headers
	cookie := c.Cookies("jwt") // get value from cookie

	if token != "" {
		t := strings.Split(token, " ")
		if len(t) == 2 {
			jwtToken = t[1]
		} else {
			jwtToken = token
		}
	} else {
		jwtToken = cookie
	}

	return jwtToken
}

func ParseJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getSecret()), nil
	})

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	return claims.Issuer, nil
}

// ParseJwtWithSecret parses a JWT token with a specific secret (used for refresh tokens)
func ParseJwtWithSecret(cookie string, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	return claims.Issuer, nil
}

// GenerateUserTokens sets token and refreshToken to user
func GenerateUserTokens(user *models.User) error {
	accessToken, err := GenerateJwt(user.ID)
	if err != nil {
		return err
	}

	refreshToken, err := GenerateRefreshJwt(user.ID)
	if err != nil {
		return err
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return nil
}

func GenerateResetPasswordToken(user *models.User) (string, error) {
	expiry := time.Now().Add(time.Hour * 24)
	token := fmt.Sprintf("%s|%s", user.Email, expiry.Format(time.RFC3339))
	return Encrypt(token, getSecret())
}

func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	jwtToken := GetJWT(c)
	return ParseJwt(jwtToken)
}
