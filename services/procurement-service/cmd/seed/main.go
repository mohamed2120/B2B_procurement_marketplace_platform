package main

import (
	"fmt"
	"log"
	"time"

	"github.com/b2b-platform/procurement-service/models"
	"github.com/b2b-platform/procurement-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	prRepo := repository.NewPRRepository(db)
	rfqRepo := repository.NewRFQRepository(db)
	quoteRepo := repository.NewQuoteRepository(db)
	poRepo := repository.NewPORepository(db)

	// Demo tenant
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	supplierID := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	// Create PR
	pr := &models.PurchaseRequest{
		TenantID:    tenantID,
		PRNumber:    "PR-001",
		Title:       "Office Supplies Purchase",
		Description: "Monthly office supplies procurement",
		Status:      "approved",
		Priority:    "normal",
		RequestedBy:  userID,
		Department:  "Operations",
		Budget:      5000.00,
		Currency:    "USD",
	}

	requiredDate := time.Now().Add(30 * 24 * time.Hour)
	pr.RequiredDate = &requiredDate

	approvedAt := time.Now()
	pr.ApprovedAt = &approvedAt
	pr.ApprovedBy = &userID

	if err := prRepo.Create(pr); err != nil {
		log.Printf("Failed to create PR: %v", err)
	} else {
		fmt.Printf("Created PR: %s\n", pr.PRNumber)

		// Add PR items
		items := []models.PRItem{
			{
				PRID:         pr.ID,
				Description:  "Printer Paper A4",
				Quantity:     100,
				Unit:         "ream",
				UnitPrice:    25.00,
				TotalPrice:   2500.00,
				Specifications: "80gsm, white",
			},
			{
				PRID:         pr.ID,
				Description:  "Office Pens",
				Quantity:     50,
				Unit:         "box",
				UnitPrice:    15.00,
				TotalPrice:   750.00,
			},
		}

		for _, item := range items {
			if err := db.Create(&item).Error; err != nil {
				log.Printf("Failed to create PR item: %v", err)
			}
		}

		// Create RFQs (2 total as required)
		rfqs := []*models.RFQ{
			{
				TenantID:    tenantID,
				PRID:        pr.ID,
				RFQNumber:   "RFQ-001",
				Title:       "RFQ for Office Supplies",
				Description: "Request for quotation for office supplies",
				Status:      "sent",
				DueDate:     time.Now().Add(7 * 24 * time.Hour),
				CreatedBy:   userID,
			},
			{
				TenantID:    tenantID,
				PRID:        pr.ID,
				RFQNumber:   "RFQ-002",
				Title:       "RFQ for Equipment Parts",
				Description: "Request for quotation for heavy equipment spare parts",
				Status:      "sent",
				DueDate:     time.Now().Add(14 * 24 * time.Hour),
				CreatedBy:   userID,
			},
		}

		for _, rfq := range rfqs {
			if err := rfqRepo.Create(rfq); err != nil {
				log.Printf("Failed to create RFQ: %v", err)
				continue
			}
			fmt.Printf("Created RFQ: %s\n", rfq.RFQNumber)

			// Create Quotes (3 total as required - 2 for RFQ-001, 1 for RFQ-002)
			quotes := []*models.Quote{}
			if rfq.RFQNumber == "RFQ-001" {
				quotes = []*models.Quote{
					{
						TenantID:    tenantID,
						RFQID:       rfq.ID,
						SupplierID:  supplierID,
						QuoteNumber: "QT-001",
						Status:      "submitted",
						TotalAmount: 3250.00,
						Currency:    "USD",
						ValidUntil:  time.Now().Add(30 * 24 * time.Hour),
						Notes:       "Best price guaranteed",
						SubmittedAt: time.Now(),
					},
					{
						TenantID:    tenantID,
						RFQID:       rfq.ID,
						SupplierID:  supplierID,
						QuoteNumber: "QT-002",
						Status:      "submitted",
						TotalAmount: 3100.00,
						Currency:    "USD",
						ValidUntil:  time.Now().Add(30 * 24 * time.Hour),
						Notes:       "Competitive pricing",
						SubmittedAt: time.Now(),
					},
				}
			} else {
				quotes = []*models.Quote{
					{
						TenantID:    tenantID,
						RFQID:       rfq.ID,
						SupplierID:  supplierID,
						QuoteNumber: "QT-003",
						Status:      "submitted",
						TotalAmount: 5500.00,
						Currency:    "USD",
						ValidUntil:  time.Now().Add(30 * 24 * time.Hour),
						Notes:       "Parts quote",
						SubmittedAt: time.Now(),
					},
				}
			}

			for _, quote := range quotes {
				if err := quoteRepo.Create(quote); err != nil {
					log.Printf("Failed to create quote: %v", err)
					continue
				}
				fmt.Printf("Created Quote: %s\n", quote.QuoteNumber)

				// Add quote items
				quoteItems := []models.QuoteItem{
					{
						QuoteID:     quote.ID,
						PRItemID:    items[0].ID,
						Description: "Printer Paper A4",
						Quantity:    100,
						UnitPrice:   24.00,
						TotalPrice:  2400.00,
						LeadTime:    7,
					},
					{
						QuoteID:     quote.ID,
						PRItemID:    items[1].ID,
						Description: "Office Pens",
						Quantity:    50,
						UnitPrice:   14.00,
						TotalPrice:  700.00,
						LeadTime:    5,
					},
				}

				for _, item := range quoteItems {
					if err := db.Create(&item).Error; err != nil {
						log.Printf("Failed to create quote item: %v", err)
					}
				}

				// Create PO only for first quote of first RFQ (2 orders total as required)
				if rfq.RFQNumber == "RFQ-001" && quote.QuoteNumber == "QT-001" {
					po := &models.PurchaseOrder{
						TenantID:      tenantID,
						PRID:          pr.ID,
						RFQID:         rfq.ID,
						QuoteID:       quote.ID,
						PONumber:      "PO-001",
						Status:        "pending",
						TotalAmount:   3100.00,
						Currency:      "USD",
						PaymentMode:   "DIRECT",
						PaymentStatus: "succeeded",
						SupplierID:    supplierID,
						CreatedBy:     userID,
					}

					if err := poRepo.Create(po); err != nil {
						log.Printf("Failed to create PO: %v", err)
					} else {
						fmt.Printf("Created PO: %s\n", po.PONumber)

						// Add PO items
						poItems := []models.POItem{
							{
								POID:        po.ID,
								PRItemID:    items[0].ID,
								Description: "Printer Paper A4",
								Quantity:    100,
								UnitPrice:   24.00,
								TotalPrice:  2400.00,
							},
							{
								POID:        po.ID,
								PRItemID:    items[1].ID,
								Description: "Office Pens",
								Quantity:    50,
								UnitPrice:   14.00,
								TotalPrice:  700.00,
							},
						}

						for _, item := range poItems {
							if err := db.Create(&item).Error; err != nil {
								log.Printf("Failed to create PO item: %v", err)
							}
						}
					}
				}
			}
		}
		
		// Create second order for RFQ-002
		if len(rfqs) > 1 && rfqs[1].RFQNumber == "RFQ-002" {
			// Create second PO
			po2 := &models.PurchaseOrder{
				TenantID:      tenantID,
				PRID:          pr.ID,
				RFQID:         rfqs[1].ID,
				PONumber:      "PO-002",
				Status:        "pending",
				TotalAmount:   5500.00,
				Currency:      "USD",
				PaymentMode:   "DIRECT",
				PaymentStatus: "pending",
				SupplierID:    supplierID,
				CreatedBy:     userID,
			}

			if err := poRepo.Create(po2); err != nil {
				log.Printf("Failed to create second PO: %v", err)
			} else {
				fmt.Printf("Created PO: %s\n", po2.PONumber)
			}
		}
	}

	fmt.Println("Procurement seeding completed!")
}
