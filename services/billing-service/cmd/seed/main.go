package main

import (
	"fmt"
	"log"
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	planRepo := repository.NewPlanRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	escrowRepo := repository.NewEscrowRepository(db)
	payoutRepo := repository.NewPayoutRepository(db)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	supplierID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Create a plan
	plan := &models.Plan{
		Name:        "Starter Plan",
		Code:        "starter",
		Description: "Starter plan for small businesses",
		Price:       99.00,
		Currency:    "USD",
		BillingCycle: "monthly",
		IsActive:    true,
	}

	if err := planRepo.Create(plan); err != nil {
		log.Printf("Plan may already exist: %v", err)
	} else {
		fmt.Printf("Created plan: %s\n", plan.Name)
	}

	// Create subscription
	subscription := &models.Subscription{
		TenantID:    tenantID,
		PlanID:      plan.ID,
		Status:      "active",
		StartedAt:   time.Now(),
		AutoRenew:   true,
	}

	if err := subscriptionRepo.Create(subscription); err != nil {
		log.Printf("Subscription may already exist: %v", err)
	} else {
		fmt.Printf("Created subscription for tenant\n")
	}

	// Create payout account for supplier
	payoutAccount := &models.PayoutAccount{
		TenantID:        tenantID,
		SupplierID:      supplierID,
		AccountType:     "bank_account",
		Provider:        "mock",
		AccountNumber:   "****1234",
		RoutingNumber:   "123456789",
		AccountHolderName: "Supplier Company",
		BankName:        "Mock Bank",
		IsDefault:       true,
		IsVerified:      true,
	}

	if err := payoutRepo.Create(payoutAccount); err != nil {
		log.Printf("Payout account may already exist: %v", err)
	} else {
		fmt.Printf("Created payout account for supplier\n")
	}

	// Create a direct payment (for direct order)
	directOrderID := uuid.New()
	directPayment := &models.Payment{
		TenantID:        tenantID,
		OrderID:         directOrderID,
		PaymentIntentID: "pi_direct_001",
		Provider:        "mock",
		Amount:          1500.00,
		Currency:        "USD",
		Status:          "succeeded",
		PaymentMode:     "DIRECT",
		PaidAt:          func() *time.Time { t := time.Now(); return &t }(),
	}

	if err := paymentRepo.Create(directPayment); err != nil {
		log.Printf("Direct payment may already exist: %v", err)
	} else {
		fmt.Printf("Created direct payment for order: %s\n", directOrderID.String())
	}

	// Create an escrow payment (for escrow order)
	// Note: In real scenario, this order would be created in procurement-service first
	escrowOrderID := uuid.MustParse("00000000-0000-0000-0000-000000000100") // Placeholder order ID
	escrowPayment := &models.Payment{
		TenantID:        tenantID,
		OrderID:         escrowOrderID,
		PaymentIntentID: "pi_escrow_001",
		Provider:        "mock",
		Amount:          2500.00,
		Currency:        "USD",
		Status:          "succeeded",
		PaymentMode:     "ESCROW",
		PaidAt:          func() *time.Time { t := time.Now(); return &t }(),
	}

	if err := paymentRepo.Create(escrowPayment); err != nil {
		log.Printf("Escrow payment may already exist: %v", err)
	} else {
		fmt.Printf("Created escrow payment for order: %s\n", escrowOrderID.String())

		// Create escrow hold
		autoReleaseDate := time.Now().Add(30 * 24 * time.Hour)
		escrowHold := &models.EscrowHold{
			TenantID:        tenantID,
			PaymentID:       escrowPayment.ID,
			OrderID:         escrowOrderID,
			SupplierID:      supplierID,
			Amount:          2500.00,
			Currency:        "USD",
			Status:          "held",
			AutoReleaseDays: 30,
			AutoReleaseDate: &autoReleaseDate,
			BlockedByDispute: false,
		}

		if err := escrowRepo.Create(escrowHold); err != nil {
			log.Printf("Escrow hold may already exist: %v", err)
		} else {
			fmt.Printf("Created escrow hold: %s\n", escrowHold.ID.String())
		}
	}

	fmt.Println("Billing seeding completed!")
}
