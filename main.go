package main

import (
	"kanban-app/api/controllers"
	"kanban-app/api/database"
	"kanban-app/api/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "supersecretjwtkeythatshouldbemoresecureandlongerinproduction")
		log.Println("WARNING: JWT_SECRET environment variable not set. Using a default development key.")
	}

	database.ConnectDatabase()

	router := gin.Default()

	router.GET("/health", controllers.HealthCheck)
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	authenticated := router.Group("/api")
	authenticated.Use(middlewares.AuthMiddleware())
	{
		// organization routes
		authenticated.POST("/organizations", controllers.CreateOrganization)
		authenticated.GET("/organizations", controllers.GetOrganizations)
		authenticated.GET("/organizations/:id", controllers.GetOrganizationByID)
		authenticated.PUT("/organizations/:id", controllers.UpdateOrganization)
		authenticated.DELETE("/organizations/:id", controllers.DeleteOrganization)

		// project routes
		authenticated.POST("/organizations/:orgID/projects", controllers.CreateProject)
		authenticated.GET("/organizations/:orgID/projects", controllers.GetProjects)
		authenticated.GET("/organizations/:orgID/projects/:projectID", controllers.GetProjectByID)
		authenticated.PUT("/organizations/:orgID/projects/:projectID", controllers.UpdateProject)
		authenticated.DELETE("/organizations/:orgID/projects/:projectID", controllers.DeleteProject)
	}

	port := ":8080"
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
