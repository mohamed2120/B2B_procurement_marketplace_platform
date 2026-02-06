package main

import (
	"fmt"
	"log"
	"time"

	"github.com/b2b-platform/company-service/models"
	"github.com/b2b-platform/company-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	companyRepo := repository.NewCompanyRepository(db)

	// Demo tenant
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	company := &models.Company{
		ID:              tenantID,
		Name:            "Demo Company Inc",
		LegalName:       "Demo Company Incorporated",
		TaxID:           "TAX-123456",
		Subdomain:       "demo",
		Status:          "approved",
		VerificationStatus: "verified",
		Address:         "123 Business St",
		City:            "New York",
		State:           "NY",
		Country:         "USA",
		PostalCode:      "10001",
		Phone:           "+1-555-0123",
		Email:           "info@democompany.com",
		Website:         "https://democompany.com",
		Industry:        "Manufacturing",
		CompanyType:     "both",
	}

	approvedAt := time.Now()
	company.ApprovedAt = &approvedAt
	company.ApprovedBy = &userID

	if err := companyRepo.Create(company); err != nil {
		log.Printf("Company may already exist: %v", err)
	} else {
		fmt.Printf("Created company: %s\n", company.Name)
	}

	fmt.Println("Company seeding completed!")
}
