package route

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"oph26-backend/internal/usecase"
)

type RouteConfig struct {
	AuthUsecase         usecase.AuthUsecase
	AttendeeUsecase     usecase.AttendeeUsecase
	CheckinUsecase      usecase.CheckinUsecase
	PieceUsecase        usecase.PieceUsecase
	UserUsecase         usecase.UserUsecase
	LeaderboardUsecase  usecase.LeaderboardUsecase
	StatsUsecase        usecase.StatsUsecase
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

	r.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok": true,
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
			attendees.Post("/", c.AttendeeUsecase.PostAttendee)
			attendees.Get("/me", c.AttendeeUsecase.GetMyAttendee)
			attendees.Put("/me", c.AttendeeUsecase.PutAttendee)
			attendees.Put("/me/certificate_name", c.AttendeeUsecase.UpdateCertificateName)
		}

		pieces := api.Group("/pieces", c.AuthMiddleware)
		{
			pieces.Get("/me", c.PieceUsecase.GetMyPiece)
			pieces.Get("/me/collected", c.PieceUsecase.GetCollectedPieces)
			pieces.Post("/me/collected", c.PieceUsecase.CollectPiece)
		}

		leaderboards := api.Group("/leaderboards", c.AuthMiddleware)
		{
			leaderboards.Get("/me", c.LeaderboardUsecase.GetMyLeaderboard)
		}

		favWorkshop := api.Group("/favorite_workshops", c.AuthMiddleware)
		{
			favWorkshop.Get("/me", c.AttendeeUsecase.GetMyFavWorkshops)
			favWorkshop.Put("/me", c.AttendeeUsecase.PutMyFavWorkshops)
		}

		api.Post("/checkin", c.AuthMiddleware, c.CheckinUsecase.CheckIn)

		stats := api.Group("/stats", c.RateLimitMiddleware)
		{
			stats.Get("/stats/attendees", c.StatsUsecase.GetCountAttendeesStats)
			stats.Get("/stats/attendees/checkins_by_date", c.StatsUsecase.GetCountUniqueAttendeesCheckinsGroupedByDateStats)
			stats.Get("/stats/pieces/available_by_faculty", c.StatsUsecase.GetCountAvailablePiecesGroupedByFacultyStats)
		}
	}
}
