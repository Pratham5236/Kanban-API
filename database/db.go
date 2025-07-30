package database

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban-app/api/models"
)

var (
	DB       *gorm.DB
	Enforcer *casbin.Enforcer
)

func ConnectDatabase() {
	var err error

	db, err := gorm.Open(sqlite.Open("kanban.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established to kanban.db")

	err = db.AutoMigrate(&models.User{}, &models.Organization{}, &models.Project{}, &models.Board{}, &models.List{}, &models.Card{}, &models.Label{}, &models.Comment{}, &models.Attachment{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database schema migrated successfully")

	DB = db

	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		log.Fatalf("Failed to create casbin adapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer("auth/casbin_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	Enforcer = enforcer
	log.Println("Casbin enforcer initialized")
}
