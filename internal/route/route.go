package route

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"oph26-backend/internal/usecase"
)

type RouteConfig struct {
	AuthUsecase         usecase.AuthUsecase
	AttendeeUsecase     usecase.AttendeeUsecase
	UserUsecase         usecase.UserUsecase
	PieceUsecase        usecase.PieceUsecase
	AuthMiddleware      fiber.Handler
	RateLimitMiddleware fiber.Handler
}

var startTime = time.Now()

func SetupRoutes(r *fiber.App, c RouteConfig) {
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

		auth := api.Group("/auth", c.RateLimitMiddleware)
		{
			auth.Post("/token", c.AuthUsecase.Login)
			auth.Get("/me", c.AuthMiddleware, c.AuthUsecase.GetCurrentUser)
			auth.Post("/refresh", c.AuthUsecase.RefreshToken)
			auth.Post("/signOut", c.AuthMiddleware, c.AuthUsecase.SignOut)
		}

		attendees := api.Group("/attendees", c.AuthMiddleware)
		{
			attendees.Get("/me", c.AttendeeUsecase.GetMyAttendee)
			attendees.Get("/:attendeeId", c.AttendeeUsecase.GetByAttendeeId)
		}
	}
}
