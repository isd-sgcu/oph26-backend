package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               maxRequests,
		Expiration:        window,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			if ip := c.Get("CF-Connecting-IP"); ip != "" {
				return ip
			}
			if ip := c.Get("X-Forwarded-For"); ip != "" {
				return ip
			}
			return c.IP()
		},
	})
}
