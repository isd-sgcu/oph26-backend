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
	})
}
