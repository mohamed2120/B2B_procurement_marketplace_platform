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

	// Create equipment (10 total as required)
	equipmentList := []models.Equipment{
		{TenantID: tenantID, EquipmentNumber: "EQ-001", Name: "Excavator CAT 320", Type: "excavator", Manufacturer: "Caterpillar", Model: "320", SerialNumber: "CAT320-12345", Year: 2020, Status: "active", Location: "Construction Site A"},
		{TenantID: tenantID, EquipmentNumber: "EQ-002", Name: "Bulldozer CAT D6", Type: "bulldozer", Manufacturer: "Caterpillar", Model: "D6", SerialNumber: "CATD6-67890", Year: 2019, Status: "active", Location: "Construction Site B"},
		{TenantID: tenantID, EquipmentNumber: "EQ-003", Name: "Loader KOM WA380", Type: "loader", Manufacturer: "Komatsu", Model: "WA380", SerialNumber: "KOM380-11111", Year: 2021, Status: "active", Location: "Warehouse"},
		{TenantID: tenantID, EquipmentNumber: "EQ-004", Name: "Excavator KOM PC200", Type: "excavator", Manufacturer: "Komatsu", Model: "PC200", SerialNumber: "KOMPC200-22222", Year: 2020, Status: "active", Location: "Construction Site A"},
		{TenantID: tenantID, EquipmentNumber: "EQ-005", Name: "Dump Truck VOL A25", Type: "truck", Manufacturer: "Volvo", Model: "A25", SerialNumber: "VOLA25-33333", Year: 2022, Status: "active", Location: "Mining Site"},
		{TenantID: tenantID, EquipmentNumber: "EQ-006", Name: "Crane CAT 350", Type: "crane", Manufacturer: "Caterpillar", Model: "350", SerialNumber: "CAT350-44444", Year: 2018, Status: "active", Location: "Port"},
		{TenantID: tenantID, EquipmentNumber: "EQ-007", Name: "Compactor KOM BW", Type: "compactor", Manufacturer: "Komatsu", Model: "BW", SerialNumber: "KOMBW-55555", Year: 2021, Status: "active", Location: "Road Construction"},
		{TenantID: tenantID, EquipmentNumber: "EQ-008", Name: "Grader CAT 140", Type: "grader", Manufacturer: "Caterpillar", Model: "140", SerialNumber: "CAT140-66666", Year: 2020, Status: "active", Location: "Highway Project"},
		{TenantID: tenantID, EquipmentNumber: "EQ-009", Name: "Excavator VOL EC", Type: "excavator", Manufacturer: "Volvo", Model: "EC", SerialNumber: "VOLEC-77777", Year: 2023, Status: "active", Location: "Urban Development"},
		{TenantID: tenantID, EquipmentNumber: "EQ-010", Name: "Loader CAT 950", Type: "loader", Manufacturer: "Caterpillar", Model: "950", SerialNumber: "CAT950-88888", Year: 2019, Status: "active", Location: "Quarry"},
	}

	for _, equipment := range equipmentList {
		if err := equipmentRepo.Create(&equipment); err != nil {
			log.Printf("Equipment may already exist: %v", err)
			continue
		}
		fmt.Printf("Created equipment: %s\n", equipment.Name)

		// Add BOM nodes for first equipment only
		if equipment.EquipmentNumber == "EQ-001" {
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

			// Add compatibility mapping
			partID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
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
	}

	fmt.Println("Equipment seeding completed!")
}
