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
	poID1 := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	poID2 := uuid.MustParse("00000000-0000-0000-0000-000000000011")

	// Create shipments (2 total as required)
	shipments := []models.Shipment{
		{
			TenantID:      tenantID,
			POID:          poID1,
			TrackingNumber: "TRACK-001",
			Status:        "in_transit",
			Carrier:       "FedEx",
			ETA:           time.Now().Add(5 * 24 * time.Hour),
			Origin:        "New York, NY",
			Destination:   "Los Angeles, CA",
			IsLate:        false,
		},
		{
			TenantID:      tenantID,
			POID:          poID2,
			TrackingNumber: "TRACK-002",
			Status:        "pending",
			Carrier:       "UPS",
			ETA:           time.Now().Add(7 * 24 * time.Hour),
			Origin:        "Chicago, IL",
			Destination:   "Houston, TX",
			IsLate:        false,
		},
	}

	for _, shipment := range shipments {
		if err := shipmentRepo.Create(&shipment); err != nil {
			log.Printf("Failed to create shipment: %v", err)
		} else {
			fmt.Printf("Created shipment: %s\n", shipment.TrackingNumber)
		}
	}

	fmt.Println("Logistics seeding completed!")
}
