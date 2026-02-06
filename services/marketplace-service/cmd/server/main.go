package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/marketplace-service/handlers"
	"github.com/b2b-platform/marketplace-service/repository"
	"github.com/b2b-platform/marketplace-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/diagnostics"
	"github.com/b2b-platform/shared/health"
	"github.com/b2b-platform/shared/middleware"
	"github.com/b2b-platform/shared/observability"
	"github.com/b2b-platform/shared/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.CreateSchema(db, "marketplace"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	storeRepo := repository.NewStoreRepository(db)
	listingRepo := repository.NewListingRepository(db)
	mediaRepo := repository.NewMediaRepository(db)

	marketplaceService := service.NewMarketplaceService(storeRepo, listingRepo, mediaRepo)
	marketplaceHandler := handlers.NewMarketplaceHandler(marketplaceService)

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize logger
	logger := observability.NewLogger("marketplace-service")

	// Initialize diagnostics reporter
	diagnosticsReporter := diagnostics.NewReporter(db)

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Add error handler middleware
	r.Use(middleware.ErrorHandler(diagnosticsReporter, "marketplace-service"))

	// Health endpoints
	var redisClient *redis.Client
	if redisClient == nil {
		redisClient, _ = redis.GetRedisClient()
	}
	healthChecker := health.NewHealthChecker("marketplace-service", db, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		// Store endpoints
		api.GET("/stores", marketplaceHandler.ListStores)
		api.GET("/stores/:id", marketplaceHandler.GetStore)
		api.POST("/stores", marketplaceHandler.CreateStore)
		api.PUT("/stores/:id", marketplaceHandler.UpdateStore)
		api.GET("/stores/:id/listings", marketplaceHandler.GetStoreListings)

		// Listing endpoints
		api.GET("/listings", marketplaceHandler.ListListings)
		api.GET("/listings/:id", marketplaceHandler.GetListing)
		api.POST("/listings", marketplaceHandler.CreateListing)
		api.PUT("/listings/:id", marketplaceHandler.UpdateListing)
		api.PUT("/listings/:id/stock", marketplaceHandler.UpdateStock)

		// Media endpoints
		api.GET("/listings/:id/media", marketplaceHandler.GetListingMedia)
		api.POST("/listings/:id/media", marketplaceHandler.AddMedia)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	fmt.Printf("Marketplace service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
