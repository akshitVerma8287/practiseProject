package main

import (
	"project/config"
	"project/controller"
	"project/models"
	"project/routes"
	"project/services"
	"github.com/joho/godotenv"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

_ "project/docs"

	"github.com/gin-gonic/gin"
)

// @title Student CRUD API
// @version 1.0
// @description This is a student CRUD API server.
// @host localhost:8080
// @BasePath /

func BuildStudentProvider() models.StudentProvider {
	return services.InitStudentService()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	config.ConnectPostgres()
	config.ConnectMongo()
	config.ConnectRedis()

	config.DB.AutoMigrate(&models.Student{}, &models.Admin{})

	// Build Provider
	studentProvider := BuildStudentProvider()

	// Inject into controller
	controller.InitStudentProvider(studentProvider)

	routes.SetupRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	router.Run(":8080")
}