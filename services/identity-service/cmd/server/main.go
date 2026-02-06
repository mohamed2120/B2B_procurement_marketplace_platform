package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/identity-service/handlers"
	"github.com/b2b-platform/identity-service/repository"
	"github.com/b2b-platform/identity-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/database"
	"github.com/b2b-platform/shared/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create schema if not exists
	if err := database.CreateSchema(db, "identity"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, roleRepo)
	jwtService := service.NewJWTService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	// Initialize Redis for RBAC caching (future use)
	_, err = redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	// Setup router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "identity-service"})
	})

	// Public routes
	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
	}

	// Protected routes
	protected := api.Group("/auth")
	protected.Use(auth.AuthMiddleware(auth.NewJWTService()))
	{
		protected.GET("/validate", authHandler.ValidateToken)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	fmt.Printf("Identity service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
