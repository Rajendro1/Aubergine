package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// For production, this should be stored securely in environment variables.
var jwtSecret = []byte(os.Getenv("JwtSecret"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Plan   string `json:"plan"`
	jwt.RegisteredClaims
}

// GenerateToken generates both access and refresh tokens.
func GenerateToken(userID uint, plan string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 1 day validity

	claims := &Claims{
		UserID: userID,
		Plan:   plan,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses the token, checks the signature and validity
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
