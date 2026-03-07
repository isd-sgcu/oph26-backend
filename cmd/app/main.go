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

	r := fiber.New()

	// Init Dependencies
	userRepo := repository.NewUserRepository(config.DB)
	pieceRepo := repository.NewPieceRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	attendeeRepo := repository.NewAttendeeRepository(config.DB)
	attendeeUsecase := usecase.NewAttendeeUsecase(userRepo, attendeeRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.GoogleClientID, cfg.JWTSecret)
	pieceUsecase := usecase.NewPieceUsecase(pieceRepo)

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	route.SetupRoutes(r, authUsecase, userUsecase, attendeeUsecase, pieceUsecase, authMiddleware)

	log.Fatal(r.Listen(":8080"))
}
