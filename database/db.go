package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban-app/api/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	db, err := gorm.Open(sqlite.Open("kanban.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established to kanban.db")

	err = db.AutoMigrate(&models.User{}, &models.Organization{}, &models.Project{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database schema migrated successfully")

	DB = db
}
