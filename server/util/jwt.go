package util

import (
	"fmt"
	"go-app/models"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// SecretKey, to encode jwt claims, change it
const SecretKey = "s3cr37kEy"

const AccessTokenDuration = time.Hour * 24 * 90
const RefreshTokenDuration = time.Hour * 24 * 180

func GenerateJwt(issuer string) (string, error) {
	// ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
	})

	return claims.SignedString([]byte(SecretKey))
}

func GenerateRefreshJwt(issuer string) (string, error) {
	// ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
	})

	return claims.SignedString([]byte(SecretKey))
}

// GetJWT returns jwt from `Authorization` header or from `jwt` cookie
func GetJWT(c *fiber.Ctx) string {
	var jwt string

	token := c.Get("Token")    // get value from headers
	cookie := c.Cookies("jwt") // get value from cookie

	if token != "" {
		t := strings.Split(token, " ")
		if len(t) == 2 {
			jwt = t[1]
		} else {
			jwt = token
		}
	} else {
		jwt = cookie
	}

	return jwt
}

func ParseJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err != nil || !token.Valid {
		return "", err
	}

	// cast to RegisteredClaims struct
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
	// add expiry time to token for 24 hours
	expiry := time.Now().Add(time.Hour * 24)
	// use user email with expiry time to generate token string
	token := fmt.Sprintf("%s|%s", user.Email, expiry.Format(time.RFC3339))
	// encrypt token string with secret key
	return Encrypt(token, SecretKey)
}
