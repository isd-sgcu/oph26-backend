package route

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"oph26-backend/internal/usecase"
)

var startTime = time.Now()

func SetupRoutes(r *fiber.App, authUsecase usecase.AuthUsecase, attendeeUsecase usecase.AttendeesUsecase, authMiddleware fiber.Handler) {
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

		// Example of protected route (as per requirement, but applies generally)
		// attendees := api.Group("/attendees", authMiddleware)
		// attendees.Post("/", ...)
		// Delivery Layer: HTTP handlers that call Use Cases
		// api.Get("/users", userUsecase.GetAllUsers)

		auth := api.Group("/auth")
		{
			auth.Post("/token", authUsecase.Login)
		}

		attendees := api.Group("/attendees", authMiddleware)
		{
			attendees.Get("/me", attendeeUsecase.GetMyAttendee)
			attendees.Get("/:attendeeId", attendeeUsecase.GetByAttendeeId)
		}
	}
}
