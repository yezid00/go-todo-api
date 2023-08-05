package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitializeRoutes(app *fiber.App) {
	api := app.Group("/api/v1", logger.New())

	todoRoutes(api)

	authRoutes(api)
}
