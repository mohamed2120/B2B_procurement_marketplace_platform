package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/collaboration-service/handlers"
	"github.com/b2b-platform/collaboration-service/repository"
	"github.com/b2b-platform/collaboration-service/service"
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

	if err := database.CreateSchema(db, "collaboration"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)

	threadRepo := repository.NewThreadRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	disputeRepo := repository.NewDisputeRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	collaborationService := service.NewCollaborationService(
		threadRepo,
		messageRepo,
		disputeRepo,
		ratingRepo,
		eventBus,
	)
	collaborationHandler := handlers.NewCollaborationHandler(collaborationService)

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "collaboration-service"})
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		// Thread endpoints
		api.GET("/threads", collaborationHandler.ListThreads)
		api.GET("/threads/user", collaborationHandler.GetUserThreads)
		api.POST("/threads", collaborationHandler.CreateThread)
		
		// Message endpoints (must come before /threads/:id to avoid route conflict)
		api.GET("/threads/:id/messages", collaborationHandler.GetThreadMessages)
		api.POST("/threads/:id/messages", collaborationHandler.SendMessage)
		
		// Single thread endpoint (must come after messages routes)
		api.GET("/threads/:id", collaborationHandler.GetThread)

		// Dispute endpoints
		api.GET("/disputes", collaborationHandler.ListDisputes)
		api.GET("/disputes/:id", collaborationHandler.GetDispute)
		api.POST("/disputes", collaborationHandler.CreateDispute)
		api.POST("/disputes/:id/resolve", collaborationHandler.ResolveDispute)

		// Rating endpoints
		api.GET("/ratings", collaborationHandler.GetRatings)
		api.GET("/ratings/average", collaborationHandler.GetAverageRating)
		api.POST("/ratings", collaborationHandler.CreateRating)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	fmt.Printf("Collaboration service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
