package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid Authorization header format",
			})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token claims",
			})
		}

		// Extract and parse claims
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid user_id format",
			})
		}

		email, ok := claims["email"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid email claim",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid role claim",
			})
		}

		// Parse optional UUID fields
		var attendeeID *uuid.UUID
		var staffID *uuid.UUID

		if attendeeIDClaim, ok := claims["attendee_id"]; ok && attendeeIDClaim != nil {
			if attendeeIDStr, ok := attendeeIDClaim.(string); ok && attendeeIDStr != "" {
				if parsedID, err := uuid.Parse(attendeeIDStr); err == nil {
					attendeeID = &parsedID
				}
			}
		}

		if staffIDClaim, ok := claims["staff_id"]; ok && staffIDClaim != nil {
			if staffIDStr, ok := staffIDClaim.(string); ok && staffIDStr != "" {
				if parsedID, err := uuid.Parse(staffIDStr); err == nil {
					staffID = &parsedID
				}
			}
		}

		// Store flattened user information in locals
		c.Locals("user_id", userID)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("attendee_id", attendeeID)
		c.Locals("staff_id", staffID)

		return c.Next()
	}
}
