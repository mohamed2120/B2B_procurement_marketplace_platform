package service

import (
	"context"
	"fmt"
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type PaymentService struct {
	paymentRepo    *repository.PaymentRepository
	escrowRepo     *repository.EscrowRepository
	settlementRepo *repository.SettlementRepository
	refundRepo     *repository.RefundRepository
	payoutRepo     *repository.PayoutRepository
	provider       PaymentProvider
	eventBus       events.EventBus
	orderClient    *OrderClient
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	escrowRepo *repository.EscrowRepository,
	settlementRepo *repository.SettlementRepository,
	refundRepo *repository.RefundRepository,
	payoutRepo *repository.PayoutRepository,
	provider PaymentProvider,
	eventBus events.EventBus,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		escrowRepo:     escrowRepo,
		settlementRepo: settlementRepo,
		refundRepo:     refundRepo,
		payoutRepo:     payoutRepo,
		provider:       provider,
		eventBus:       eventBus,
		orderClient:    NewOrderClient(),
	}
}

type CreatePaymentIntentRequest struct {
	OrderID     uuid.UUID              `json:"order_id"`
	SupplierID  uuid.UUID              `json:"supplier_id"` // Required for ESCROW mode
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	PaymentMode string                 `json:"payment_mode"` // DIRECT, ESCROW
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CreatePaymentIntentResponse struct {
	PaymentIntentID string `json:"payment_intent_id"`
	ClientSecret    string `json:"client_secret"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

func (s *PaymentService) CreatePaymentIntent(ctx context.Context, tenantID uuid.UUID, req CreatePaymentIntentRequest) (*CreatePaymentIntentResponse, error) {
	// Create payment intent with provider
	intent, err := s.provider.CreatePaymentIntent(ctx, req.Amount, req.Currency, req.OrderID, req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create payment record
	payment := &models.Payment{
		TenantID:        tenantID,
		OrderID:         req.OrderID,
		PaymentIntentID: intent.ID,
		Provider:        "mock", // In production, use actual provider
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          "pending",
		PaymentMode:     req.PaymentMode,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// If ESCROW mode, create escrow hold
	if req.PaymentMode == "ESCROW" {
		if req.SupplierID == uuid.Nil {
			return nil, fmt.Errorf("supplier_id is required for ESCROW payment mode")
		}

		autoReleaseDays := 30 // Configurable
		autoReleaseDate := time.Now().Add(time.Duration(autoReleaseDays) * 24 * time.Hour)

		escrowHold := &models.EscrowHold{
			TenantID:        tenantID,
			PaymentID:       payment.ID,
			OrderID:         req.OrderID,
			SupplierID:      req.SupplierID,
			Amount:          req.Amount,
			Currency:        req.Currency,
			Status:          "pending", // Will be "held" when payment succeeds
			AutoReleaseDays: autoReleaseDays,
			AutoReleaseDate: &autoReleaseDate,
			BlockedByDispute: false,
		}

		if err := s.escrowRepo.Create(escrowHold); err != nil {
			return nil, fmt.Errorf("failed to create escrow hold: %w", err)
		}
	}

	return &CreatePaymentIntentResponse{
		PaymentIntentID: intent.ID,
		ClientSecret:    intent.ClientSecret,
		Amount:          req.Amount,
		Currency:        req.Currency,
	}, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := s.provider.HandleWebhook(ctx, payload, signature)
	if err != nil {
		return fmt.Errorf("failed to process webhook: %w", err)
	}

	// Find payment by intent ID
	payment, err := s.paymentRepo.GetByPaymentIntentID(event.PaymentIntentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentSucceeded(ctx, payment, event)
	case "payment_intent.payment_failed":
		return s.handlePaymentFailed(ctx, payment, event)
	default:
		return nil // Unknown event type, ignore
	}
}

func (s *PaymentService) handlePaymentSucceeded(ctx context.Context, payment *models.Payment, event *WebhookEvent) error {
	now := time.Now()
	payment.Status = "succeeded"
	payment.PaidAt = &now

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	// Update order payment status in procurement service
	// Note: In production, this would use service-to-service auth
	// For now, we'll skip if no token is available
	if err := s.orderClient.UpdateOrderPaymentStatus(payment.OrderID, "succeeded", ""); err != nil {
		// Log error but don't fail - order status update is best effort
	}

	// If ESCROW mode, update escrow hold status
	if payment.PaymentMode == "ESCROW" {
		escrowHold, err := s.escrowRepo.GetByPaymentID(payment.ID)
		if err == nil {
			escrowHold.Status = "held"
			if err := s.escrowRepo.Update(escrowHold); err != nil {
				return err
			}

			// Publish escrow held event
			event := events.NewEventEnvelope(
				events.EventEscrowHeld,
				"billing-service",
				map[string]interface{}{
					"escrow_hold_id": escrowHold.ID.String(),
					"payment_id":     payment.ID.String(),
					"order_id":       escrowHold.OrderID.String(),
					"amount":         escrowHold.Amount,
				},
			).WithTenantID(payment.TenantID)

			if err := s.eventBus.Publish(ctx, event); err != nil {
				// Log error but don't fail
			}
		}
	}

	// Publish payment succeeded event
	paymentEvent := events.NewEventEnvelope(
		events.EventPaymentSucceeded,
		"billing-service",
		map[string]interface{}{
			"payment_id":     payment.ID.String(),
			"payment_intent_id": payment.PaymentIntentID,
			"order_id":       payment.OrderID.String(),
			"amount":         payment.Amount,
			"payment_mode":   payment.PaymentMode,
		},
	).WithTenantID(payment.TenantID)

	return s.eventBus.Publish(ctx, paymentEvent)
}

func (s *PaymentService) handlePaymentFailed(ctx context.Context, payment *models.Payment, event *WebhookEvent) error {
	payment.Status = "failed"
	if reason, ok := event.Metadata["error"].(string); ok {
		payment.FailedReason = reason
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	// Update order payment status in procurement service
	if err := s.orderClient.UpdateOrderPaymentStatus(payment.OrderID, "failed", ""); err != nil {
		// Log error but don't fail
	}

	// Publish payment failed event
	failedEvent := events.NewEventEnvelope(
		events.EventPaymentFailed,
		"billing-service",
		map[string]interface{}{
			"payment_id":     payment.ID.String(),
			"payment_intent_id": payment.PaymentIntentID,
			"order_id":       payment.OrderID.String(),
			"failed_reason":   payment.FailedReason,
		},
	).WithTenantID(payment.TenantID)

	return s.eventBus.Publish(ctx, failedEvent)
}

func (s *PaymentService) ReleaseEscrow(ctx context.Context, escrowHoldID uuid.UUID, releasedBy uuid.UUID, reason string) error {
	escrowHold, err := s.escrowRepo.GetByID(escrowHoldID)
	if err != nil {
		return fmt.Errorf("escrow hold not found: %w", err)
	}

	// Check if blocked by dispute
	if escrowHold.BlockedByDispute {
		return fmt.Errorf("escrow release blocked by open dispute")
	}

	// Check if already released
	if escrowHold.Status != "held" {
		return fmt.Errorf("escrow hold is not in 'held' status")
	}

	// Get supplier payout account
	payoutAccount, err := s.payoutRepo.GetDefaultBySupplierID(escrowHold.SupplierID)
	if err != nil {
		return fmt.Errorf("no default payout account found for supplier: %w", err)
	}

	// Create settlement
	settlement := &models.Settlement{
		TenantID:        escrowHold.TenantID,
		EscrowHoldID:    escrowHold.ID,
		SupplierID:      escrowHold.SupplierID,
		PayoutAccountID: payoutAccount.ID,
		Amount:          escrowHold.Amount,
		Currency:        escrowHold.Currency,
		Status:          "pending",
	}

	if err := s.settlementRepo.Create(settlement); err != nil {
		return fmt.Errorf("failed to create settlement: %w", err)
	}

	// Update escrow hold
	now := time.Now()
	escrowHold.Status = "released"
	escrowHold.ReleasedAt = &now
	escrowHold.ReleasedBy = &releasedBy
	escrowHold.ReleaseReason = reason

	if err := s.escrowRepo.Update(escrowHold); err != nil {
		return err
	}

	// Process settlement (in production, this would call payment provider)
	settlement.Status = "completed"
	settlement.SettledAt = &now
	settlement.ProviderPayoutID = fmt.Sprintf("payout_%s", uuid.New().String()[:8])

	if err := s.settlementRepo.Update(settlement); err != nil {
		return err
	}

	// Publish settlement released event
	event := events.NewEventEnvelope(
		events.EventSettlementReleased,
		"billing-service",
		map[string]interface{}{
			"settlement_id":  settlement.ID.String(),
			"escrow_hold_id": escrowHold.ID.String(),
			"order_id":       escrowHold.OrderID.String(),
			"supplier_id":    escrowHold.SupplierID.String(),
			"amount":         settlement.Amount,
		},
	).WithTenantID(escrowHold.TenantID)

	return s.eventBus.Publish(ctx, event)
}

func (s *PaymentService) CreateRefund(ctx context.Context, tenantID uuid.UUID, paymentID uuid.UUID, amount float64, reason string, createdBy uuid.UUID) (*models.Refund, error) {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.TenantID != tenantID {
		return nil, fmt.Errorf("payment not found")
	}

	if payment.Status != "succeeded" {
		return nil, fmt.Errorf("can only refund succeeded payments")
	}

	// Create refund with provider
	refundResult, err := s.provider.Refund(ctx, payment.PaymentIntentID, amount, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	// Create refund record
	refundNumber := fmt.Sprintf("REF-%d", time.Now().Unix())
	refund := &models.Refund{
		TenantID:        tenantID,
		PaymentID:       paymentID,
		OrderID:         payment.OrderID,
		RefundNumber:    refundNumber,
		Amount:          amount,
		Currency:        payment.Currency,
		Reason:          reason,
		Status:          refundResult.Status,
		ProviderRefundID: refundResult.RefundID,
		CreatedBy:       createdBy,
	}

	if !refundResult.Success {
		refund.Status = "failed"
		refund.FailedReason = refundResult.FailedReason
	} else {
		now := time.Now()
		refund.Status = "completed"
		refund.RefundedAt = &now
	}

	if err := s.refundRepo.Create(refund); err != nil {
		return nil, fmt.Errorf("failed to create refund record: %w", err)
	}

	// Publish refund issued event
	event := events.NewEventEnvelope(
		events.EventRefundIssued,
		"billing-service",
		map[string]interface{}{
			"refund_id":    refund.ID.String(),
			"refund_number": refund.RefundNumber,
			"payment_id":   payment.ID.String(),
			"order_id":     payment.OrderID.String(),
			"amount":       refund.Amount,
		},
	).WithTenantID(tenantID)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail
	}

	return refund, nil
}

func (s *PaymentService) CheckDisputeStatus(orderID uuid.UUID) (bool, error) {
	// In production, this would call collaboration-service to check for open disputes
	// For now, return false (no dispute)
	return false, nil
}

func (s *PaymentService) ProcessAutoRelease(ctx context.Context) error {
	// Get all escrow holds that are eligible for auto-release
	// This would be called by a scheduled job
	holds, err := s.escrowRepo.ListPendingRelease(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		return err
	}

	for _, hold := range holds {
		// Check for disputes
		hasDispute, err := s.CheckDisputeStatus(hold.OrderID)
		if err != nil {
			continue
		}

		if hasDispute {
			hold.BlockedByDispute = true
			s.escrowRepo.Update(&hold)
			continue
		}

		// Auto-release
		if err := s.ReleaseEscrow(ctx, hold.ID, uuid.Nil, "Auto-released after grace period"); err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}

func (s *PaymentService) GetPayment(id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetByID(id)
}

func (s *PaymentService) ListPayments(tenantID uuid.UUID, limit, offset int) ([]models.Payment, error) {
	return s.paymentRepo.List(tenantID, limit, offset)
}

func (s *PaymentService) GetEscrowHold(id uuid.UUID) (*models.EscrowHold, error) {
	return s.escrowRepo.GetByID(id)
}

func (s *PaymentService) ListEscrowHolds(supplierID uuid.UUID, limit, offset int) ([]models.EscrowHold, error) {
	return s.escrowRepo.ListBySupplier(supplierID, limit, offset)
}

func (s *PaymentService) UpdateEscrowDisputeStatus(orderID uuid.UUID, hasDispute bool) error {
	escrowHold, err := s.escrowRepo.GetByOrderID(orderID)
	if err != nil {
		return nil // No escrow hold for this order
	}

	escrowHold.BlockedByDispute = hasDispute
	return s.escrowRepo.Update(escrowHold)
}
