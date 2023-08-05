package main

import (
	"log"
	"todo-api/database"
	"todo-api/middleware"
	"todo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	routes.InitializeRoutes(app)

	app.Use(cors.New())

	app.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendString("🚀 App is running...")
	})

	app.Get("/api/v1/test", middleware.Protected(), func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		id := int(claims["user_id"].(float64))

		return c.SendString(id)
	})

	log.Fatal(app.Listen("localhost:5000"))

}
