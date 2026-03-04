package main

import (
	"log"
	"oph26-backend/internal/config"
	"oph26-backend/internal/initializer"
	"oph26-backend/internal/middleware"
	"oph26-backend/internal/repository"
	"oph26-backend/internal/route"
	"oph26-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	cfg := config.LoadEnv()
	config.InitDB(cfg)

	switch cfg.AppEnv {
	case "production":
		log.Println("Running in PRODUCTION mode")
	case "development":
		log.Println("Running in development mode")
	default:
		log.Printf("Running in unknown mode: %s\n", cfg.AppEnv)
	}

	r := fiber.New()

	// Init Dependencies
	userRepo := repository.NewUserRepository(config.DB)
	staffRepo := repository.NewStaffRepository(config.DB)
	authUsecase := usecase.NewAuthUsecase(userRepo, staffRepo, cfg.GoogleClientID, cfg.JWTSecret, cfg.AppEnv)

	attendeeRepo := repository.NewAttendeeRepository(config.DB)
	attendeeUsecase := usecase.NewAttendeeUsecase(attendeeRepo)

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	route.SetupRoutes(r, authUsecase, attendeeUsecase, authMiddleware)

	log.Fatal(r.Listen(":8080"))
}
