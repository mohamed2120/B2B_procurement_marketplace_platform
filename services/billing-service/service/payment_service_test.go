package service

import (
	"context"
	"testing"
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
)

// MockPaymentRepository for testing
type MockPaymentRepository struct {
	payments map[uuid.UUID]*models.Payment
}

func (m *MockPaymentRepository) Create(payment *models.Payment) error {
	if m.payments == nil {
		m.payments = make(map[uuid.UUID]*models.Payment)
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *MockPaymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	if payment, ok := m.payments[id]; ok {
		return payment, nil
	}
	return nil, nil
}

func (m *MockPaymentRepository) GetByPaymentIntentID(paymentIntentID string) (*models.Payment, error) {
	for _, payment := range m.payments {
		if payment.PaymentIntentID == paymentIntentID {
			return payment, nil
		}
	}
	return nil, nil
}

func (m *MockPaymentRepository) GetByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	for _, payment := range m.payments {
		if payment.OrderID == orderID {
			return payment, nil
		}
	}
	return nil, nil
}

func (m *MockPaymentRepository) Update(payment *models.Payment) error {
	m.payments[payment.ID] = payment
	return nil
}

func (m *MockPaymentRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Payment, error) {
	var result []models.Payment
	for _, payment := range m.payments {
		if payment.TenantID == tenantID {
			result = append(result, *payment)
		}
	}
	return result, nil
}

// MockEscrowRepository for testing
type MockEscrowRepository struct {
	holds map[uuid.UUID]*models.EscrowHold
}

func (m *MockEscrowRepository) Create(hold *models.EscrowHold) error {
	if m.holds == nil {
		m.holds = make(map[uuid.UUID]*models.EscrowHold)
	}
	m.holds[hold.ID] = hold
	return nil
}

func (m *MockEscrowRepository) GetByID(id uuid.UUID) (*models.EscrowHold, error) {
	if hold, ok := m.holds[id]; ok {
		return hold, nil
	}
	return nil, nil
}

func (m *MockEscrowRepository) GetByPaymentID(paymentID uuid.UUID) (*models.EscrowHold, error) {
	for _, hold := range m.holds {
		if hold.PaymentID == paymentID {
			return hold, nil
		}
	}
	return nil, nil
}

func (m *MockEscrowRepository) GetByOrderID(orderID uuid.UUID) (*models.EscrowHold, error) {
	for _, hold := range m.holds {
		if hold.OrderID == orderID {
			return hold, nil
		}
	}
	return nil, nil
}

func (m *MockEscrowRepository) Update(hold *models.EscrowHold) error {
	m.holds[hold.ID] = hold
	return nil
}

func (m *MockEscrowRepository) ListPendingRelease(tenantID uuid.UUID) ([]models.EscrowHold, error) {
	var result []models.EscrowHold
	now := time.Now()
	for _, hold := range m.holds {
		if hold.TenantID == tenantID && hold.Status == "held" && !hold.BlockedByDispute {
			if hold.AutoReleaseDate != nil && hold.AutoReleaseDate.Before(now) {
				result = append(result, *hold)
			}
		}
	}
	return result, nil
}

func (m *MockEscrowRepository) ListBySupplier(supplierID uuid.UUID, limit, offset int) ([]models.EscrowHold, error) {
	var result []models.EscrowHold
	for _, hold := range m.holds {
		if hold.SupplierID == supplierID {
			result = append(result, *hold)
		}
	}
	return result, nil
}

// MockPayoutRepository for testing
type MockPayoutRepository struct {
	accounts map[uuid.UUID]*models.PayoutAccount
}

func (m *MockPayoutRepository) Create(account *models.PayoutAccount) error {
	if m.accounts == nil {
		m.accounts = make(map[uuid.UUID]*models.PayoutAccount)
	}
	m.accounts[account.ID] = account
	return nil
}

func (m *MockPayoutRepository) GetByID(id uuid.UUID) (*models.PayoutAccount, error) {
	if account, ok := m.accounts[id]; ok {
		return account, nil
	}
	return nil, nil
}

func (m *MockPayoutRepository) GetBySupplierID(supplierID uuid.UUID) ([]models.PayoutAccount, error) {
	var result []models.PayoutAccount
	for _, account := range m.accounts {
		if account.SupplierID == supplierID {
			result = append(result, *account)
		}
	}
	return result, nil
}

func (m *MockPayoutRepository) GetDefaultBySupplierID(supplierID uuid.UUID) (*models.PayoutAccount, error) {
	for _, account := range m.accounts {
		if account.SupplierID == supplierID && account.IsDefault {
			return account, nil
		}
	}
	return nil, nil
}

func (m *MockPayoutRepository) Update(account *models.PayoutAccount) error {
	m.accounts[account.ID] = account
	return nil
}

func (m *MockPayoutRepository) Delete(id uuid.UUID) error {
	delete(m.accounts, id)
	return nil
}

func (m *MockPayoutRepository) SetDefault(supplierID uuid.UUID, accountID uuid.UUID) error {
	for _, account := range m.accounts {
		if account.SupplierID == supplierID {
			account.IsDefault = (account.ID == accountID)
		}
	}
	return nil
}

// MockSettlementRepository for testing
type MockSettlementRepository struct {
	settlements map[uuid.UUID]*models.Settlement
}

func (m *MockSettlementRepository) Create(settlement *models.Settlement) error {
	if m.settlements == nil {
		m.settlements = make(map[uuid.UUID]*models.Settlement)
	}
	m.settlements[settlement.ID] = settlement
	return nil
}

func (m *MockSettlementRepository) GetByID(id uuid.UUID) (*models.Settlement, error) {
	if settlement, ok := m.settlements[id]; ok {
		return settlement, nil
	}
	return nil, nil
}

func (m *MockSettlementRepository) GetByEscrowHoldID(escrowHoldID uuid.UUID) (*models.Settlement, error) {
	for _, settlement := range m.settlements {
		if settlement.EscrowHoldID == escrowHoldID {
			return settlement, nil
		}
	}
	return nil, nil
}

func (m *MockSettlementRepository) Update(settlement *models.Settlement) error {
	m.settlements[settlement.ID] = settlement
	return nil
}

func (m *MockSettlementRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	var result []models.Settlement
	for _, settlement := range m.settlements {
		if settlement.TenantID == tenantID {
			result = append(result, *settlement)
		}
	}
	return result, nil
}

// MockRefundRepository for testing
type MockRefundRepository struct {
	refunds map[uuid.UUID]*models.Refund
}

func (m *MockRefundRepository) Create(refund *models.Refund) error {
	if m.refunds == nil {
		m.refunds = make(map[uuid.UUID]*models.Refund)
	}
	m.refunds[refund.ID] = refund
	return nil
}

func (m *MockRefundRepository) GetByID(id uuid.UUID) (*models.Refund, error) {
	if refund, ok := m.refunds[id]; ok {
		return refund, nil
	}
	return nil, nil
}

func (m *MockRefundRepository) GetByRefundNumber(tenantID uuid.UUID, refundNumber string) (*models.Refund, error) {
	for _, refund := range m.refunds {
		if refund.TenantID == tenantID && refund.RefundNumber == refundNumber {
			return refund, nil
		}
	}
	return nil, nil
}

func (m *MockRefundRepository) Update(refund *models.Refund) error {
	m.refunds[refund.ID] = refund
	return nil
}

func (m *MockRefundRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Refund, error) {
	var result []models.Refund
	for _, refund := range m.refunds {
		if refund.TenantID == tenantID {
			result = append(result, *refund)
		}
	}
	return result, nil
}

func TestPaymentService_CreatePaymentIntent_ESCROW(t *testing.T) {
	mockEventBus := &MockEventBus{}
	mockPaymentRepo := &MockPaymentRepository{}
	mockEscrowRepo := &MockEscrowRepository{}
	mockSettlementRepo := &MockSettlementRepository{}
	mockRefundRepo := &MockRefundRepository{}
	mockPayoutRepo := &MockPayoutRepository{}
	mockProvider := NewMockPaymentProvider()

	service := NewPaymentService(
		mockPaymentRepo,
		mockEscrowRepo,
		mockSettlementRepo,
		mockRefundRepo,
		mockPayoutRepo,
		mockProvider,
		mockEventBus,
	)

	tenantID := uuid.New()
	orderID := uuid.New()
	supplierID := uuid.New()

	req := CreatePaymentIntentRequest{
		OrderID:     orderID,
		SupplierID:  supplierID,
		Amount:      1000.00,
		Currency:    "USD",
		PaymentMode: "ESCROW",
	}

	response, err := service.CreatePaymentIntent(context.Background(), tenantID, req)
	if err != nil {
		t.Fatalf("Failed to create payment intent: %v", err)
	}

	if response.PaymentIntentID == "" {
		t.Error("Payment intent ID should not be empty")
	}

	// Verify escrow hold was created
	holds := mockEscrowRepo.holds
	if len(holds) == 0 {
		t.Error("Escrow hold should have been created")
	}
}

func TestPaymentService_ReleaseEscrow_BlockedByDispute(t *testing.T) {
	mockEventBus := &MockEventBus{}
	mockPaymentRepo := &MockPaymentRepository{}
	mockEscrowRepo := &MockEscrowRepository{}
	mockSettlementRepo := &MockSettlementRepository{}
	mockRefundRepo := &MockRefundRepository{}
	mockPayoutRepo := &MockPayoutRepository{}
	mockProvider := NewMockPaymentProvider()

	// Create payout account
	payoutAccount := &models.PayoutAccount{
		ID:            uuid.New(),
		SupplierID:   uuid.New(),
		IsDefault:     true,
		IsVerified:    true,
	}
	mockPayoutRepo.Create(payoutAccount)

	service := NewPaymentService(
		mockPaymentRepo,
		mockEscrowRepo,
		mockSettlementRepo,
		mockRefundRepo,
		mockPayoutRepo,
		mockProvider,
		mockEventBus,
	)

	// Create escrow hold with dispute
	escrowHold := &models.EscrowHold{
		ID:               uuid.New(),
		TenantID:         uuid.New(),
		PaymentID:        uuid.New(),
		OrderID:          uuid.New(),
		SupplierID:       payoutAccount.SupplierID,
		Amount:           1000.00,
		Status:           "held",
		BlockedByDispute: true,
	}
	mockEscrowRepo.Create(escrowHold)

	// Try to release - should fail
	err := service.ReleaseEscrow(context.Background(), escrowHold.ID, uuid.New(), "Test release")
	if err == nil {
		t.Error("Release should have failed due to dispute")
	}

	if err != nil && err.Error() != "escrow release blocked by open dispute" {
		t.Errorf("Expected dispute error, got: %v", err)
	}
}

func TestPaymentService_ReleaseEscrow_Success(t *testing.T) {
	mockEventBus := &MockEventBus{}
	mockPaymentRepo := &MockPaymentRepository{}
	mockEscrowRepo := &MockEscrowRepository{}
	mockSettlementRepo := &MockSettlementRepository{}
	mockRefundRepo := &MockRefundRepository{}
	mockPayoutRepo := &MockPayoutRepository{}
	mockProvider := NewMockPaymentProvider()

	// Create payout account
	supplierID := uuid.New()
	payoutAccount := &models.PayoutAccount{
		ID:          uuid.New(),
		SupplierID:  supplierID,
		IsDefault:   true,
		IsVerified:  true,
	}
	mockPayoutRepo.Create(payoutAccount)

	service := NewPaymentService(
		mockPaymentRepo,
		mockEscrowRepo,
		mockSettlementRepo,
		mockRefundRepo,
		mockPayoutRepo,
		mockProvider,
		mockEventBus,
	)

	// Create escrow hold without dispute
	escrowHold := &models.EscrowHold{
		ID:               uuid.New(),
		TenantID:         uuid.New(),
		PaymentID:        uuid.New(),
		OrderID:          uuid.New(),
		SupplierID:       supplierID,
		Amount:           1000.00,
		Status:           "held",
		BlockedByDispute: false,
	}
	mockEscrowRepo.Create(escrowHold)

	// Release should succeed
	err := service.ReleaseEscrow(context.Background(), escrowHold.ID, uuid.New(), "Test release")
	if err != nil {
		t.Fatalf("Release should have succeeded: %v", err)
	}

	// Verify settlement was created
	settlements := mockSettlementRepo.settlements
	if len(settlements) == 0 {
		t.Error("Settlement should have been created")
	}

	// Verify escrow hold status updated
	updatedHold, _ := mockEscrowRepo.GetByID(escrowHold.ID)
	if updatedHold.Status != "released" {
		t.Error("Escrow hold status should be 'released'")
	}
}
