package service

import (
	"fmt"
	"time"

	"github.com/b2b-platform/procurement-service/models"
	"github.com/b2b-platform/procurement-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type ProcurementService struct {
	prRepo        *repository.PRRepository
	rfqRepo       *repository.RFQRepository
	quoteRepo     *repository.QuoteRepository
	poRepo        *repository.PORepository
	eventBus      events.EventBus
	billingClient *BillingClient
}

func NewProcurementService(
	prRepo *repository.PRRepository,
	rfqRepo *repository.RFQRepository,
	quoteRepo *repository.QuoteRepository,
	poRepo *repository.PORepository,
	eventBus events.EventBus,
) *ProcurementService {
	return &ProcurementService{
		prRepo:        prRepo,
		rfqRepo:       rfqRepo,
		quoteRepo:     quoteRepo,
		poRepo:        poRepo,
		eventBus:      eventBus,
		billingClient: NewBillingClient(),
	}
}

func (s *ProcurementService) CreatePR(pr *models.PurchaseRequest) error {
	// Generate PR number
	pr.PRNumber = fmt.Sprintf("PR-%d", time.Now().Unix())
	return s.prRepo.Create(pr)
}

func (s *ProcurementService) GetPR(id uuid.UUID) (*models.PurchaseRequest, error) {
	return s.prRepo.GetByID(id)
}

func (s *ProcurementService) ListPRs(tenantID uuid.UUID, limit, offset int) ([]models.PurchaseRequest, error) {
	return s.prRepo.List(tenantID, limit, offset)
}

func (s *ProcurementService) ApprovePR(prID, approverID uuid.UUID) error {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return err
	}

	now := time.Now()
	pr.Status = "approved"
	pr.ApprovedAt = &now
	pr.ApprovedBy = &approverID

	if err := s.prRepo.Update(pr); err != nil {
		return err
	}

	// Add approval record
	approval := &models.PRApproval{
		PRID:       prID,
		ApproverID: approverID,
		Status:     "approved",
		ApprovedAt: &now,
	}
	s.prRepo.AddApproval(approval)

	// Publish event
	event := events.NewEventEnvelope(
		events.EventPRApproved,
		"procurement-service",
		map[string]interface{}{
			"pr_id":    pr.ID.String(),
			"pr_number": pr.PRNumber,
		},
	).WithTenantID(pr.TenantID)

	return s.eventBus.Publish(nil, event)
}

func (s *ProcurementService) CreateRFQ(rfq *models.RFQ) error {
	rfq.RFQNumber = fmt.Sprintf("RFQ-%d", time.Now().Unix())
	rfq.Status = "sent"

	if err := s.rfqRepo.Create(rfq); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventRFQCreated,
		"procurement-service",
		map[string]interface{}{
			"rfq_id":    rfq.ID.String(),
			"rfq_number": rfq.RFQNumber,
			"pr_id":     rfq.PRID.String(),
		},
	).WithTenantID(rfq.TenantID)

	return s.eventBus.Publish(nil, event)
}

func (s *ProcurementService) GetRFQ(id uuid.UUID) (*models.RFQ, error) {
	return s.rfqRepo.GetByID(id)
}

func (s *ProcurementService) GetQuote(id uuid.UUID) (*models.Quote, error) {
	return s.quoteRepo.GetByID(id)
}

func (s *ProcurementService) GetPO(id uuid.UUID) (*models.PurchaseOrder, error) {
	return s.poRepo.GetByID(id)
}

func (s *ProcurementService) UpdatePO(po *models.PurchaseOrder) error {
	return s.poRepo.Update(po)
}

func (s *ProcurementService) SubmitQuote(quote *models.Quote) error {
	quote.QuoteNumber = fmt.Sprintf("QT-%d", time.Now().Unix())
	quote.Status = "submitted"
	quote.SubmittedAt = time.Now()

	if err := s.quoteRepo.Create(quote); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventQuoteSubmitted,
		"procurement-service",
		map[string]interface{}{
			"quote_id":    quote.ID.String(),
			"quote_number": quote.QuoteNumber,
			"rfq_id":      quote.RFQID.String(),
			"supplier_id": quote.SupplierID.String(),
		},
	).WithTenantID(quote.TenantID)

	return s.eventBus.Publish(nil, event)
}

func (s *ProcurementService) CreatePO(po *models.PurchaseOrder, authToken string) error {
	po.PONumber = fmt.Sprintf("PO-%d", time.Now().Unix())
	po.Status = "pending"
	
	// Set default payment_mode if not provided
	if po.PaymentMode == "" {
		po.PaymentMode = "DIRECT"
	}
	
	// Set default payment_status
	if po.PaymentStatus == "" {
		po.PaymentStatus = "pending"
	}

	if err := s.poRepo.Create(po); err != nil {
		return err
	}

	// If ESCROW mode, create payment intent with billing service
	if po.PaymentMode == "ESCROW" {
		paymentIntentReq := CreatePaymentIntentRequest{
			OrderID:     po.ID,
			SupplierID:  po.SupplierID,
			Amount:      po.TotalAmount,
			Currency:    po.Currency,
			PaymentMode: "ESCROW",
			Metadata: map[string]interface{}{
				"po_number": po.PONumber,
				"pr_id":     po.PRID.String(),
			},
		}

		_, err := s.billingClient.CreatePaymentIntent(authToken, paymentIntentReq)
		if err != nil {
			// Update payment status to failed
			po.PaymentStatus = "failed"
			s.poRepo.Update(po)
			return fmt.Errorf("failed to create payment intent: %w", err)
		}

		// Update payment status to processing
		po.PaymentStatus = "processing"
		if err := s.poRepo.Update(po); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventOrderPlaced,
		"procurement-service",
		map[string]interface{}{
			"po_id":          po.ID.String(),
			"po_number":      po.PONumber,
			"pr_id":          po.PRID.String(),
			"quote_id":       po.QuoteID.String(),
			"payment_mode":   po.PaymentMode,
			"payment_status": po.PaymentStatus,
		},
	).WithTenantID(po.TenantID)

	return s.eventBus.Publish(nil, event)
}
