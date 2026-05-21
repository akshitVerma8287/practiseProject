package main

import (
	"log"
	"os"
	"project/config"
	"project/controller"
	"project/dbops"
	"project/models"
	"project/routes"
	"project/services"
	"strconv"

	"github.com/joho/godotenv"

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

	factory := dbops.NewDatabaseConnectorFactory()

	// REDIS CONFIG
	redisPoolSize, _ := strconv.Atoi(
		os.Getenv("REDIS_POOL_SIZE"),
	)

	redisConfig := dbops.DatabaseConfig{
		Type:     dbops.RedisType,
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		PoolSize: redisPoolSize,
		Database: "RedisMain",
	}

	redisConnector, err := factory.GetConnector(
		redisConfig,
	)

	if err != nil {
		log.Fatalf(
			"Failed to connect Redis: %v",
			err,
		)
	}

	defer redisConnector.Close()

	dbops.RedisRepo = redisConnector.GetRepository().( dbops.RedisRepository )


	// POSTGRES CONFIG

	maxConn, _ := strconv.Atoi(
		os.Getenv("POSTGRES_MAX_CONNECTION"),
	)

	maxIdleConn, _ := strconv.Atoi(
		os.Getenv("POSTGRES_MAX_IDLE_CONNECTIONS"),
	)

	postgresConfig := dbops.DatabaseConfig{
		Type:               dbops.PostgresType,
		Host:               os.Getenv("POSTGRES_HOST"),
		Port:               os.Getenv("POSTGRES_PORT"),
		Username:           os.Getenv("POSTGRES_USERNAME"),
		Password:           os.Getenv("POSTGRES_PASSWORD"),
		Database:           os.Getenv("POSTGRES_DATABASE"),
		MaxConnection:      maxConn,
		MaxIdleConnections: maxIdleConn,
	}

	postgresConnector, err := factory.GetConnector(
		postgresConfig,
	)

	if err != nil {
		log.Fatalf(
			"Failed to connect Postgres: %v",
			err,
		)
	}

	defer postgresConnector.Close()

	dbops.PostgresRepo = postgresConnector.GetRepository().( dbops.PostgresRepository )
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