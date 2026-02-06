package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/b2b-platform/notification-service/handlers"
	"github.com/b2b-platform/notification-service/repository"
	"github.com/b2b-platform/notification-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/diagnostics"
	"github.com/b2b-platform/shared/health"
	"github.com/b2b-platform/shared/middleware"
	"github.com/b2b-platform/shared/observability"
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

	if err := database.CreateSchema(db, "notification"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)

	templateRepo := repository.NewTemplateRepository(db)
	preferenceRepo := repository.NewPreferenceRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	notificationService := service.NewNotificationService(
		templateRepo,
		preferenceRepo,
		notificationRepo,
	)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Start event consumer in background
	eventConsumer := service.NewEventConsumer(notificationService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := eventConsumer.StartEventConsumer(ctx, eventBus); err != nil {
			if err != context.Canceled {
				log.Printf("Event consumer error: %v", err)
			}
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down notification service...")
		cancel()
	}()

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize logger
	logger := observability.NewLogger("notification-service")

	// Initialize diagnostics reporter
	diagnosticsReporter := diagnostics.NewReporter(db)

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Add error handler middleware
	r.Use(middleware.ErrorHandler(diagnosticsReporter, "notification-service"))

	// Health endpoints
	var redisClient *redis.Client
	if redisClient == nil {
		redisClient, _ = redis.GetRedisClient()
	}
	healthChecker := health.NewHealthChecker("notification-service", db, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		// Template endpoints
		api.GET("/templates", notificationHandler.ListTemplates)
		api.GET("/templates/:id", notificationHandler.GetTemplate)
		api.POST("/templates", notificationHandler.CreateTemplate)

		// Preference endpoints
		api.GET("/preferences", notificationHandler.GetPreferences)
		api.PUT("/preferences", notificationHandler.UpdatePreference)

		// Notification endpoints
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.POST("/notifications", notificationHandler.SendNotification)
		api.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		api.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8009"
	}

	fmt.Printf("Notification service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
