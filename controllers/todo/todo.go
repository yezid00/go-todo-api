package todo

import "github.com/gofiber/fiber"

func Todos(c *fiber.Ctx) error {
	user, err := auth.AuthenticatedUser(c)
}
