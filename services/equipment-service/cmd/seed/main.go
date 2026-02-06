package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/equipment-service/models"
	"github.com/b2b-platform/equipment-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	equipmentRepo := repository.NewEquipmentRepository(db)
	bomRepo := repository.NewBOMRepository(db)
	compatibilityRepo := repository.NewCompatibilityRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Create equipment
	equipment := &models.Equipment{
		TenantID:        tenantID,
		EquipmentNumber: "EQ-001",
		Name:            "Excavator CAT 320",
		Type:            "excavator",
		Manufacturer:    "Caterpillar",
		Model:           "320",
		SerialNumber:    "CAT320-12345",
		Year:            2020,
		Status:          "active",
		Location:        "Construction Site A",
	}

	if err := equipmentRepo.Create(equipment); err != nil {
		log.Printf("Equipment may already exist: %v", err)
	} else {
		fmt.Printf("Created equipment: %s\n", equipment.Name)

		// Add BOM nodes
		bomNodes := []models.BOMNode{
			{
				TenantID:    tenantID,
				EquipmentID: equipment.ID,
				PartName:    "Engine Oil Filter",
				PartNumber:  "CAT-ENG-001",
				Description: "Primary engine oil filter",
				Quantity:    1,
				Unit:        "piece",
				Position:    "Engine compartment",
				Level:       0,
			},
			{
				TenantID:    tenantID,
				EquipmentID: equipment.ID,
				PartName:    "Hydraulic Fluid",
				PartNumber:  "HYD-FLUID-001",
				Description: "Hydraulic system fluid",
				Quantity:    50,
				Unit:        "liter",
				Position:    "Hydraulic reservoir",
				Level:       0,
			},
		}

		for _, node := range bomNodes {
			if err := bomRepo.Create(&node); err != nil {
				log.Printf("Failed to create BOM node: %v", err)
			} else {
				fmt.Printf("Created BOM node: %s\n", node.PartName)
			}
		}

		// Add compatibility mapping (assuming part exists in catalog)
		partID := uuid.MustParse("00000000-0000-0000-0000-000000000100") // Example part ID
		mapping := &models.CompatibilityMapping{
			TenantID:     tenantID,
			EquipmentID:  equipment.ID,
			PartID:       partID,
			IsCompatible: true,
			Notes:        "Verified compatible",
		}

		if err := compatibilityRepo.Create(mapping); err != nil {
			log.Printf("Failed to create compatibility mapping: %v", err)
		} else {
			fmt.Printf("Created compatibility mapping for part: %s\n", partID.String())
		}
	}

	fmt.Println("Equipment seeding completed!")
}
