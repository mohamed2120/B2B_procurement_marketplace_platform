package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/company-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockCompanyRepository for testing
type MockCompanyRepository struct {
	companies        map[uuid.UUID]*models.Company
	subdomainRequests map[uuid.UUID]*models.SubdomainRequest
}

func (m *MockCompanyRepository) Create(company *models.Company) error {
	if m.companies == nil {
		m.companies = make(map[uuid.UUID]*models.Company)
	}
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) GetByID(id uuid.UUID) (*models.Company, error) {
	if company, ok := m.companies[id]; ok {
		return company, nil
	}
	return nil, nil
}

func (m *MockCompanyRepository) List(limit, offset int) ([]models.Company, error) {
	var result []models.Company
	for _, company := range m.companies {
		result = append(result, *company)
	}
	return result, nil
}

func (m *MockCompanyRepository) Update(company *models.Company) error {
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) CreateSubdomainRequest(req *models.SubdomainRequest) error {
	if m.subdomainRequests == nil {
		m.subdomainRequests = make(map[uuid.UUID]*models.SubdomainRequest)
	}
	m.subdomainRequests[req.ID] = req
	return nil
}

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

func TestCompanyService_Create(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	mockEventBus := &MockEventBus{}
	service := NewCompanyService(mockRepo, mockEventBus)

	company := &models.Company{
		ID:     uuid.New(),
		Name:   "Test Company",
		Status: "pending",
	}

	err := service.Create(company)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify company was created
	created, _ := mockRepo.GetByID(company.ID)
	if created == nil {
		t.Errorf("expected company to be created")
	}
	if created.Name != "Test Company" {
		t.Errorf("expected name 'Test Company', got %s", created.Name)
	}
}

func TestCompanyService_GetByID(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	mockEventBus := &MockEventBus{}
	service := NewCompanyService(mockRepo, mockEventBus)

	companyID := uuid.New()
	company := &models.Company{
		ID:     companyID,
		Name:   "Test Company",
		Status: "pending",
	}
	mockRepo.Create(company)

	result, err := service.GetByID(companyID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected company to be found")
	}
	if result.Name != "Test Company" {
		t.Errorf("expected name 'Test Company', got %s", result.Name)
	}
}

func TestCompanyService_Approve(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	mockEventBus := &MockEventBus{}
	service := NewCompanyService(mockRepo, mockEventBus)

	companyID := uuid.New()
	approvedBy := uuid.New()
	company := &models.Company{
		ID:       companyID,
		Name:     "Test Company",
		Status:   "pending",
		Subdomain: "test",
	}
	mockRepo.Create(company)

	err := service.Approve(companyID, approvedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify company was updated
	updated, _ := mockRepo.GetByID(companyID)
	if updated.Status != "approved" {
		t.Errorf("expected status 'approved', got %s", updated.Status)
	}
	if updated.VerificationStatus != "verified" {
		t.Errorf("expected verification status 'verified', got %s", updated.VerificationStatus)
	}
	if updated.ApprovedAt == nil {
		t.Errorf("expected ApprovedAt to be set")
	}
	if updated.ApprovedBy == nil || *updated.ApprovedBy != approvedBy {
		t.Errorf("expected ApprovedBy to be set")
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventCompanyApproved {
			t.Errorf("expected event type %s, got %s", events.EventCompanyApproved, event.Type)
		}
		if event.Payload["company_id"] != companyID.String() {
			t.Errorf("expected company_id %s, got %v", companyID.String(), event.Payload["company_id"])
		}
	}
}

func TestCompanyService_RequestSubdomain(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	mockEventBus := &MockEventBus{}
	service := NewCompanyService(mockRepo, mockEventBus)

	companyID := uuid.New()
	requestedBy := uuid.New()
	subdomain := "test-company"

	err := service.RequestSubdomain(companyID, subdomain, requestedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify subdomain request was created
	if len(mockRepo.subdomainRequests) == 0 {
		t.Errorf("expected subdomain request to be created")
	}
}

func TestCompanyService_Update(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	mockEventBus := &MockEventBus{}
	service := NewCompanyService(mockRepo, mockEventBus)

	companyID := uuid.New()
	company := &models.Company{
		ID:     companyID,
		Name:   "Test Company",
		Status: "pending",
	}
	mockRepo.Create(company)

	company.Name = "Updated Company"
	err := service.Update(company)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify company was updated
	updated, _ := mockRepo.GetByID(companyID)
	if updated.Name != "Updated Company" {
		t.Errorf("expected name 'Updated Company', got %s", updated.Name)
	}
}
