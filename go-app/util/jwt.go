package util

import (
	"log"
	"time"

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
