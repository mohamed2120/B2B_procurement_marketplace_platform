package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/b2b-platform/virtual-warehouse-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	inventoryRepo := repository.NewInventoryRepository(db)
	groupRepo := repository.NewEquipmentGroupRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	partID := uuid.MustParse("00000000-0000-0000-0000-000000000100")

	// Create shared inventory
	inventory := &models.SharedInventory{
		TenantID:    tenantID,
		PartID:      partID,
		Quantity:    50,
		Location:    "Warehouse A",
		IsAvailable: true,
		ReservedQty: 0,
	}

	if err := inventoryRepo.Create(inventory); err != nil {
		log.Printf("Inventory may already exist: %v", err)
	} else {
		fmt.Printf("Created shared inventory: %f units\n", inventory.Quantity)
	}

	// Create equipment group
	group := &models.EquipmentGroup{
		TenantID:    tenantID,
		Name:        "Fleet Group A",
		Description: "Main equipment fleet",
	}

	if err := groupRepo.Create(group); err != nil {
		log.Printf("Group may already exist: %v", err)
	} else {
		fmt.Printf("Created equipment group: %s\n", group.Name)
	}

	fmt.Println("Virtual warehouse seeding completed!")
}
