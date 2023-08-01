package main

import (
	"log"
	"todo-api/database"
	"todo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	routes.InitializeRoutes(app)

	app.Use(cors.New())

	app.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendString("🚀 App is running...")
	})

	log.Fatal(app.Listen("localhost:5000"))

}
