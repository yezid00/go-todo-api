package routes

import (
	authController "todo-api/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func authRoutes(api fiber.Router) {

	auth := api.Group("/auth")

	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)

	// auth.Post("/login", authController.Login)
	// auth.Post("/register", authController.Register)
}
