package routes

import (
	todoController "todo-api/controllers/todo"
	"todo-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func todoRoutes(api fiber.Router) {
	todo := api.Group("/todos")

	todo.Get("", middleware.Protected(), todoController.Todos)
	todo.Post("", middleware.Protected(), todoController.Create)
	todo.Patch("/completed/:id", middleware.Protected(), todoController.MarkAsCompleted)
	todo.Patch("/not-completed/:id", middleware.Protected(), todoController.MarkAsNotCompleted)
}
