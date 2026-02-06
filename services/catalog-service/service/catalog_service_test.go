package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/catalog-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockManufacturerRepository for testing
type MockManufacturerRepository struct {
	manufacturers map[uuid.UUID]*models.Manufacturer
}

func (m *MockManufacturerRepository) Create(manufacturer *models.Manufacturer) error {
	if m.manufacturers == nil {
		m.manufacturers = make(map[uuid.UUID]*models.Manufacturer)
	}
	m.manufacturers[manufacturer.ID] = manufacturer
	return nil
}

func (m *MockManufacturerRepository) GetByID(id uuid.UUID) (*models.Manufacturer, error) {
	if mfr, ok := m.manufacturers[id]; ok {
		return mfr, nil
	}
	return nil, nil
}

func (m *MockManufacturerRepository) List() ([]models.Manufacturer, error) {
	var result []models.Manufacturer
	for _, mfr := range m.manufacturers {
		result = append(result, *mfr)
	}
	return result, nil
}

// MockCategoryRepository for testing
type MockCategoryRepository struct {
	categories map[uuid.UUID]*models.Category
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	if m.categories == nil {
		m.categories = make(map[uuid.UUID]*models.Category)
	}
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	if cat, ok := m.categories[id]; ok {
		return cat, nil
	}
	return nil, nil
}

func (m *MockCategoryRepository) List() ([]models.Category, error) {
	var result []models.Category
	for _, cat := range m.categories {
		result = append(result, *cat)
	}
	return result, nil
}

// MockAttributeRepository for testing
type MockAttributeRepository struct {
	attributes map[uuid.UUID]*models.Attribute
	partAttributes []models.PartAttribute
}

func (m *MockAttributeRepository) Create(attribute *models.Attribute) error {
	if m.attributes == nil {
		m.attributes = make(map[uuid.UUID]*models.Attribute)
	}
	m.attributes[attribute.ID] = attribute
	return nil
}

func (m *MockAttributeRepository) List() ([]models.Attribute, error) {
	var result []models.Attribute
	for _, attr := range m.attributes {
		result = append(result, *attr)
	}
	return result, nil
}

func (m *MockAttributeRepository) AddPartAttribute(partAttr *models.PartAttribute) error {
	if m.partAttributes == nil {
		m.partAttributes = make([]models.PartAttribute, 0)
	}
	m.partAttributes = append(m.partAttributes, *partAttr)
	return nil
}

// MockPartRepository for testing
type MockPartRepository struct {
	parts map[uuid.UUID]*models.LibraryPart
}

func (m *MockPartRepository) Create(part *models.LibraryPart) error {
	if m.parts == nil {
		m.parts = make(map[uuid.UUID]*models.LibraryPart)
	}
	m.parts[part.ID] = part
	return nil
}

func (m *MockPartRepository) GetByID(id uuid.UUID) (*models.LibraryPart, error) {
	if part, ok := m.parts[id]; ok {
		return part, nil
	}
	return nil, nil
}

func (m *MockPartRepository) List(limit, offset int, status string) ([]models.LibraryPart, error) {
	var result []models.LibraryPart
	for _, part := range m.parts {
		if status == "" || part.Status == status {
			result = append(result, *part)
		}
	}
	return result, nil
}

func (m *MockPartRepository) Update(part *models.LibraryPart) error {
	m.parts[part.ID] = part
	return nil
}

func (m *MockPartRepository) FindDuplicate(partNumber string, manufacturerID uuid.UUID) (*models.LibraryPart, error) {
	for _, part := range m.parts {
		if part.PartNumber == partNumber && part.ManufacturerID == manufacturerID {
			return part, nil
		}
	}
	return nil, nil
}

func (m *MockPartRepository) MarkAsDuplicate(partID, originalID uuid.UUID) error {
	if part, ok := m.parts[partID]; ok {
		part.IsDuplicate = true
		part.OriginalPartID = &originalID
	}
	return nil
}

func (m *MockPartRepository) GetPendingApproval() ([]models.LibraryPart, error) {
	var result []models.LibraryPart
	for _, part := range m.parts {
		if part.Status == "pending" {
			result = append(result, *part)
		}
	}
	return result, nil
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

func TestCatalogService_CreatePart_NoDuplicate(t *testing.T) {
	mockMfrRepo := &MockManufacturerRepository{}
	mockCatRepo := &MockCategoryRepository{}
	mockPartRepo := &MockPartRepository{}
	mockAttrRepo := &MockAttributeRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCatalogService(mockMfrRepo, mockCatRepo, mockPartRepo, mockAttrRepo, mockEventBus)

	manufacturerID := uuid.New()
	part := &models.LibraryPart{
		ID:            uuid.New(),
		PartNumber:    "PART-001",
		Name:          "Test Part",
		ManufacturerID: manufacturerID,
		Status:        "pending",
	}

	err := service.CreatePart(part)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify part was created
	created, _ := mockPartRepo.GetByID(part.ID)
	if created == nil {
		t.Errorf("expected part to be created")
	}
	if created.IsDuplicate {
		t.Errorf("expected part not to be marked as duplicate")
	}
}

func TestCatalogService_CreatePart_WithDuplicate(t *testing.T) {
	mockMfrRepo := &MockManufacturerRepository{}
	mockCatRepo := &MockCategoryRepository{}
	mockPartRepo := &MockPartRepository{}
	mockAttrRepo := &MockAttributeRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCatalogService(mockMfrRepo, mockCatRepo, mockPartRepo, mockAttrRepo, mockEventBus)

	manufacturerID := uuid.New()
	originalPart := &models.LibraryPart{
		ID:            uuid.New(),
		PartNumber:    "PART-001",
		Name:          "Original Part",
		ManufacturerID: manufacturerID,
		Status:        "approved",
	}
	mockPartRepo.Create(originalPart)

	duplicatePart := &models.LibraryPart{
		ID:            uuid.New(),
		PartNumber:    "PART-001",
		Name:          "Duplicate Part",
		ManufacturerID: manufacturerID,
		Status:        "pending",
	}

	err := service.CreatePart(duplicatePart)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify duplicate part was created and marked
	created, _ := mockPartRepo.GetByID(duplicatePart.ID)
	if created == nil {
		t.Errorf("expected duplicate part to be created")
	}
	if !created.IsDuplicate {
		t.Errorf("expected part to be marked as duplicate")
	}
	if created.OriginalPartID == nil || *created.OriginalPartID != originalPart.ID {
		t.Errorf("expected OriginalPartID to be set to original part ID")
	}
}

func TestCatalogService_ApprovePart(t *testing.T) {
	mockMfrRepo := &MockManufacturerRepository{}
	mockCatRepo := &MockCategoryRepository{}
	mockPartRepo := &MockPartRepository{}
	mockAttrRepo := &MockAttributeRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCatalogService(mockMfrRepo, mockCatRepo, mockPartRepo, mockAttrRepo, mockEventBus)

	partID := uuid.New()
	approvedBy := uuid.New()
	manufacturerID := uuid.New()
	part := &models.LibraryPart{
		ID:            partID,
		PartNumber:    "PART-001",
		Name:          "Test Part",
		ManufacturerID: manufacturerID,
		Status:        "pending",
	}
	mockPartRepo.Create(part)

	err := service.ApprovePart(partID, approvedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify part was updated
	updated, _ := mockPartRepo.GetByID(partID)
	if updated.Status != "approved" {
		t.Errorf("expected status 'approved', got %s", updated.Status)
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
		if event.Type != events.EventCatalogPartApproved {
			t.Errorf("expected event type %s, got %s", events.EventCatalogPartApproved, event.Type)
		}
		if event.Payload["part_id"] != partID.String() {
			t.Errorf("expected part_id %s, got %v", partID.String(), event.Payload["part_id"])
		}
	}
}

func TestCatalogService_RejectPart(t *testing.T) {
	mockMfrRepo := &MockManufacturerRepository{}
	mockCatRepo := &MockCategoryRepository{}
	mockPartRepo := &MockPartRepository{}
	mockAttrRepo := &MockAttributeRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCatalogService(mockMfrRepo, mockCatRepo, mockPartRepo, mockAttrRepo, mockEventBus)

	partID := uuid.New()
	manufacturerID := uuid.New()
	part := &models.LibraryPart{
		ID:            partID,
		PartNumber:    "PART-001",
		Name:          "Test Part",
		ManufacturerID: manufacturerID,
		Status:        "pending",
	}
	mockPartRepo.Create(part)

	reason := "Does not meet quality standards"
	err := service.RejectPart(partID, reason)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify part was updated
	updated, _ := mockPartRepo.GetByID(partID)
	if updated.Status != "rejected" {
		t.Errorf("expected status 'rejected', got %s", updated.Status)
	}
	if updated.RejectedReason != reason {
		t.Errorf("expected rejected reason '%s', got %s", reason, updated.RejectedReason)
	}
}

func TestCatalogService_GetPendingParts(t *testing.T) {
	mockMfrRepo := &MockManufacturerRepository{}
	mockCatRepo := &MockCategoryRepository{}
	mockPartRepo := &MockPartRepository{}
	mockAttrRepo := &MockAttributeRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCatalogService(mockMfrRepo, mockCatRepo, mockPartRepo, mockAttrRepo, mockEventBus)

	manufacturerID := uuid.New()
	pendingPart := &models.LibraryPart{
		ID:            uuid.New(),
		PartNumber:    "PART-001",
		ManufacturerID: manufacturerID,
		Status:        "pending",
	}
	approvedPart := &models.LibraryPart{
		ID:            uuid.New(),
		PartNumber:    "PART-002",
		ManufacturerID: manufacturerID,
		Status:        "approved",
	}
	mockPartRepo.Create(pendingPart)
	mockPartRepo.Create(approvedPart)

	pending, err := service.GetPendingParts()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(pending) != 1 {
		t.Errorf("expected 1 pending part, got %d", len(pending))
	}
	if pending[0].ID != pendingPart.ID {
		t.Errorf("expected pending part ID %s, got %s", pendingPart.ID, pending[0].ID)
	}
}
