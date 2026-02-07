package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/b2b-platform/search-indexer-service/handlers"
	"github.com/b2b-platform/search-indexer-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/events"
	"github.com/b2b-platform/shared/health"
	"github.com/b2b-platform/shared/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)
	indexerService := service.NewIndexerService()
	searchService := service.NewSearchService()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3002", "http://127.0.0.1:3000", "http://127.0.0.1:3002"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Tenant-ID", "X-Request-Id"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health endpoints
	healthChecker := health.NewHealthChecker("search-indexer-service", nil, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)

	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(searchService)

	// Search API routes (public, but JWT optional for enhanced access)
	api := r.Group("/api/v1")
	api.Use(auth.OptionalAuthMiddleware(auth.NewJWTService()))
	{
		// Public search endpoint (works without auth, but enhanced with auth)
		api.GET("/search", searchHandler.Search)
		api.GET("/search/autocomplete", searchHandler.Autocomplete)
	}

	// Protected routes (optional - for admin operations like reindex)
	protected := api.Group("/admin")
	protected.Use(auth.OptionalAuthMiddleware(auth.NewJWTService()))
	{
		// Future: reindex endpoint, stats, etc.
	}

	// Start HTTP server in goroutine
	port := os.Getenv("PORT")
	if port == "" {
		port = "8012"
	}

	go func() {
		fmt.Printf("Search indexer service starting on port %s\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down search indexer...")
		cancel()
	}()

	// Subscribe to all events
	fmt.Println("Subscribing to events...")

	handler := func(event *events.EventEnvelope) error {
		fmt.Printf("Received event: %s\n", event.Type)
		if err := indexerService.HandleEvent(event); err != nil {
			log.Printf("Error handling event %s: %v", event.Type, err)
			return err
		}
		fmt.Printf("Successfully indexed event: %s\n", event.Type)
		return nil
	}

	// Subscribe to all events
	if err := eventBus.SubscribeAll(ctx, handler); err != nil {
		if err == context.Canceled {
			fmt.Println("Search indexer stopped")
			return
		}
		log.Fatalf("Error in event subscription: %v", err)
	}
}
