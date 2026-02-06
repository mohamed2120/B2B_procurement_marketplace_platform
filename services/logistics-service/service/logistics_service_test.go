package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/logistics-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockShipmentRepository for testing
type MockShipmentRepository struct {
	shipments     map[uuid.UUID]*models.Shipment
	trackingEvents []models.TrackingEvent
}

func (m *MockShipmentRepository) Create(shipment *models.Shipment) error {
	if m.shipments == nil {
		m.shipments = make(map[uuid.UUID]*models.Shipment)
	}
	m.shipments[shipment.ID] = shipment
	return nil
}

func (m *MockShipmentRepository) GetByID(id uuid.UUID) (*models.Shipment, error) {
	if shipment, ok := m.shipments[id]; ok {
		return shipment, nil
	}
	return nil, nil
}

func (m *MockShipmentRepository) List(tenantID uuid.UUID) ([]models.Shipment, error) {
	var result []models.Shipment
	for _, shipment := range m.shipments {
		if shipment.TenantID == tenantID {
			result = append(result, *shipment)
		}
	}
	return result, nil
}

func (m *MockShipmentRepository) Update(shipment *models.Shipment) error {
	m.shipments[shipment.ID] = shipment
	return nil
}

func (m *MockShipmentRepository) AddTrackingEvent(event *models.TrackingEvent) error {
	if m.trackingEvents == nil {
		m.trackingEvents = make([]models.TrackingEvent, 0)
	}
	m.trackingEvents = append(m.trackingEvents, *event)
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

func TestLogisticsService_Create(t *testing.T) {
	mockRepo := &MockShipmentRepository{}
	mockEventBus := &MockEventBus{}
	service := NewLogisticsService(mockRepo, mockEventBus)

	tenantID := uuid.New()
	shipment := &models.Shipment{
		ID:             uuid.New(),
		TenantID:      tenantID,
		TrackingNumber: "TRACK-001",
		Status:        "in_transit",
		ETA:           time.Now().Add(24 * time.Hour),
	}

	err := service.Create(shipment)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify shipment was created
	created, _ := mockRepo.GetByID(shipment.ID)
	if created == nil {
		t.Errorf("expected shipment to be created")
	}
	if created.TrackingNumber != "TRACK-001" {
		t.Errorf("expected tracking number 'TRACK-001', got %s", created.TrackingNumber)
	}
}

func TestLogisticsService_UpdateTracking_NotLate(t *testing.T) {
	mockRepo := &MockShipmentRepository{}
	mockEventBus := &MockEventBus{}
	service := NewLogisticsService(mockRepo, mockEventBus)

	tenantID := uuid.New()
	shipment := &models.Shipment{
		ID:             uuid.New(),
		TenantID:      tenantID,
		TrackingNumber: "TRACK-001",
		Status:        "in_transit",
		ETA:           time.Now().Add(24 * time.Hour), // Future ETA
		IsLate:        false,
	}
	mockRepo.Create(shipment)

	event := &models.TrackingEvent{
		ID:             uuid.New(),
		ShipmentID:    shipment.ID,
		Status:        "in_transit",
		Location:      "Warehouse A",
		Timestamp:     time.Now(),
	}

	err := service.UpdateTracking(shipment.ID, event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify tracking event was added
	if len(mockRepo.trackingEvents) == 0 {
		t.Errorf("expected tracking event to be added")
	}

	// Verify no late event was published (ETA is in the future)
	if len(mockEventBus.publishedEvents) > 0 {
		t.Errorf("expected no late event to be published")
	}
}

func TestLogisticsService_UpdateTracking_Late(t *testing.T) {
	mockRepo := &MockShipmentRepository{}
	mockEventBus := &MockEventBus{}
	service := NewLogisticsService(mockRepo, mockEventBus)

	tenantID := uuid.New()
	shipment := &models.Shipment{
		ID:             uuid.New(),
		TenantID:      tenantID,
		TrackingNumber: "TRACK-001",
		Status:        "in_transit",
		ETA:           time.Now().Add(-24 * time.Hour), // Past ETA
		IsLate:        false,
	}
	mockRepo.Create(shipment)

	event := &models.TrackingEvent{
		ID:          uuid.New(),
		ShipmentID: shipment.ID,
		Status:     "in_transit",
		Location:   "Warehouse A",
		Timestamp:  time.Now(),
	}

	err := service.UpdateTracking(shipment.ID, event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify shipment was marked as late
	updated, _ := mockRepo.GetByID(shipment.ID)
	if !updated.IsLate {
		t.Errorf("expected shipment to be marked as late")
	}

	// Verify late event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected late event to be published")
	} else {
		lateEvent := mockEventBus.publishedEvents[0]
		if lateEvent.Type != events.EventShipmentLate {
			t.Errorf("expected event type %s, got %s", events.EventShipmentLate, lateEvent.Type)
		}
		if lateEvent.Payload["shipment_id"] != shipment.ID.String() {
			t.Errorf("expected shipment_id %s, got %v", shipment.ID.String(), lateEvent.Payload["shipment_id"])
		}
	}
}

func TestLogisticsService_UpdateTracking_AlreadyLate(t *testing.T) {
	mockRepo := &MockShipmentRepository{}
	mockEventBus := &MockEventBus{}
	service := NewLogisticsService(mockRepo, mockEventBus)

	tenantID := uuid.New()
	shipment := &models.Shipment{
		ID:             uuid.New(),
		TenantID:      tenantID,
		TrackingNumber: "TRACK-001",
		Status:        "in_transit",
		ETA:           time.Now().Add(-24 * time.Hour), // Past ETA
		IsLate:        true, // Already marked as late
	}
	mockRepo.Create(shipment)

	event := &models.TrackingEvent{
		ID:          uuid.New(),
		ShipmentID: shipment.ID,
		Status:     "in_transit",
		Location:   "Warehouse A",
		Timestamp:  time.Now(),
	}

	err := service.UpdateTracking(shipment.ID, event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify no additional late event was published (already late)
	if len(mockEventBus.publishedEvents) > 0 {
		t.Errorf("expected no additional late event to be published (already late)")
	}
}
