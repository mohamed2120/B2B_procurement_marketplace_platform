package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/notification-service/models"
	"github.com/b2b-platform/notification-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	templateRepo := repository.NewTemplateRepository(db)
	preferenceRepo := repository.NewPreferenceRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Create templates
	templates := []models.NotificationTemplate{
		{
			Code:      "order_placed",
			Name:      "Order Placed",
			Subject:   "New Order Placed",
			Body:      "A new order has been placed: {{order_number}}",
			BodyHTML:  "<p>A new order has been placed: <strong>{{order_number}}</strong></p>",
			Channel:   "email",
			EventType: "procurement.order.placed.v1",
			IsActive:  true,
		},
		{
			Code:      "shipment_late",
			Name:      "Shipment Late",
			Subject:   "Shipment Delayed",
			Body:      "Your shipment {{tracking_number}} is running late. Expected delivery: {{eta}}",
			BodyHTML:  "<p>Your shipment <strong>{{tracking_number}}</strong> is running late. Expected delivery: {{eta}}</p>",
			Channel:   "in_app",
			EventType: "logistics.shipment.late.v1",
			IsActive:  true,
		},
		{
			Code:      "quote_submitted",
			Name:      "Quote Submitted",
			Subject:   "New Quote Submitted",
			Body:      "A new quote has been submitted for RFQ {{rfq_number}}",
			BodyHTML:  "<p>A new quote has been submitted for RFQ <strong>{{rfq_number}}</strong></p>",
			Channel:   "in_app",
			EventType: "procurement.quote.submitted.v1",
			IsActive:  true,
		},
		{
			Code:      "pr_approved",
			Name:      "PR Approved",
			Subject:   "Purchase Request Approved",
			Body:      "Your purchase request {{pr_number}} has been approved",
			BodyHTML:  "<p>Your purchase request <strong>{{pr_number}}</strong> has been approved</p>",
			Channel:   "in_app",
			EventType: "procurement.pr.approved.v1",
			IsActive:  true,
		},
		{
			Code:      "part_approved",
			Name:      "Part Approved",
			Subject:   "Catalog Part Approved",
			Body:      "Your catalog part {{part_number}} - {{name}} has been approved",
			BodyHTML:  "<p>Your catalog part <strong>{{part_number}}</strong> - {{name}} has been approved</p>",
			Channel:   "in_app",
			EventType: "catalog.lib_part.approved.v1",
			IsActive:  true,
		},
	}

	for _, template := range templates {
		if err := templateRepo.Create(&template); err != nil {
			log.Printf("Template may already exist: %v", err)
		} else {
			fmt.Printf("Created template: %s\n", template.Name)
		}
	}

	// Create preference
	preference := &models.NotificationPreference{
		TenantID:  tenantID,
		UserID:    userID,
		Channel:   "email",
		EventType: "procurement.order.placed.v1",
		IsEnabled: true,
	}

	if err := preferenceRepo.Create(preference); err != nil {
		log.Printf("Preference may already exist: %v", err)
	} else {
		fmt.Printf("Created preference for user\n")
	}

	// Create sample notification
	notification := &models.Notification{
		TenantID: tenantID,
		UserID:   userID,
		Channel:  "in_app",
		Type:     "order_placed",
		Title:    "New Order Placed",
		Message:  "Order PO-001 has been placed successfully",
		Data:     "{}", // Valid JSON for JSONB field
		Status:   "sent",
	}

	if err := notificationRepo.Create(notification); err != nil {
		log.Printf("Notification may already exist: %v", err)
	} else {
		fmt.Printf("Created notification: %s\n", notification.Title)
	}

	fmt.Println("Notification seeding completed!")
}
