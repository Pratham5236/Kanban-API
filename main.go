// @title Kanban API Documentation
// @version 1.0
// @description This is the API documentation for the Kanban application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter your JWT token in the format "Bearer <your_token>".

// @tag.name Authentication
// @tag.description "User login and registration"
// @tag.name Organizations
// @tag.description "Operations related to organizations"
// @tag.name Projects
// @tag.description "Operations related to projects within organizations"
// @tag.name Boards
// @tag.description "Operations related to Kanban boards within projects"
// @tag.name Lists
// @tag.description "Operations related to lists (columns) within boards"
// @tag.name Cards
// @tag.description "Operations related to cards (tasks) within lists"
package main

import (
	"kanban-app/api/controllers"
	"kanban-app/api/database"
	_ "kanban-app/api/docs"
	"kanban-app/api/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "supersecretjwtkeythatshouldbemoresecureandlongerinproduction")
		log.Println("WARNING: JWT_SECRET environment variable not set. Using a default development key.")
	}

	database.ConnectDatabase()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	router.GET("/health", controllers.HealthCheck)
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	authenticated := router.Group("/api")
	authenticated.POST("/uploads", controllers.UploadFile)
	authenticated.Use(middlewares.AuthMiddleware())
	{
		// organization routes
		authenticated.POST("/organizations", controllers.CreateOrganization)
		authenticated.GET("/organizations", controllers.GetOrganizations)
		orgRoutes := authenticated.Group("/organizations/:orgID")
		orgRoutes.Use(middlewares.CasbinMiddleware("orgID", "owner"))
		{
			orgRoutes.GET("", controllers.GetOrganizationByID)
			orgRoutes.PUT("", controllers.UpdateOrganization)
			orgRoutes.DELETE("", controllers.DeleteOrganization)
		}

		// project routes
		projectRoutes := authenticated.Group("/organizations/:orgID/projects")
		projectRoutes.Use(middlewares.CasbinMiddleware("orgID", "owner"))
		{
			projectRoutes.POST("", controllers.CreateProject)
			projectRoutes.GET("", controllers.GetProjects)
		}

		projectDetailRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID")
		projectDetailRoutes.Use(middlewares.CasbinMiddleware("projectID", "owner"))
		{
			projectDetailRoutes.GET("", controllers.GetProjectByID)
			projectDetailRoutes.PUT("", controllers.UpdateProject)
			projectDetailRoutes.DELETE("", controllers.DeleteProject)
		}

		// Board routes (nested under projects)
		boardRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards")
		boardRoutes.Use(middlewares.CasbinMiddleware("projectID", "owner"))
		{
			boardRoutes.POST("", controllers.CreateBoard)
			boardRoutes.GET("", controllers.GetBoards)
		}

		boardDetailRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards/:boardID")
		boardDetailRoutes.Use(middlewares.CasbinMiddleware("boardID", "owner"))
		{
			boardDetailRoutes.GET("", controllers.GetBoardByID)
			boardDetailRoutes.PUT("", controllers.UpdateBoard)
			boardDetailRoutes.DELETE("", controllers.DeleteBoard)
			boardDetailRoutes.GET("/details", controllers.GetBoardDetails)
		}

		// List routes (nested under boards)
		listRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards/:boardID/lists")
		listRoutes.Use(middlewares.CasbinMiddleware("boardID", "owner"))
		{
			listRoutes.POST("", controllers.CreateList)
			listRoutes.GET("", controllers.GetLists)
		}

		listDetailRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards/:boardID/lists/:listID")
		listDetailRoutes.Use(middlewares.CasbinMiddleware("listID", "owner"))
		{
			listDetailRoutes.GET("", controllers.GetListByID)
			listDetailRoutes.PUT("", controllers.UpdateList)
			listDetailRoutes.DELETE("", controllers.DeleteList)
		}

		// Card routes (nested under lists)
		cardRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards/:boardID/lists/:listID/cards")
		cardRoutes.Use(middlewares.CasbinMiddleware("listID", "owner"))
		{
			cardRoutes.POST("", controllers.CreateCard)
			cardRoutes.GET("", controllers.GetCards)
		}

		cardDetailRoutes := authenticated.Group("/organizations/:orgID/projects/:projectID/boards/:boardID/lists/:listID/cards/:cardID")
		cardDetailRoutes.Use(middlewares.CasbinMiddleware("cardID", "owner"))
		{
			cardDetailRoutes.GET("", controllers.GetCardByID)
			cardDetailRoutes.PUT("", controllers.UpdateCard)
			cardDetailRoutes.DELETE("", controllers.DeleteCard)
		}

		// Comment routes (nested under cards)
		commentRoutes := authenticated.Group("/cards/:cardID/comments")
		commentRoutes.Use(middlewares.CasbinMiddleware("cardID", "owner"))
		{
			commentRoutes.POST("", controllers.CreateComment)
		}
		authenticated.DELETE("/cards/:cardID/comments/:commentID", middlewares.CasbinMiddleware("cardID", "owner"), controllers.DeleteComment)

		// Label routes
		labelRoutes := authenticated.Group("/labels")
		{
			labelRoutes.POST("", controllers.CreateLabel)
			labelRoutes.GET("", controllers.GetAllLabels)
			labelRoutes.GET("/:labelID", controllers.GetLabelByID)
			labelRoutes.PUT("/:labelID", controllers.UpdateLabel)
			labelRoutes.DELETE("/:labelID", controllers.DeleteLabel)
		}

		// Attachment routes (nested under cards)
		attachmentRoutes := authenticated.Group("/cards/:cardID/attachments")
		attachmentRoutes.Use(middlewares.CasbinMiddleware("cardID", "owner"))
		{
			attachmentRoutes.POST("", controllers.CreateAttachment)
		}
		authenticated.DELETE("/cards/:cardID/attachments/:attachmentID", middlewares.CasbinMiddleware("cardID", "owner"), controllers.DeleteAttachment)

		// Card-Label association routes // TODO implement these
		//authenticated.POST("/cards/:cardID/labels/:labelID", middlewares.CasbinMiddleware("cardID", "owner"), controllers.AddLabelToCard)
		//authenticated.DELETE("/cards/:cardID/labels/:labelID", middlewares.CasbinMiddleware("cardID", "owner"), controllers.RemoveLabelFromCard)

	}

	port := ":8080"
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
