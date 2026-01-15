package route

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"oph26-backend/internal/usecase"
)

var startTime = time.Now()

func SetupRoutes(r *fiber.App) {
	r.Get("/healthz", func(c *fiber.Ctx) error {
		uptime := time.Since(startTime).String()
		return c.JSON(fiber.Map{
			"status": "up",
			"uptime": uptime,
		})
	})

	api := r.Group("/api")
	{
		api.Get("/ping", usecase.Ping)
	}
}