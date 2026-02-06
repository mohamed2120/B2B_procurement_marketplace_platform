package main

import (
	"fmt"
	"log"
	"time"

	"github.com/b2b-platform/catalog-service/models"
	"github.com/b2b-platform/catalog-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	manufacturerRepo := repository.NewManufacturerRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	partRepo := repository.NewPartRepository(db)
	attributeRepo := repository.NewAttributeRepository(db)

	// Create manufacturers
	manufacturers := []models.Manufacturer{
		{Name: "Caterpillar", Code: "CAT", Country: "USA", IsActive: true},
		{Name: "Komatsu", Code: "KOM", Country: "Japan", IsActive: true},
		{Name: "Volvo", Code: "VOL", Country: "Sweden", IsActive: true},
	}

	manufacturerMap := make(map[string]uuid.UUID)
	for _, mfr := range manufacturers {
		if err := manufacturerRepo.Create(&mfr); err != nil {
			log.Printf("Manufacturer may already exist: %v", err)
		} else {
			manufacturerMap[mfr.Code] = mfr.ID
			fmt.Printf("Created manufacturer: %s\n", mfr.Name)
		}
	}

	// Create categories
	categories := []models.Category{
		{Name: "Engine Parts", Code: "ENG", Description: "Engine components", IsActive: true},
		{Name: "Hydraulic Components", Code: "HYD", Description: "Hydraulic system parts", IsActive: true},
		{Name: "Electrical", Code: "ELEC", Description: "Electrical components", IsActive: true},
	}

	categoryMap := make(map[string]uuid.UUID)
	for _, cat := range categories {
		if err := categoryRepo.Create(&cat); err != nil {
			log.Printf("Category may already exist: %v", err)
		} else {
			categoryMap[cat.Code] = cat.ID
			fmt.Printf("Created category: %s\n", cat.Name)
		}
	}

	// Create attributes
	attributes := []models.Attribute{
		{Name: "Weight", Code: "WEIGHT", DataType: "number", Unit: "kg", IsSearchable: true},
		{Name: "Material", Code: "MATERIAL", DataType: "string", IsSearchable: true},
		{Name: "Voltage", Code: "VOLTAGE", DataType: "number", Unit: "V", IsSearchable: true},
	}

	attributeMap := make(map[string]uuid.UUID)
	for _, attr := range attributes {
		if err := attributeRepo.Create(&attr); err != nil {
			log.Printf("Attribute may already exist: %v", err)
		} else {
			attributeMap[attr.Code] = attr.ID
			fmt.Printf("Created attribute: %s\n", attr.Name)
		}
	}

	// Create parts
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	parts := []models.LibraryPart{
		{
			PartNumber:     "CAT-ENG-001",
			ManufacturerID: manufacturerMap["CAT"],
			CategoryID:     categoryMap["ENG"],
			Name:           "Engine Oil Filter",
			Description:    "High-performance engine oil filter for Caterpillar engines",
			Status:         "approved",
			CreatedBy:      userID,
		},
		{
			PartNumber:     "KOM-HYD-001",
			ManufacturerID: manufacturerMap["KOM"],
			CategoryID:     categoryMap["HYD"],
			Name:           "Hydraulic Pump Seal",
			Description:    "Replacement seal for Komatsu hydraulic pumps",
			Status:         "pending",
			CreatedBy:      userID,
		},
		{
			PartNumber:     "VOL-ELEC-001",
			ManufacturerID: manufacturerMap["VOL"],
			CategoryID:     categoryMap["ELEC"],
			Name:           "Alternator",
			Description:    "12V alternator for Volvo equipment",
			Status:         "approved",
			CreatedBy:      userID,
		},
	}

	for _, part := range parts {
		if err := partRepo.Create(&part); err != nil {
			log.Printf("Part may already exist: %v", err)
		} else {
			fmt.Printf("Created part: %s - %s\n", part.PartNumber, part.Name)

			// Add attributes for approved parts
			if part.Status == "approved" {
				now := time.Now()
				part.ApprovedAt = &now
				part.ApprovedBy = &userID
				partRepo.Update(&part)

				// Add part attributes
				if part.CategoryID == categoryMap["ELEC"] {
					partAttr := models.PartAttribute{
						PartID:      part.ID,
						AttributeID: attributeMap["VOLTAGE"],
						Value:       "12",
					}
					attributeRepo.AddPartAttribute(&partAttr)
				}
			}
		}
	}

	fmt.Println("Catalog seeding completed!")
}
