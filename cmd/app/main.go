package main

import (
	"log"
	"log/slog"
	"oph26-backend/internal/config"
	"oph26-backend/internal/initializer"
	"oph26-backend/internal/metrics"
	"oph26-backend/internal/middleware"
	"oph26-backend/internal/repository"
	"oph26-backend/internal/route"
	"oph26-backend/internal/usecase"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	cfg := config.LoadEnv()
	config.InitDB(cfg)

	// Structured logging
	logLevel := slog.LevelInfo
	if cfg.AppEnv == "development" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	switch cfg.AppEnv {
	case "production":
		slog.Info("server starting", "mode", "production")
	case "development":
		slog.Info("server starting", "mode", "development")
	default:
		slog.Warn("server starting", "mode", cfg.AppEnv)
	}

	r := fiber.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization",
		AllowCredentials: true,
	}))

	r.Use(middleware.NewRequestLogger())

	serverRuntimeMetrics := metrics.NewServerRuntimeMetrics()
	r.Use(serverRuntimeMetrics.Middleware())

	// User & Staff
	userRepo := repository.NewUserRepository(config.DB)
	staffRepo := repository.NewStaffRepository(config.DB)
	// Auth
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	// Attendee
	attendeeRepo := repository.NewAttendeeRepository(config.DB)
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseConfig{
		UserRepository:         userRepo,
		StaffRepository:        staffRepo,
		AttendeeRepository:     attendeeRepo,
		RefreshTokenRepository: refreshTokenRepo,
		GoogleClientID:         cfg.GoogleClientID,
		JWTSecret:              cfg.JWTSecret,
		AppEnv:                 cfg.AppEnv,
	})
	leaderboardRepo := repository.NewLeaderboardRepository(config.DB)
	scoreRepo := repository.NewScoreRepository(config.DB)
	pieceRepo := repository.NewPieceRepository(config.DB)
	checkinRepo := repository.NewCheckinRepository(config.DB)
	// Checkin
	userUsecase := usecase.NewUserUsecase(userRepo)
	checkinUsecase := usecase.NewCheckinUsecase(attendeeRepo, staffRepo, checkinRepo)
	attendeeUsecase := usecase.NewAttendeeUsecase(attendeeRepo, userRepo, leaderboardRepo, scoreRepo)
	leaderboardUsecase := usecase.NewLeaderboardUsecase(leaderboardRepo, scoreRepo)
	pieceUsecase := usecase.NewPieceUsecase(pieceRepo, leaderboardUsecase, scoreRepo, attendeeRepo)

	// Stats
	statsRepo := repository.NewStatsRepository(config.DB)
	statsUsecase := usecase.NewStatsUsecase(statsRepo)

	// Questionnaire
	questionnaireRepo := repository.NewQuestionnaireRepository(config.DB)
	questionnaireUsecase := usecase.NewQuestionnaireUsecase(questionnaireRepo)
	attendeeMetrics := metrics.NewAttendeeMetrics(statsRepo)

	if err := attendeeMetrics.Refresh(); err != nil {
		slog.Error("initial attendee metrics refresh failed", "error", err)
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := attendeeMetrics.Refresh(); err != nil {
				slog.Error("periodic attendee metrics refresh failed", "error", err)
			}
		}
	}()

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimitMiddleware := middleware.RateLimitMiddleware(10, time.Minute) // 10 requests per minute

	metricsBasicAuthMiddleware := basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.MetricsBasicAuthUser: cfg.MetricsBasicAuthPass,
		},
	})

	route.SetupRoutes(r, route.RouteConfig{
		AuthUsecase:                authUsecase,
		AttendeeUsecase:            attendeeUsecase,
		CheckinUsecase:             checkinUsecase,
		AuthMiddleware:             authMiddleware,
		UserUsecase:                userUsecase,
		PieceUsecase:               pieceUsecase,
		StatsUsecase:               statsUsecase,
		RateLimitMiddleware:        rateLimitMiddleware,
		LeaderboardUsecase:         leaderboardUsecase,
		QuestionnaireUsecase:       questionnaireUsecase,
		MetricsBasicAuthMiddleware: metricsBasicAuthMiddleware,
	})

	log.Fatal(r.Listen(":" + cfg.Port))
}
