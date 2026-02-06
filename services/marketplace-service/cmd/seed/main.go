package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/marketplace-service/models"
	"github.com/b2b-platform/marketplace-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	storeRepo := repository.NewStoreRepository(db)
	listingRepo := repository.NewListingRepository(db)
	mediaRepo := repository.NewMediaRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	companyID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Create store
	store := &models.Store{
		TenantID:    tenantID,
		CompanyID:   companyID,
		Name:        "Demo Parts Store",
		Description: "Quality spare parts for heavy equipment",
		Status:      "active",
		IsVerified:  true,
	}

	if err := storeRepo.Create(store); err != nil {
		log.Printf("Store may already exist: %v", err)
	} else {
		fmt.Printf("Created store: %s\n", store.Name)

		// Create listings
		listings := []models.Listing{
			{
				TenantID:        tenantID,
				StoreID:         store.ID,
				ListingType:     "product",
				Title:           "Engine Oil Filter - CAT",
				Description:     "High-quality engine oil filter for Caterpillar equipment",
				SKU:             "CAT-FILTER-001",
				Status:          "active",
				Price:           25.00,
				Currency:        "USD",
				StockQuantity:   100,
				MinOrderQuantity: 1,
				LeadTimeDays:    7,
				IsActive:        true,
				CreatedBy:       userID,
			},
			{
				TenantID:        tenantID,
				StoreID:         store.ID,
				ListingType:     "service",
				Title:           "Equipment Maintenance Service",
				Description:     "Professional maintenance service for heavy equipment",
				SKU:             "SVC-MAINT-001",
				Status:          "active",
				Price:           500.00,
				Currency:        "USD",
				StockQuantity:   0,
				MinOrderQuantity: 1,
				LeadTimeDays:    14,
				IsActive:        true,
				CreatedBy:       userID,
			},
			{
				TenantID:        tenantID,
				StoreID:         store.ID,
				ListingType:     "surplus",
				Title:           "Surplus Hydraulic Pump",
				Description:     "Surplus hydraulic pump in good condition",
				SKU:             "SURP-PUMP-001",
				Status:          "active",
				Price:           1200.00,
				Currency:        "USD",
				StockQuantity:   5,
				MinOrderQuantity: 1,
				LeadTimeDays:    3,
				IsActive:        true,
				CreatedBy:       userID,
			},
		}

		for _, listing := range listings {
			if err := listingRepo.Create(&listing); err != nil {
				log.Printf("Failed to create listing: %v", err)
			} else {
				fmt.Printf("Created listing: %s\n", listing.Title)

				// Add media for first listing
				if listing.ListingType == "product" {
					media := &models.ListingMedia{
						ListingID:   listing.ID,
						MediaType:   "image",
						URL:         "https://example.com/images/cat-filter.jpg",
						ThumbnailURL: "https://example.com/images/cat-filter-thumb.jpg",
						FileName:    "cat-filter.jpg",
						IsPrimary:   true,
						SortOrder:   0,
					}

					if err := mediaRepo.Create(media); err != nil {
						log.Printf("Failed to create media: %v", err)
					} else {
						fmt.Printf("Created media for listing: %s\n", listing.Title)
					}
				}
			}
		}
	}

	fmt.Println("Marketplace seeding completed!")
}
