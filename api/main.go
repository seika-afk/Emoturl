package main

import (
	"emoturl/routes"
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	port := strings.TrimPrefix(os.Getenv("APP_PORT"), ":")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
