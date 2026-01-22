package middleware

import (
	// "backend/res/utils"
	"fmt"
	// "net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(401, echo.Map{"error": "Authorization header is missing"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(401, echo.Map{"error": "Invalid token format"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(401, echo.Map{"error": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		return next(c)
	}
}
