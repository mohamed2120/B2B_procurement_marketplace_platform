package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/diagnostics-service/handlers"
	"github.com/b2b-platform/diagnostics-service/repository"
	"github.com/b2b-platform/shared/auth"
	"github.com/b2b-platform/shared/database"
	"github.com/b2b-platform/shared/health"
	"github.com/b2b-platform/shared/observability"
	"github.com/b2b-platform/shared/middleware"
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
	if err := database.CreateSchema(db, "diagnostics"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Initialize logger
	logger := observability.NewLogger("diagnostics-service")

	// Initialize repositories
	diagnosticsRepo := repository.NewDiagnosticsRepository(db)

	// Initialize handlers
	diagnosticsHandler := handlers.NewDiagnosticsHandler(diagnosticsRepo)

	// Setup router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Health endpoints
	healthChecker := health.NewHealthChecker("diagnostics-service", db, nil)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)

	// API routes (admin only)
	api := r.Group("/api/diagnostics/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	api.Use(auth.RoleMiddleware("admin")) // Only admins can access diagnostics
	{
		api.GET("/summary", diagnosticsHandler.GetSummary)
		api.GET("/services", diagnosticsHandler.GetServices)
		api.GET("/incidents", diagnosticsHandler.GetIncidents)
		api.GET("/incidents/:id", diagnosticsHandler.GetIncident)
		api.POST("/incidents/:id/resolve", diagnosticsHandler.ResolveIncident)
		api.GET("/events/failures", diagnosticsHandler.GetEventFailures)
		api.POST("/events/failures/:id/retry", diagnosticsHandler.RetryEventFailure)
		api.GET("/metrics", diagnosticsHandler.GetMetrics)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8013"
	}

	fmt.Printf("Diagnostics service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
