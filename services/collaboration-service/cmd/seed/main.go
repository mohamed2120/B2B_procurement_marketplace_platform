package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/collaboration-service/models"
	"github.com/b2b-platform/collaboration-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	threadRepo := repository.NewThreadRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	disputeRepo := repository.NewDisputeRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID1 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	userID2 := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000010")

	// Create chat thread
	thread := &models.ChatThread{
		TenantID:    tenantID,
		Title:       "Order Discussion",
		ThreadType:  "order",
		ReferenceID: &orderID,
		CreatedBy:   userID1,
	}

	if err := threadRepo.Create(thread); err != nil {
		log.Printf("Thread may already exist: %v", err)
	} else {
		fmt.Printf("Created thread: %s\n", thread.Title)

		// Add participants
		participants := []models.ThreadParticipant{
			{ThreadID: thread.ID, UserID: userID1, TenantID: tenantID, Role: "buyer"},
			{ThreadID: thread.ID, UserID: userID2, TenantID: tenantID, Role: "supplier"},
		}

		for _, p := range participants {
			if err := threadRepo.AddParticipant(&p); err != nil {
				log.Printf("Failed to add participant: %v", err)
			}
		}

		// Add messages
		messages := []models.ChatMessage{
			{
				ThreadID:    thread.ID,
				SenderID:    userID1,
				Message:     "Hello, I have a question about the order",
				MessageType: "text",
			},
			{
				ThreadID:    thread.ID,
				SenderID:    userID2,
				Message:     "Sure, how can I help?",
				MessageType: "text",
			},
		}

		for _, msg := range messages {
			if err := messageRepo.Create(&msg); err != nil {
				log.Printf("Failed to create message: %v", err)
			} else {
				fmt.Printf("Created message: %s\n", msg.Message[:20])
			}
		}
	}

	// Create dispute
	dispute := &models.Dispute{
		TenantID:    tenantID,
		OrderID:     orderID,
		DisputeType: "quality",
		Status:      "open",
		Title:       "Quality Issue with Delivered Parts",
		Description: "The parts received do not meet the quality standards",
		RaisedBy:    userID1,
	}

	if err := disputeRepo.Create(dispute); err != nil {
		log.Printf("Dispute may already exist: %v", err)
	} else {
		fmt.Printf("Created dispute: %s\n", dispute.Title)
	}

	// Create rating
	rating := &models.Rating{
		TenantID:        tenantID,
		OrderID:         orderID,
		RatedBy:         userID1,
		RatedEntityType: "supplier",
		RatedEntityID:  userID2,
		Rating:          4,
		Comment:         "Good service, fast delivery",
	}

	if err := ratingRepo.Create(rating); err != nil {
		log.Printf("Rating may already exist: %v", err)
	} else {
		fmt.Printf("Created rating: %d stars\n", rating.Rating)
	}

	fmt.Println("Collaboration seeding completed!")
}
