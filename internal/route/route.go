package route

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"oph26-backend/internal/repository"
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

	// Initialize Repository (Data Layer)
	userRepo := repository.NewUserRepository()

	// Initialize Use Case (Business Logic Layer)
	userUsecase := usecase.NewUserUsecase(userRepo)

	api := r.Group("/api")
	{
		api.Get("/ping", usecase.Ping)
		// Delivery Layer: HTTP handlers that call Use Cases
		api.Get("/users", userUsecase.GetAllUsers)
	}
}
