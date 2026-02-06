package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/equipment-service/handlers"
	"github.com/b2b-platform/equipment-service/repository"
	"github.com/b2b-platform/equipment-service/service"
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

	if err := database.CreateSchema(db, "equipment"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	equipmentRepo := repository.NewEquipmentRepository(db)
	bomRepo := repository.NewBOMRepository(db)
	compatibilityRepo := repository.NewCompatibilityRepository(db)

	equipmentService := service.NewEquipmentService(equipmentRepo, bomRepo, compatibilityRepo)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService)

	r := gin.Default()

	
	// Health endpoints
	var redisClient *redis.Client
	if redisClient == nil {
		redisClient, _ = redis.GetRedisClient()
	}
	healthChecker := health.NewHealthChecker("equipment-service", db, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize logger
	logger := observability.NewLogger("equipment-service")

	// Initialize diagnostics reporter
	diagnosticsReporter := diagnostics.NewReporter(db)

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Add error handler middleware
	r.Use(middleware.ErrorHandler(diagnosticsReporter, "equipment-service"))

		})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.TenantMiddleware())
	{
		api.GET("/equipment", equipmentHandler.List)
		api.GET("/equipment/:id", equipmentHandler.Get)
		api.POST("/equipment", equipmentHandler.Create)
		api.PUT("/equipment/:id", equipmentHandler.Update)

		// BOM endpoints
		api.GET("/equipment/:id/bom", equipmentHandler.GetBOM)
		api.POST("/equipment/:id/bom", equipmentHandler.AddBOMNode)

		// Compatibility endpoints
		api.GET("/equipment/:id/compatibility", equipmentHandler.GetCompatibilityMappings)
		api.POST("/equipment/:id/compatibility", equipmentHandler.CreateCompatibilityMapping)
		api.GET("/equipment/:id/compatibility/check", equipmentHandler.CheckCompatibility)
		api.POST("/compatibility/:mapping_id/verify", equipmentHandler.VerifyCompatibility)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	fmt.Printf("Equipment service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
