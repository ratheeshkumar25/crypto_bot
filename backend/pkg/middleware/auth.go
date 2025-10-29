package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates JWT tokens
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			c.Locals("user_id", claims.UserID)
			c.Locals("username", claims.Username)
			return c.Next()
		}

		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}
}

// GetUserIDFromContext extracts user ID from Fiber context
func GetUserIDFromContext(c *fiber.Ctx) int {
	if userID, ok := c.Locals("user_id").(int); ok {
		return userID
	}
	return 0
}
