package database

import (
	"log"
	"todo-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var DB DBInstance

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Database connection failed")
	}

	log.Println("Database connected.")

	db.Logger = logger.Default.LogMode((logger.Info))

	log.Println("Running Migration")

	db.AutoMigrate(&models.User{}, models.Todo{})

	DB = DBInstance{
		Db: db,
	}
}
