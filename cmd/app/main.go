package main

import (
	"log"
	"oph26-backend/internal/config"
	"oph26-backend/internal/initializer"
	"oph26-backend/internal/middleware"
	"oph26-backend/internal/repository"
	"oph26-backend/internal/route"
	"oph26-backend/internal/usecase"
	"time"

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
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	authUsecase := usecase.NewAuthUsecase(userRepo, staffRepo, refreshTokenRepo, cfg.GoogleClientID, cfg.JWTSecret, cfg.AppEnv)

	pieceRepo := repository.NewPieceRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	attendeeRepo := repository.NewAttendeeRepository(config.DB)
	attendeeUsecase := usecase.NewAttendeeUsecase(attendeeRepo, userRepo)
	pieceUsecase := usecase.NewPieceUsecase(pieceRepo)

	leaderboardRepo := repository.NewLeaderboardRepository(config.DB)
	scoreRepo := repository.NewScoreRepository(config.DB)
	leaderboardUsecase := usecase.NewLeaderboardUsecase(leaderboardRepo, scoreRepo)

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimitMiddleWare := middleware.RateLimitMiddleware(10, time.Minute) // 10 requests per minute

	route.SetupRoutes(r, route.RouteConfig{
		AuthUsecase:         authUsecase,
		AttendeeUsecase:     attendeeUsecase,
		AuthMiddleware:      authMiddleware,
		UserUsecase:         userUsecase,
		PieceUsecase:        pieceUsecase,
		LeaderboardUsecase: leaderboardUsecase,
		RateLimitMiddleware: rateLimitMiddleWare,
	})

	log.Fatal(r.Listen(":8080"))
}
