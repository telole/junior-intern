package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type JwtCustomClaims struct {
	UserID   uint   `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, tenantID string, email string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
