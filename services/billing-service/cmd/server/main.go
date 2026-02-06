package main

import (
	"fmt"
	"log"
	"os"

	"github.com/b2b-platform/billing-service/handlers"
	"github.com/b2b-platform/billing-service/repository"
	"github.com/b2b-platform/billing-service/service"
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

	if err := database.CreateSchema(db, "billing"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)

	planRepo := repository.NewPlanRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	escrowRepo := repository.NewEscrowRepository(db)
	settlementRepo := repository.NewSettlementRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	payoutRepo := repository.NewPayoutRepository(db)

	// Payment provider (mock for local dev)
	paymentProvider := service.NewMockPaymentProvider()

	billingService := service.NewBillingService(planRepo, subscriptionRepo, eventBus)
	paymentService := service.NewPaymentService(paymentRepo, escrowRepo, settlementRepo, refundRepo, payoutRepo, paymentProvider, eventBus)
	payoutService := service.NewPayoutService(payoutRepo)

	billingHandler := handlers.NewBillingHandler(billingService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, payoutService)

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize logger
	logger := observability.NewLogger("billing-service")

	// Initialize diagnostics reporter
	diagnosticsReporter := diagnostics.NewReporter(db)

	// Add logging middleware
	r.Use(middleware.RequestLogging(logger))

	// Add error handler middleware
	r.Use(middleware.ErrorHandler(diagnosticsReporter, "billing-service"))

	// Health endpoints
	var redisClient *redis.Client
	if redisClient == nil {
		redisClient, _ = redis.GetRedisClient()
	}
	healthChecker := health.NewHealthChecker("billing-service", db, redisClient)
	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Ready)
	})

	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(auth.NewJWTService()))
	{
		// Plan endpoints (public for listing)
		api.GET("/plans", billingHandler.ListPlans)
		api.GET("/plans/:id", billingHandler.GetPlan)
		api.POST("/plans", billingHandler.CreatePlan)

		// Subscription endpoints
		api.Use(auth.TenantMiddleware())
		api.GET("/subscriptions", billingHandler.GetSubscription)
		api.POST("/subscriptions", billingHandler.CreateSubscription)
		api.DELETE("/subscriptions/:id", billingHandler.CancelSubscription)
		api.GET("/entitlements/check", billingHandler.CheckEntitlement)

		// Payment endpoints
		api.POST("/billing/v1/payments/intent", paymentHandler.CreatePaymentIntent)
		api.GET("/billing/v1/payments/:id", paymentHandler.GetPayment)
		api.GET("/billing/v1/payments", paymentHandler.ListPayments)

		// Escrow endpoints
		api.POST("/billing/v1/escrow/release", paymentHandler.ReleaseEscrow)
		api.GET("/billing/v1/escrow/:id", paymentHandler.GetEscrowHold)
		api.GET("/billing/v1/escrow", paymentHandler.ListEscrowHolds)

		// Refund endpoints
		api.POST("/billing/v1/refunds", paymentHandler.CreateRefund)

		// Payout account endpoints
		api.POST("/billing/v1/payout-accounts", paymentHandler.CreatePayoutAccount)
		api.GET("/billing/v1/payout-accounts", paymentHandler.ListPayoutAccounts)
		api.GET("/billing/v1/payout-accounts/:id", paymentHandler.GetPayoutAccount)
		api.PUT("/billing/v1/payout-accounts/:id", paymentHandler.UpdatePayoutAccount)
		api.DELETE("/billing/v1/payout-accounts/:id", paymentHandler.DeletePayoutAccount)
		api.PUT("/billing/v1/payout-accounts/:id/default", paymentHandler.SetDefaultPayoutAccount)
	}

	// Webhook endpoint (no auth required - called by payment provider)
	r.POST("/api/billing/v1/payments/webhook", paymentHandler.HandleWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}

	fmt.Printf("Billing service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
