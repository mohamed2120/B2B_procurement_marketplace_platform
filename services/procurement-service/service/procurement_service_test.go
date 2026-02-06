package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/procurement-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockEventBus for testing
type MockEventBus struct {
	publishedEvents []*events.EventEnvelope
}

func (m *MockEventBus) Publish(ctx interface{}, event *events.EventEnvelope) error {
	if m.publishedEvents == nil {
		m.publishedEvents = make([]*events.EventEnvelope, 0)
	}
	m.publishedEvents = append(m.publishedEvents, event)
	return nil
}

func (m *MockEventBus) Subscribe(ctx interface{}, eventType events.EventType, handler func(*events.EventEnvelope) error) error {
	return nil
}

// MockPRRepository for testing
type MockPRRepository struct {
	prs map[uuid.UUID]*models.PurchaseRequest
}

func (m *MockPRRepository) Create(pr *models.PurchaseRequest) error {
	if m.prs == nil {
		m.prs = make(map[uuid.UUID]*models.PurchaseRequest)
	}
	m.prs[pr.ID] = pr
	return nil
}

func (m *MockPRRepository) GetByID(id uuid.UUID) (*models.PurchaseRequest, error) {
	if pr, ok := m.prs[id]; ok {
		return pr, nil
	}
	return nil, nil
}

func (m *MockPRRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.PurchaseRequest, error) {
	var result []models.PurchaseRequest
	for _, pr := range m.prs {
		if pr.TenantID == tenantID {
			result = append(result, *pr)
		}
	}
	return result, nil
}

func (m *MockPRRepository) Update(pr *models.PurchaseRequest) error {
	m.prs[pr.ID] = pr
	return nil
}

func (m *MockPRRepository) AddApproval(approval *models.PRApproval) error {
	return nil
}

func TestProcurementService_ApprovePR_PublishesEvent(t *testing.T) {
	mockEventBus := new(MockEventBus)
	mockPRRepo := &MockPRRepository{
		prs: make(map[uuid.UUID]*models.PurchaseRequest),
	}

	service := &ProcurementService{
		prRepo:   mockPRRepo,
		eventBus: mockEventBus,
	}

	tenantID := uuid.New()
	prID := uuid.New()
	approverID := uuid.New()

	pr := &models.PurchaseRequest{
		ID:       prID,
		TenantID: tenantID,
		PRNumber: "PR-001",
		Status:   "pending",
	}
	mockPRRepo.Create(pr)

	err := service.ApprovePR(prID, approverID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventPRApproved {
			t.Errorf("expected event type %s, got %s", events.EventPRApproved, event.Type)
		}
		if event.Payload["pr_id"] != prID.String() {
			t.Errorf("expected pr_id %s, got %v", prID.String(), event.Payload["pr_id"])
		}
	}

	// Verify PR was updated
	updatedPR, _ := mockPRRepo.GetByID(prID)
	if updatedPR.Status != "approved" {
		t.Errorf("expected status 'approved', got %s", updatedPR.Status)
	}
	if updatedPR.ApprovedAt == nil {
		t.Errorf("expected ApprovedAt to be set")
	}
}

func TestProcurementService_CreateRFQ_PublishesEvent(t *testing.T) {
	mockEventBus := new(MockEventBus)
	mockRFQRepo := &MockRFQRepository{
		rfqs: make(map[uuid.UUID]*models.RFQ),
	}

	service := &ProcurementService{
		rfqRepo:  mockRFQRepo,
		eventBus: mockEventBus,
	}

	tenantID := uuid.New()
	prID := uuid.New()
	rfqID := uuid.New()

	rfq := &models.RFQ{
		ID:        rfqID,
		TenantID:  tenantID,
		PRID:      prID,
		DueDate:   time.Now().Add(7 * 24 * time.Hour),
		CreatedBy: uuid.New(),
	}

	err := service.CreateRFQ(rfq)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventRFQCreated {
			t.Errorf("expected event type %s, got %s", events.EventRFQCreated, event.Type)
		}
		if event.Payload["rfq_id"] != rfqID.String() {
			t.Errorf("expected rfq_id %s, got %v", rfqID.String(), event.Payload["rfq_id"])
		}
	}
}

type MockRFQRepository struct {
	rfqs map[uuid.UUID]*models.RFQ
}

func (m *MockRFQRepository) Create(rfq *models.RFQ) error {
	if m.rfqs == nil {
		m.rfqs = make(map[uuid.UUID]*models.RFQ)
	}
	m.rfqs[rfq.ID] = rfq
	return nil
}

func (m *MockRFQRepository) GetByID(id uuid.UUID) (*models.RFQ, error) {
	if rfq, ok := m.rfqs[id]; ok {
		return rfq, nil
	}
	return nil, nil
}

func (m *MockRFQRepository) List(tenantID uuid.UUID) ([]models.RFQ, error) {
	var result []models.RFQ
	for _, rfq := range m.rfqs {
		if rfq.TenantID == tenantID {
			result = append(result, *rfq)
		}
	}
	return result, nil
}

func (m *MockRFQRepository) Update(rfq *models.RFQ) error {
	m.rfqs[rfq.ID] = rfq
	return nil
}

// MockQuoteRepository for testing
type MockQuoteRepository struct {
	quotes map[uuid.UUID]*models.Quote
}

func (m *MockQuoteRepository) Create(quote *models.Quote) error {
	if m.quotes == nil {
		m.quotes = make(map[uuid.UUID]*models.Quote)
	}
	m.quotes[quote.ID] = quote
	return nil
}

func (m *MockQuoteRepository) GetByID(id uuid.UUID) (*models.Quote, error) {
	if quote, ok := m.quotes[id]; ok {
		return quote, nil
	}
	return nil, nil
}

func (m *MockQuoteRepository) GetByRFQ(rfqID uuid.UUID) ([]models.Quote, error) {
	var result []models.Quote
	for _, quote := range m.quotes {
		if quote.RFQID == rfqID {
			result = append(result, *quote)
		}
	}
	return result, nil
}

func (m *MockQuoteRepository) Update(quote *models.Quote) error {
	m.quotes[quote.ID] = quote
	return nil
}

// MockPORepository for testing
type MockPORepository struct {
	pos map[uuid.UUID]*models.PurchaseOrder
}

func (m *MockPORepository) Create(po *models.PurchaseOrder) error {
	if m.pos == nil {
		m.pos = make(map[uuid.UUID]*models.PurchaseOrder)
	}
	m.pos[po.ID] = po
	return nil
}

func (m *MockPORepository) GetByID(id uuid.UUID) (*models.PurchaseOrder, error) {
	if po, ok := m.pos[id]; ok {
		return po, nil
	}
	return nil, nil
}

func (m *MockPORepository) List(tenantID uuid.UUID) ([]models.PurchaseOrder, error) {
	var result []models.PurchaseOrder
	for _, po := range m.pos {
		if po.TenantID == tenantID {
			result = append(result, *po)
		}
	}
	return result, nil
}

func (m *MockPORepository) Update(po *models.PurchaseOrder) error {
	m.pos[po.ID] = po
	return nil
}

func TestProcurementService_CreatePR(t *testing.T) {
	mockPRRepo := &MockPRRepository{}
	mockRFQRepo := &MockRFQRepository{}
	mockQuoteRepo := &MockQuoteRepository{}
	mockPORepo := &MockPORepository{}
	mockEventBus := &MockEventBus{}

	service := NewProcurementService(mockPRRepo, mockRFQRepo, mockQuoteRepo, mockPORepo, mockEventBus)

	tenantID := uuid.New()
	pr := &models.PurchaseRequest{
		ID:       uuid.New(),
		TenantID: tenantID,
		Status:   "pending",
	}

	err := service.CreatePR(pr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify PR was created with PR number
	created, _ := mockPRRepo.GetByID(pr.ID)
	if created == nil {
		t.Errorf("expected PR to be created")
	}
	if created.PRNumber == "" {
		t.Errorf("expected PR number to be generated")
	}
}

func TestProcurementService_SubmitQuote_PublishesEvent(t *testing.T) {
	mockPRRepo := &MockPRRepository{}
	mockRFQRepo := &MockRFQRepository{}
	mockQuoteRepo := &MockQuoteRepository{}
	mockPORepo := &MockPORepository{}
	mockEventBus := &MockEventBus{}

	service := NewProcurementService(mockPRRepo, mockRFQRepo, mockQuoteRepo, mockPORepo, mockEventBus)

	tenantID := uuid.New()
	rfqID := uuid.New()
	quoteID := uuid.New()
	quote := &models.Quote{
		ID:       quoteID,
		TenantID: tenantID,
		RFQID:    rfqID,
		Status:   "submitted",
		Amount:   1000.00,
	}

	err := service.SubmitQuote(quote)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify quote was created
	created, _ := mockQuoteRepo.GetByID(quoteID)
	if created == nil {
		t.Errorf("expected quote to be created")
	}
	if created.QuoteNumber == "" {
		t.Errorf("expected quote number to be generated")
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventQuoteSubmitted {
			t.Errorf("expected event type %s, got %s", events.EventQuoteSubmitted, event.Type)
		}
	}
}

func (m *MockRFQRepository) Update(rfq *models.RFQ) error {
	m.rfqs[rfq.ID] = rfq
	return nil
}

// MockQuoteRepository for testing
type MockQuoteRepository struct {
	quotes map[uuid.UUID]*models.Quote
}

func (m *MockQuoteRepository) Create(quote *models.Quote) error {
	if m.quotes == nil {
		m.quotes = make(map[uuid.UUID]*models.Quote)
	}
	m.quotes[quote.ID] = quote
	return nil
}

func (m *MockQuoteRepository) GetByID(id uuid.UUID) (*models.Quote, error) {
	if quote, ok := m.quotes[id]; ok {
		return quote, nil
	}
	return nil, nil
}

func (m *MockQuoteRepository) GetByRFQ(rfqID uuid.UUID) ([]models.Quote, error) {
	var result []models.Quote
	for _, quote := range m.quotes {
		if quote.RFQID == rfqID {
			result = append(result, *quote)
		}
	}
	return result, nil
}

func (m *MockQuoteRepository) Update(quote *models.Quote) error {
	m.quotes[quote.ID] = quote
	return nil
}

// MockPORepository for testing
type MockPORepository struct {
	pos map[uuid.UUID]*models.PurchaseOrder
}

func (m *MockPORepository) Create(po *models.PurchaseOrder) error {
	if m.pos == nil {
		m.pos = make(map[uuid.UUID]*models.PurchaseOrder)
	}
	m.pos[po.ID] = po
	return nil
}

func (m *MockPORepository) GetByID(id uuid.UUID) (*models.PurchaseOrder, error) {
	if po, ok := m.pos[id]; ok {
		return po, nil
	}
	return nil, nil
}

func (m *MockPORepository) List(tenantID uuid.UUID) ([]models.PurchaseOrder, error) {
	var result []models.PurchaseOrder
	for _, po := range m.pos {
		if po.TenantID == tenantID {
			result = append(result, *po)
		}
	}
	return result, nil
}

func (m *MockPORepository) Update(po *models.PurchaseOrder) error {
	m.pos[po.ID] = po
	return nil
}

func TestProcurementService_CreatePR(t *testing.T) {
	mockPRRepo := &MockPRRepository{}
	mockRFQRepo := &MockRFQRepository{}
	mockQuoteRepo := &MockQuoteRepository{}
	mockPORepo := &MockPORepository{}
	mockEventBus := &MockEventBus{}

	service := NewProcurementService(mockPRRepo, mockRFQRepo, mockQuoteRepo, mockPORepo, mockEventBus)

	tenantID := uuid.New()
	pr := &models.PurchaseRequest{
		ID:       uuid.New(),
		TenantID: tenantID,
		Status:   "pending",
	}

	err := service.CreatePR(pr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify PR was created with PR number
	created, _ := mockPRRepo.GetByID(pr.ID)
	if created == nil {
		t.Errorf("expected PR to be created")
	}
	if created.PRNumber == "" {
		t.Errorf("expected PR number to be generated")
	}
}

func TestProcurementService_SubmitQuote_PublishesEvent(t *testing.T) {
	mockPRRepo := &MockPRRepository{}
	mockRFQRepo := &MockRFQRepository{}
	mockQuoteRepo := &MockQuoteRepository{}
	mockPORepo := &MockPORepository{}
	mockEventBus := &MockEventBus{}

	service := NewProcurementService(mockPRRepo, mockRFQRepo, mockQuoteRepo, mockPORepo, mockEventBus)

	rfqID := uuid.New()
	quoteID := uuid.New()
	quote := &models.Quote{
		ID:       quoteID,
		RFQID:    rfqID,
		Status:   "submitted",
		Amount:   1000.00,
	}

	err := service.SubmitQuote(quote)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify quote was created
	created, _ := mockQuoteRepo.GetByID(quoteID)
	if created == nil {
		t.Errorf("expected quote to be created")
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventQuoteSubmitted {
			t.Errorf("expected event type %s, got %s", events.EventQuoteSubmitted, event.Type)
		}
	}
}
