package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/logistics-service/handlers"
	"github.com/b2b-platform/logistics-service/repository"
	"github.com/b2b-platform/logistics-service/service"
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

	if err := database.CreateSchema(db, "logistics"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)

	shipmentRepo := repository.NewShipmentRepository(db)
	logisticsService := service.NewLogisticsService(shipmentRepo, eventBus)
	logisticsHandler := handlers.NewLogisticsHandler(logisticsService)

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "logistics-service"})
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		api.GET("/shipments", logisticsHandler.List)
		api.GET("/shipments/:id", logisticsHandler.Get)
		api.POST("/shipments", logisticsHandler.Create)
		api.PUT("/shipments/:id/tracking", logisticsHandler.UpdateTracking)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8007"
	}

	fmt.Printf("Logistics service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
