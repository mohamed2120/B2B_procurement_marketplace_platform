package main

import (
	"log"
	"os"
	"time"

	"github.com/b2b-platform/diagnostics-service/models"
	"github.com/b2b-platform/shared/database"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.CreateSchema(db, "diagnostics"); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Seed service heartbeats
	services := []string{
		"identity-service", "company-service", "catalog-service",
		"equipment-service", "marketplace-service", "procurement-service",
		"logistics-service", "collaboration-service", "notification-service",
		"billing-service", "virtual-warehouse-service", "search-indexer-service",
	}

	for _, serviceName := range services {
		heartbeat := models.ServiceHeartbeat{
			ServiceName: serviceName,
			InstanceID:  "instance-1",
			LastSeenAt:  time.Now(),
			Status:      "healthy",
			Version:     "1.0.0",
			Env:         env,
		}
		if err := db.Create(&heartbeat).Error; err != nil {
			log.Printf("Warning: Failed to create heartbeat for %s: %v", serviceName, err)
		}
	}

	// Seed sample incidents
	incidents := []models.Incident{
		{
			Severity:    "ERROR",
			ServiceName: "identity-service",
			Category:    "AUTH",
			ErrorCode:   stringPtr("AUTH_001"),
			Title:       "Invalid token received",
			DetailsJSON: models.JSONB{"ip": "192.168.1.1", "user_agent": "Mozilla/5.0"},
			RequestID:   stringPtr("req-123"),
			OccurredAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			Severity:    "WARN",
			ServiceName: "procurement-service",
			Category:    "API",
			ErrorCode:   stringPtr("API_001"),
			Title:       "Bad request: missing required field",
			DetailsJSON: models.JSONB{"field": "quantity", "route": "/api/v1/purchase-requests"},
			RequestID:   stringPtr("req-456"),
			OccurredAt:  time.Now().Add(-1 * time.Hour),
		},
		{
			Severity:    "CRITICAL",
			ServiceName: "billing-service",
			Category:    "DB",
			ErrorCode:   stringPtr("DB_001"),
			Title:       "Database connection timeout",
			DetailsJSON: models.JSONB{"timeout_seconds": 30},
			OccurredAt:  time.Now().Add(-30 * time.Minute),
		},
	}

	for _, incident := range incidents {
		if err := db.Create(&incident).Error; err != nil {
			log.Printf("Warning: Failed to create incident: %v", err)
		}
	}

	// Seed event failures
	eventFailures := []models.EventFailure{
		{
			EventName:     "company.created",
			Direction:     "PUBLISH",
			ServiceName:   "company-service",
			PayloadJSON:   models.JSONB{"company_id": "123", "name": "Acme Corp"},
			ErrorMessage:  stringPtr("Redis connection failed"),
			RetryCount:    2,
			LastAttemptAt: time.Now().Add(-1 * time.Hour),
			Status:        "FAILED",
		},
		{
			EventName:     "order.created",
			Direction:     "CONSUME",
			ServiceName:   "notification-service",
			PayloadJSON:   models.JSONB{"order_id": "456"},
			ErrorMessage:  stringPtr("Failed to deserialize event"),
			RetryCount:    1,
			LastAttemptAt: time.Now().Add(-30 * time.Minute),
			Status:        "RETRYING",
		},
	}

	for _, failure := range eventFailures {
		if err := db.Create(&failure).Error; err != nil {
			log.Printf("Warning: Failed to create event failure: %v", err)
		}
	}

	log.Println("Diagnostics seed data created successfully")
}

func stringPtr(s string) *string {
	return &s
}
