package main

import (
	env "oph26-backend/internal/config"
	"oph26-backend/internal/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := env.Load()

	r := fiber.New()

	route.SetupRoutes(r)

	r.Listen(":" + cfg.Port)
}