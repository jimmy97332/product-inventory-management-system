package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("my_secret_key") // Define a key used to encrypt and decrypt the JWT

// Claims JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var GenerateJWT = realGenerateJWT

// Generate JWT token
func realGenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Set JWT expired time
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
