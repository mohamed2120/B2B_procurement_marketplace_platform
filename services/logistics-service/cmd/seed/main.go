package main

import (
	"fmt"
	"log"
	"time"

	"github.com/b2b-platform/logistics-service/models"
	"github.com/b2b-platform/logistics-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	shipmentRepo := repository.NewShipmentRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	poID := uuid.MustParse("00000000-0000-0000-0000-000000000010")

	shipment := &models.Shipment{
		TenantID:      tenantID,
		POID:          poID,
		TrackingNumber: "TRACK-001",
		Status:        "in_transit",
		Carrier:       "FedEx",
		ETA:           time.Now().Add(5 * 24 * time.Hour),
		Origin:        "New York, NY",
		Destination:   "Los Angeles, CA",
		IsLate:        false,
	}

	if err := shipmentRepo.Create(shipment); err != nil {
		log.Printf("Failed to create shipment: %v", err)
	} else {
		fmt.Printf("Created shipment: %s\n", shipment.TrackingNumber)
	}

	fmt.Println("Logistics seeding completed!")
}
