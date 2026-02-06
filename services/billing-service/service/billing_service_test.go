package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockPlanRepository for testing
type MockPlanRepository struct {
	plans map[uuid.UUID]*models.Plan
}

func (m *MockPlanRepository) Create(plan *models.Plan) error {
	if m.plans == nil {
		m.plans = make(map[uuid.UUID]*models.Plan)
	}
	m.plans[plan.ID] = plan
	return nil
}

func (m *MockPlanRepository) GetByID(id uuid.UUID) (*models.Plan, error) {
	if plan, ok := m.plans[id]; ok {
		return plan, nil
	}
	return nil, nil
}

func (m *MockPlanRepository) GetByCode(code string) (*models.Plan, error) {
	for _, plan := range m.plans {
		if plan.Code == code {
			return plan, nil
		}
	}
	return nil, nil
}

func (m *MockPlanRepository) List() ([]models.Plan, error) {
	var result []models.Plan
	for _, plan := range m.plans {
		result = append(result, *plan)
	}
	return result, nil
}

// MockSubscriptionRepository for testing
type MockSubscriptionRepository struct {
	subscriptions map[uuid.UUID]*models.Subscription
}

func (m *MockSubscriptionRepository) Create(subscription *models.Subscription) error {
	if m.subscriptions == nil {
		m.subscriptions = make(map[uuid.UUID]*models.Subscription)
	}
	m.subscriptions[subscription.ID] = subscription
	return nil
}

func (m *MockSubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	if sub, ok := m.subscriptions[id]; ok {
		return sub, nil
	}
	return nil, nil
}

func (m *MockSubscriptionRepository) GetByTenant(tenantID uuid.UUID) (*models.Subscription, error) {
	for _, sub := range m.subscriptions {
		if sub.TenantID == tenantID {
			return sub, nil
		}
	}
	return nil, nil
}

func (m *MockSubscriptionRepository) Cancel(subscriptionID uuid.UUID) error {
	if sub, ok := m.subscriptions[subscriptionID]; ok {
		now := time.Now()
		sub.Status = "cancelled"
		sub.CancelledAt = &now
	}
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

func (m *MockEventBus) SubscribeAll(ctx interface{}, handler func(*events.EventEnvelope) error) error {
	return nil
}

func TestBillingService_CreatePlan(t *testing.T) {
	mockPlanRepo := &MockPlanRepository{}
	mockSubRepo := &MockSubscriptionRepository{}
	mockEventBus := &MockEventBus{}

	service := NewBillingService(mockPlanRepo, mockSubRepo, mockEventBus)

	plan := &models.Plan{
		ID:   uuid.New(),
		Code: "basic",
		Name: "Basic Plan",
		Price: 99.99,
		Entitlements: []models.Entitlement{
			{
				Feature: "users",
				Limit:   10,
			},
		},
	}

	err := service.CreatePlan(plan)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify plan was created
	created, _ := service.GetPlan(plan.ID)
	if created == nil {
		t.Errorf("expected plan to be created")
	}
	if created.Code != "basic" {
		t.Errorf("expected code 'basic', got %s", created.Code)
	}
}

func TestBillingService_CreateSubscription_PublishesEvent(t *testing.T) {
	mockPlanRepo := &MockPlanRepository{}
	mockSubRepo := &MockSubscriptionRepository{}
	mockEventBus := &MockEventBus{}

	service := NewBillingService(mockPlanRepo, mockSubRepo, mockEventBus)

	planID := uuid.New()
	tenantID := uuid.New()
	subscription := &models.Subscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		PlanID:   planID,
	}

	err := service.CreateSubscription(subscription)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify subscription was created
	created, _ := service.GetSubscription(subscription.ID)
	if created == nil {
		t.Errorf("expected subscription to be created")
	}
	if created.Status != "active" {
		t.Errorf("expected status 'active', got %s", created.Status)
	}
	if created.StartedAt.IsZero() {
		t.Errorf("expected StartedAt to be set")
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventSubscriptionStarted {
			t.Errorf("expected event type %s, got %s", events.EventSubscriptionStarted, event.Type)
		}
		if event.Payload["subscription_id"] != subscription.ID.String() {
			t.Errorf("expected subscription_id %s, got %v", subscription.ID.String(), event.Payload["subscription_id"])
		}
	}
}

func TestBillingService_CancelSubscription(t *testing.T) {
	mockPlanRepo := &MockPlanRepository{}
	mockSubRepo := &MockSubscriptionRepository{}
	mockEventBus := &MockEventBus{}

	service := NewBillingService(mockPlanRepo, mockSubRepo, mockEventBus)

	subscriptionID := uuid.New()
	subscription := &models.Subscription{
		ID:     subscriptionID,
		Status: "active",
	}
	mockSubRepo.Create(subscription)

	err := service.CancelSubscription(subscriptionID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify subscription was cancelled
	updated, _ := service.GetSubscription(subscriptionID)
	if updated.Status != "cancelled" {
		t.Errorf("expected status 'cancelled', got %s", updated.Status)
	}
	if updated.CancelledAt == nil {
		t.Errorf("expected CancelledAt to be set")
	}
}

func TestBillingService_CheckEntitlement(t *testing.T) {
	mockPlanRepo := &MockPlanRepository{}
	mockSubRepo := &MockSubscriptionRepository{}
	mockEventBus := &MockEventBus{}

	service := NewBillingService(mockPlanRepo, mockSubRepo, mockEventBus)

	planID := uuid.New()
	tenantID := uuid.New()
	plan := &models.Plan{
		ID:   planID,
		Code: "basic",
		Name: "Basic Plan",
		Entitlements: []models.Entitlement{
			{
				Feature: "users",
				Limit:   10,
			},
			{
				Feature: "storage",
				Limit:   100,
			},
		},
	}
	mockPlanRepo.Create(plan)

	subscription := &models.Subscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		PlanID:   planID,
		Status:   "active",
	}
	mockSubRepo.Create(subscription)

	// Test existing entitlement
	enabled, limit, err := service.CheckEntitlement(tenantID, "users")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !enabled {
		t.Errorf("expected entitlement 'users' to be enabled")
	}
	if limit != 10 {
		t.Errorf("expected limit 10, got %d", limit)
	}

	// Test non-existing entitlement
	enabled, limit, err = service.CheckEntitlement(tenantID, "nonexistent")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if enabled {
		t.Errorf("expected entitlement 'nonexistent' to be disabled")
	}
	if limit != 0 {
		t.Errorf("expected limit 0, got %d", limit)
	}
}

func TestBillingService_GetPlanByCode(t *testing.T) {
	mockPlanRepo := &MockPlanRepository{}
	mockSubRepo := &MockSubscriptionRepository{}
	mockEventBus := &MockEventBus{}

	service := NewBillingService(mockPlanRepo, mockSubRepo, mockEventBus)

	plan := &models.Plan{
		ID:   uuid.New(),
		Code: "premium",
		Name: "Premium Plan",
	}
	mockPlanRepo.Create(plan)

	retrieved, err := service.GetPlanByCode("premium")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Errorf("expected plan to be found")
	}
	if retrieved.Code != "premium" {
		t.Errorf("expected code 'premium', got %s", retrieved.Code)
	}
}
