package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/virtual-warehouse-service/handlers"
	"github.com/b2b-platform/virtual-warehouse-service/repository"
	"github.com/b2b-platform/virtual-warehouse-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/diagnostics"
	"github.com/b2b-platform/shared/health"
	"github.com/b2b-platform/shared/middleware"
	"github.com/b2b-platform/shared/observability"
	"github.com/b2b-platform/shared/database"
	"github.com/b2b-platform/shared/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.CreateSchema(db, "virtual_warehouse"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	inventoryRepo := repository.NewInventoryRepository(db)
	groupRepo := repository.NewEquipmentGroupRepository(db)
	transferRepo := repository.NewTransferRepository(db)
	emergencyRepo := repository.NewEmergencyRepository(db)

	warehouseService := service.NewWarehouseService(
		inventoryRepo,
		groupRepo,
		transferRepo,
		emergencyRepo,
	)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseService)

	r := gin.Default()

	
	// Health endpoints
	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}
	if redisClient == nil {
		redisClient, err = redis.GetRedisClient()
		if err != nil {
			log.Printf("Warning: still cannot connect Redis for readiness: %v", err)
		}
	}
	healthChecker := health.NewHealthChecker("virtual-warehouse-service", db, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3002", "http://127.0.0.1:3000", "http://127.0.0.1:3002"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize logger
	logger := observability.NewLogger("virtual-warehouse-service")

	// Initialize diagnostics reporter
	diagnosticsReporter := diagnostics.NewReporter(db)

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Add error handler middleware
	r.Use(middleware.ErrorHandler(diagnosticsReporter, "virtual-warehouse-service"))

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		// Inventory endpoints
		api.GET("/inventory", warehouseHandler.ListInventory)
		api.GET("/inventory/available", warehouseHandler.GetAvailable)
		api.POST("/inventory", warehouseHandler.CreateInventory)

		// Equipment group endpoints
		api.GET("/equipment-groups", warehouseHandler.ListGroups)
		api.GET("/equipment-groups/:id", warehouseHandler.GetGroup)
		api.POST("/equipment-groups", warehouseHandler.CreateGroup)

		// Transfer endpoints
		api.GET("/transfers", warehouseHandler.ListTransfers)
		api.POST("/transfers", warehouseHandler.CreateTransfer)
		api.POST("/transfers/:id/approve", warehouseHandler.ApproveTransfer)

		// Emergency sourcing endpoints
		api.GET("/emergency-sourcing", warehouseHandler.ListEmergencySourcing)
		api.POST("/emergency-sourcing", warehouseHandler.CreateEmergencySourcing)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8011"
	}

	fmt.Printf("Virtual warehouse service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
