package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewRequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID := uuid.New().String()[:8]

		c.Locals("request_id", reqID)

		err := c.Next()

		status := c.Response().StatusCode()
		latency := time.Since(start)

		attrs := []any{
			slog.String("req_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("ip", c.IP()),
		}

		// Add user context if available (set by auth middleware)
		if userID, ok := c.Locals("user_id").(uuid.UUID); ok {
			attrs = append(attrs, slog.String("user_id", userID.String()))
		}
		if role, ok := c.Locals("role").(string); ok {
			attrs = append(attrs, slog.String("role", role))
		}

		switch {
		case status >= 500:
			slog.Error("request", attrs...)
		case status >= 400:
			slog.Warn("request", attrs...)
		default:
			slog.Info("request", attrs...)
		}

		return err
	}
}
