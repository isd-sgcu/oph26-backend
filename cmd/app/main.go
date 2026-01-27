package main

import (
	"oph26-backend/internal/config"
	"oph26-backend/internal/initializer"
	"oph26-backend/internal/route"

	"github.com/gofiber/fiber/v2"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	cfg := config.LoadEnv()
	config.InitDB(cfg)

	r := fiber.New()

	route.SetupRoutes(r)

	r.Listen(":8080")
}
