package todo

import (
	"errors"
	"todo-api/controllers/auth"
	"todo-api/database"
	"todo-api/models"

	// "github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Todos(c *fiber.Ctx) error {
	user, err := auth.AuthenticatedUser(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "User not found",
		})
	}

	var todos []models.Todo

	db := database.DB.Db

	db.Preload("User").Where("user_id = ?", user.ID).Find(&todos)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "failed",
		"message": "Todos retrieved",
		"data":    todos,
	})
}

func Create(c *fiber.Ctx) error {
	type CreateTodoInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var input CreateTodoInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation errors",
		})
	}

	id := auth.ExtractUserId(c)

	user, err := auth.GetUserById(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	todo := models.Todo{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	todo.User = *user

	db := database.DB.Db

	if err := db.Create(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't create todo", "data": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "message": "Todo Created", "data": todo})
}

func MarkAsCompleted(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get Todo
	db := database.DB.Db

	var todo models.Todo

	if err := db.Preload("User").First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Todo not found", "data": nil})
		}
	}

	user, err := auth.AuthenticatedUser(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Could not get authenticated user", "data": nil})
	}

	// Check if authenticated user owns the todo item
	if user.ID != uint(todo.UserID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}

	todo.Completed = true

	db.Save(&todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Todo marked as completed", "data": todo})
}

func MarkAsNotCompleted(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get Todo
	db := database.DB.Db

	var todo models.Todo

	if err := db.Preload("User").First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Todo not found", "data": nil})
		}
	}

	user, err := auth.AuthenticatedUser(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Could not get authenticated user", "data": nil})
	}

	// Check if authenticated user owns the todo item
	if user.ID != uint(todo.UserID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}

	todo.Completed = false

	db.Save(&todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Todo marked as not completed", "data": todo})
}
