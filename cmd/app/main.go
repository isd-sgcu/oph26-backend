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
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// User & Staff
	userRepo := repository.NewUserRepository(config.DB)
	staffRepo := repository.NewStaffRepository(config.DB)
	// Auth
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseConfig{
		UserRepository:         userRepo,
		StaffRepository:        staffRepo,
		RefreshTokenRepository: refreshTokenRepo,
		GoogleClientID:         cfg.GoogleClientID,
		JWTSecret:              cfg.JWTSecret,
		AppEnv:                 cfg.AppEnv,
	})

	// Attendee
	attendeeRepo := repository.NewAttendeeRepository(config.DB)
	leaderboardRepo := repository.NewLeaderboardRepository(config.DB)
	scoreRepo := repository.NewScoreRepository(config.DB)
	pieceRepo := repository.NewPieceRepository(config.DB)
	checkinRepo := repository.NewCheckinRepository(config.DB)
	// Checkin
	userUsecase := usecase.NewUserUsecase(userRepo)
	checkinUsecase := usecase.NewCheckinUsecase(attendeeRepo, staffRepo, checkinRepo)
	attendeeUsecase := usecase.NewAttendeeUsecase(attendeeRepo, userRepo, leaderboardRepo, scoreRepo)
	leaderboardUsecase := usecase.NewLeaderboardUsecase(leaderboardRepo, scoreRepo)
	pieceUsecase := usecase.NewPieceUsecase(pieceRepo, leaderboardUsecase, scoreRepo)

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimitMiddleware := middleware.RateLimitMiddleware(10, time.Minute) // 10 requests per minute

	route.SetupRoutes(r, route.RouteConfig{
		AuthUsecase:         authUsecase,
		AttendeeUsecase:     attendeeUsecase,
		CheckinUsecase:      checkinUsecase,
		AuthMiddleware:      authMiddleware,
		UserUsecase:         userUsecase,
		PieceUsecase:        pieceUsecase,
		RateLimitMiddleware: rateLimitMiddleware,
		LeaderboardUsecase:  leaderboardUsecase,
	})

	log.Fatal(r.Listen(":" + cfg.Port))
}
