package service

import (
	"testing"

	"github.com/b2b-platform/collaboration-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// MockThreadRepository for testing
type MockThreadRepository struct {
	threads      map[uuid.UUID]*models.ChatThread
	participants []models.ThreadParticipant
}

func (m *MockThreadRepository) Create(thread *models.ChatThread) error {
	if m.threads == nil {
		m.threads = make(map[uuid.UUID]*models.ChatThread)
	}
	m.threads[thread.ID] = thread
	return nil
}

func (m *MockThreadRepository) GetByID(id uuid.UUID) (*models.ChatThread, error) {
	if thread, ok := m.threads[id]; ok {
		return thread, nil
	}
	return nil, nil
}

func (m *MockThreadRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.ChatThread, error) {
	var result []models.ChatThread
	for _, thread := range m.threads {
		if thread.TenantID == tenantID {
			result = append(result, *thread)
		}
	}
	return result, nil
}

func (m *MockThreadRepository) GetUserThreads(userID, tenantID uuid.UUID) ([]models.ChatThread, error) {
	var result []models.ChatThread
	for _, thread := range m.threads {
		if thread.TenantID == tenantID {
			// Check if user is a participant
			for _, p := range m.participants {
				if p.ThreadID == thread.ID && p.UserID == userID {
					result = append(result, *thread)
					break
				}
			}
		}
	}
	return result, nil
}

func (m *MockThreadRepository) AddParticipant(participant *models.ThreadParticipant) error {
	if m.participants == nil {
		m.participants = make([]models.ThreadParticipant, 0)
	}
	m.participants = append(m.participants, *participant)
	return nil
}

// MockMessageRepository for testing
type MockMessageRepository struct {
	messages map[uuid.UUID]*models.ChatMessage
}

func (m *MockMessageRepository) Create(message *models.ChatMessage) error {
	if m.messages == nil {
		m.messages = make(map[uuid.UUID]*models.ChatMessage)
	}
	m.messages[message.ID] = message
	return nil
}

func (m *MockMessageRepository) GetByThread(threadID uuid.UUID, limit, offset int) ([]models.ChatMessage, error) {
	var result []models.ChatMessage
	for _, msg := range m.messages {
		if msg.ThreadID == threadID {
			result = append(result, *msg)
		}
	}
	return result, nil
}

// MockDisputeRepository for testing
type MockDisputeRepository struct {
	disputes map[uuid.UUID]*models.Dispute
}

func (m *MockDisputeRepository) Create(dispute *models.Dispute) error {
	if m.disputes == nil {
		m.disputes = make(map[uuid.UUID]*models.Dispute)
	}
	m.disputes[dispute.ID] = dispute
	return nil
}

func (m *MockDisputeRepository) GetByID(id uuid.UUID) (*models.Dispute, error) {
	if dispute, ok := m.disputes[id]; ok {
		return dispute, nil
	}
	return nil, nil
}

func (m *MockDisputeRepository) List(tenantID uuid.UUID, status string) ([]models.Dispute, error) {
	var result []models.Dispute
	for _, dispute := range m.disputes {
		if dispute.TenantID == tenantID {
			if status == "" || dispute.Status == status {
				result = append(result, *dispute)
			}
		}
	}
	return result, nil
}

func (m *MockDisputeRepository) Resolve(disputeID, resolvedBy uuid.UUID, resolution string) error {
	if dispute, ok := m.disputes[disputeID]; ok {
		dispute.Status = "resolved"
		dispute.Resolution = resolution
		dispute.ResolvedBy = &resolvedBy
	}
	return nil
}

// MockRatingRepository for testing
type MockRatingRepository struct {
	ratings map[uuid.UUID]*models.Rating
}

func (m *MockRatingRepository) Create(rating *models.Rating) error {
	if m.ratings == nil {
		m.ratings = make(map[uuid.UUID]*models.Rating)
	}
	m.ratings[rating.ID] = rating
	return nil
}

func (m *MockRatingRepository) GetByID(id uuid.UUID) (*models.Rating, error) {
	if rating, ok := m.ratings[id]; ok {
		return rating, nil
	}
	return nil, nil
}

func (m *MockRatingRepository) GetByEntity(entityType string, entityID uuid.UUID) ([]models.Rating, error) {
	var result []models.Rating
	for _, rating := range m.ratings {
		if rating.RatedEntityType == entityType && rating.RatedEntityID == entityID {
			result = append(result, *rating)
		}
	}
	return result, nil
}

func (m *MockRatingRepository) GetAverageRating(entityType string, entityID uuid.UUID) (float64, error) {
	var sum float64
	var count int
	for _, rating := range m.ratings {
		if rating.RatedEntityType == entityType && rating.RatedEntityID == entityID {
			sum += float64(rating.Rating)
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

func (m *MockRatingRepository) Moderate(ratingID, moderatedBy uuid.UUID) error {
	if rating, ok := m.ratings[ratingID]; ok {
		rating.IsModerated = true
		rating.ModeratedBy = &moderatedBy
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

func TestCollaborationService_SendMessage_PublishesEvent(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	threadID := uuid.New()
	senderID := uuid.New()
	message := &models.ChatMessage{
		ID:       uuid.New(),
		ThreadID: threadID,
		SenderID: senderID,
		Message:  "Test message",
	}

	err := service.SendMessage(message)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify message was created
	if len(mockMessageRepo.messages) == 0 {
		t.Errorf("expected message to be created")
	}

	// Verify event was published
	if len(mockEventBus.publishedEvents) == 0 {
		t.Errorf("expected event to be published")
	} else {
		event := mockEventBus.publishedEvents[0]
		if event.Type != events.EventChatMessageSent {
			t.Errorf("expected event type %s, got %s", events.EventChatMessageSent, event.Type)
		}
		if event.Payload["message_id"] != message.ID.String() {
			t.Errorf("expected message_id %s, got %v", message.ID.String(), event.Payload["message_id"])
		}
	}
}

func TestCollaborationService_CreateDispute(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	tenantID := uuid.New()
	dispute := &models.Dispute{
		ID:          uuid.New(),
		TenantID:    tenantID,
		OrderID:     uuid.New(),
		DisputeType: "quality",
		Status:      "open",
		Title:       "Order issue",
		Description: "Product not as described",
		RaisedBy:    uuid.New(),
	}

	err := service.CreateDispute(dispute)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify dispute was created
	created, _ := mockDisputeRepo.GetByID(dispute.ID)
	if created == nil {
		t.Errorf("expected dispute to be created")
	}
	if created.Description != "Product not as described" {
		t.Errorf("expected description 'Product not as described', got %s", created.Description)
	}
}

func TestCollaborationService_ResolveDispute(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	disputeID := uuid.New()
	resolvedBy := uuid.New()
	dispute := &models.Dispute{
		ID:     disputeID,
		Status: "open",
	}
	mockDisputeRepo.Create(dispute)

	resolution := "Refund issued"
	err := service.ResolveDispute(disputeID, resolvedBy, resolution)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify dispute was resolved
	updated, _ := mockDisputeRepo.GetByID(disputeID)
	if updated.Status != "resolved" {
		t.Errorf("expected status 'resolved', got %s", updated.Status)
	}
	if updated.Resolution != resolution {
		t.Errorf("expected resolution '%s', got %s", resolution, updated.Resolution)
	}
}

func TestCollaborationService_CreateRating(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	entityID := uuid.New()
	rating := &models.Rating{
		ID:              uuid.New(),
		RatedEntityType: "supplier",
		RatedEntityID:   entityID,
		Rating:          5,
		Comment:         "Great service!",
	}

	err := service.CreateRating(rating)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify rating was created
	ratings, _ := service.GetRatings("supplier", entityID)
	if len(ratings) != 1 {
		t.Errorf("expected 1 rating, got %d", len(ratings))
	}
	if ratings[0].Rating != 5 {
		t.Errorf("expected rating 5, got %d", ratings[0].Rating)
	}
}

func TestCollaborationService_GetAverageRating(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	entityID := uuid.New()
	rating1 := &models.Rating{
		ID:              uuid.New(),
		RatedEntityType: "supplier",
		RatedEntityID:   entityID,
		Rating:          4,
	}
	rating2 := &models.Rating{
		ID:              uuid.New(),
		RatedEntityType: "supplier",
		RatedEntityID:   entityID,
		Rating:          5,
	}
	mockRatingRepo.Create(rating1)
	mockRatingRepo.Create(rating2)

	average, err := service.GetAverageRating("supplier", entityID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if average != 4.5 {
		t.Errorf("expected average 4.5, got %.2f", average)
	}
}

func TestCollaborationService_ModerateRating(t *testing.T) {
	mockThreadRepo := &MockThreadRepository{}
	mockMessageRepo := &MockMessageRepository{}
	mockDisputeRepo := &MockDisputeRepository{}
	mockRatingRepo := &MockRatingRepository{}
	mockEventBus := &MockEventBus{}

	service := NewCollaborationService(mockThreadRepo, mockMessageRepo, mockDisputeRepo, mockRatingRepo, mockEventBus)

	ratingID := uuid.New()
	moderatedBy := uuid.New()
	rating := &models.Rating{
		ID:          ratingID,
		IsModerated: false,
	}
	mockRatingRepo.Create(rating)

	err := service.ModerateRating(ratingID, moderatedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify rating was moderated
	updated, _ := mockRatingRepo.GetByID(ratingID)
	if !updated.IsModerated {
		t.Errorf("expected rating to be moderated")
	}
	if updated.ModeratedBy == nil || *updated.ModeratedBy != moderatedBy {
		t.Errorf("expected ModeratedBy to be set")
	}
}
