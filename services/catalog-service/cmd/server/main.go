package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/catalog-service/handlers"
	"github.com/b2b-platform/catalog-service/repository"
	"github.com/b2b-platform/catalog-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/database"
	"github.com/b2b-platform/shared/events"
	"github.com/b2b-platform/shared/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.CreateSchema(db, "catalog"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)

	manufacturerRepo := repository.NewManufacturerRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	partRepo := repository.NewPartRepository(db)
	attributeRepo := repository.NewAttributeRepository(db)

	catalogService := service.NewCatalogService(
		manufacturerRepo,
		categoryRepo,
		partRepo,
		attributeRepo,
		eventBus,
	)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "catalog-service"})
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	{
		// Manufacturers
		api.GET("/manufacturers", catalogHandler.ListManufacturers)
		api.GET("/manufacturers/:id", catalogHandler.GetManufacturer)
		api.POST("/manufacturers", catalogHandler.CreateManufacturer)

		// Categories
		api.GET("/categories", catalogHandler.ListCategories)
		api.GET("/categories/:id", catalogHandler.GetCategory)
		api.POST("/categories", catalogHandler.CreateCategory)

		// Attributes
		api.GET("/attributes", catalogHandler.ListAttributes)
		api.POST("/attributes", catalogHandler.CreateAttribute)

		// Parts
		api.GET("/parts", catalogHandler.ListParts)
		api.GET("/parts/:id", catalogHandler.GetPart)
		api.POST("/parts", catalogHandler.CreatePart)
		api.GET("/parts/pending", catalogHandler.GetPendingParts)
		api.POST("/parts/:id/approve", catalogHandler.ApprovePart)
		api.POST("/parts/:id/reject", catalogHandler.RejectPart)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	fmt.Printf("Catalog service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
