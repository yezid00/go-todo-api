package main

import (
	"log"
	"todo-api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New())
	log.Fatal(app.Listen(":8000"))

}
