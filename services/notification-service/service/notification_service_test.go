package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/notification-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockEventBus is a simple mock implementation of EventBus for testing
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

// MockRepository is a mock repository for testing
type MockNotificationRepository struct {
	notifications []models.Notification
}

func (m *MockNotificationRepository) Create(notification *models.Notification) error {
	m.notifications = append(m.notifications, *notification)
	return nil
}

func (m *MockNotificationRepository) GetByID(id uuid.UUID) (*models.Notification, error) {
	for i := range m.notifications {
		if m.notifications[i].ID == id {
			return &m.notifications[i], nil
		}
	}
	return nil, nil
}

func (m *MockNotificationRepository) GetUserNotifications(userID, tenantID uuid.UUID, limit, offset int, unreadOnly bool) ([]models.Notification, error) {
	var result []models.Notification
	for i := range m.notifications {
		if m.notifications[i].UserID == userID && m.notifications[i].TenantID == tenantID {
			result = append(result, m.notifications[i])
		}
	}
	return result, nil
}

func (m *MockNotificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	for i := range m.notifications {
		if m.notifications[i].ID == notificationID {
			m.notifications[i].IsRead = true
			now := time.Now()
			m.notifications[i].ReadAt = &now
			return nil
		}
	}
	return nil
}

func (m *MockNotificationRepository) MarkAllAsRead(userID, tenantID uuid.UUID) error {
	for i := range m.notifications {
		if m.notifications[i].UserID == userID && m.notifications[i].TenantID == tenantID {
			m.notifications[i].IsRead = true
			now := time.Now()
			m.notifications[i].ReadAt = &now
		}
	}
	return nil
}

func (m *MockNotificationRepository) GetPending() ([]models.Notification, error) {
	var result []models.Notification
	for i := range m.notifications {
		if m.notifications[i].Status == "pending" {
			result = append(result, m.notifications[i])
		}
	}
	return result, nil
}

func (m *MockNotificationRepository) UpdateStatus(notificationID uuid.UUID, status string) error {
	for i := range m.notifications {
		if m.notifications[i].ID == notificationID {
			m.notifications[i].Status = status
			if status == "sent" {
				now := time.Now()
				m.notifications[i].SentAt = &now
			}
			return nil
		}
	}
	return nil
}

type MockTemplateRepository struct {
	templates map[string]*models.NotificationTemplate
}

func (m *MockTemplateRepository) GetByCode(code string) (*models.NotificationTemplate, error) {
	if template, ok := m.templates[code]; ok {
		return template, nil
	}
	return nil, nil
}

type MockPreferenceRepository struct {
	preferences map[string]bool
}

func (m *MockPreferenceRepository) IsEnabled(userID, tenantID uuid.UUID, channel, eventType string) (bool, error) {
	key := channel + ":" + eventType
	if enabled, ok := m.preferences[key]; ok {
		return enabled, nil
	}
	return true, nil // Default to enabled
}

func TestEventConsumer_HandleEvent(t *testing.T) {
	mockRepo := &MockNotificationRepository{}
	mockTemplateRepo := &MockTemplateRepository{
		templates: map[string]*models.NotificationTemplate{
			"order_placed": {
				Code:      "order_placed",
				Subject:   "New Order",
				Body:      "Order {{order_number}} placed",
				Channel:   "in_app",
				EventType: "procurement.order.placed.v1",
			},
		},
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: make(map[string]bool),
	}

	notificationService := &NotificationService{
		templateRepo:     mockTemplateRepo,
		preferenceRepo:   mockPrefRepo,
		notificationRepo: mockRepo,
	}

	consumer := NewEventConsumer(notificationService)

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	tests := []struct {
		name    string
		event   *events.EventEnvelope
		wantErr bool
	}{
		{
			name: "Order Placed Event",
			event: events.NewEventEnvelope(
				events.EventOrderPlaced,
				"procurement-service",
				map[string]interface{}{
					"po_id":     "po-123",
					"po_number": "PO-001",
				},
			).WithTenantID(tenantID),
			wantErr: false,
		},
		{
			name: "PR Approved Event",
			event: events.NewEventEnvelope(
				events.EventPRApproved,
				"procurement-service",
				map[string]interface{}{
					"pr_id":     "pr-123",
					"pr_number": "PR-001",
				},
			).WithTenantID(tenantID),
			wantErr: false,
		},
		{
			name: "Shipment Late Event",
			event: events.NewEventEnvelope(
				events.EventShipmentLate,
				"logistics-service",
				map[string]interface{}{
					"shipment_id":    "ship-123",
					"tracking_number": "TRACK-001",
					"eta":            time.Now(),
				},
			).WithTenantID(tenantID),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := consumer.HandleEvent(tt.event)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Verify notification was created
				if len(mockRepo.notifications) == 0 {
					t.Errorf("expected notification to be created")
				}
			}
		})
	}
}
