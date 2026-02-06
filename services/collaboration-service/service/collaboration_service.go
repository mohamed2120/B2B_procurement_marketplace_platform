package service

import (
	"github.com/b2b-platform/collaboration-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type ThreadRepository interface {
	Create(thread *models.ChatThread) error
	GetByID(id uuid.UUID) (*models.ChatThread, error)
	List(tenantID uuid.UUID, limit, offset int) ([]models.ChatThread, error)
	GetUserThreads(userID, tenantID uuid.UUID) ([]models.ChatThread, error)
	AddParticipant(participant *models.ThreadParticipant) error
}

type MessageRepository interface {
	Create(message *models.ChatMessage) error
	GetByThread(threadID uuid.UUID, limit, offset int) ([]models.ChatMessage, error)
}

type DisputeRepository interface {
	Create(dispute *models.Dispute) error
	GetByID(id uuid.UUID) (*models.Dispute, error)
	List(tenantID uuid.UUID, status string) ([]models.Dispute, error)
	Resolve(disputeID, resolvedBy uuid.UUID, resolution string) error
}

type RatingRepository interface {
	Create(rating *models.Rating) error
	GetByEntity(entityType string, entityID uuid.UUID) ([]models.Rating, error)
	GetAverageRating(entityType string, entityID uuid.UUID) (float64, error)
	Moderate(ratingID, moderatedBy uuid.UUID) error
}

type CollaborationService struct {
	threadRepo  ThreadRepository
	messageRepo MessageRepository
	disputeRepo DisputeRepository
	ratingRepo  RatingRepository
	eventBus    events.EventBus
}

func NewCollaborationService(
	threadRepo ThreadRepository,
	messageRepo MessageRepository,
	disputeRepo DisputeRepository,
	ratingRepo RatingRepository,
	eventBus events.EventBus,
) *CollaborationService {
	return &CollaborationService{
		threadRepo:  threadRepo,
		messageRepo: messageRepo,
		disputeRepo: disputeRepo,
		ratingRepo:  ratingRepo,
		eventBus:    eventBus,
	}
}

func (s *CollaborationService) CreateThread(thread *models.ChatThread) error {
	return s.threadRepo.Create(thread)
}

func (s *CollaborationService) GetThread(id uuid.UUID) (*models.ChatThread, error) {
	return s.threadRepo.GetByID(id)
}

func (s *CollaborationService) ListThreads(tenantID uuid.UUID, limit, offset int) ([]models.ChatThread, error) {
	return s.threadRepo.List(tenantID, limit, offset)
}

func (s *CollaborationService) GetUserThreads(userID, tenantID uuid.UUID) ([]models.ChatThread, error) {
	return s.threadRepo.GetUserThreads(userID, tenantID)
}

func (s *CollaborationService) AddParticipant(participant *models.ThreadParticipant) error {
	return s.threadRepo.AddParticipant(participant)
}

func (s *CollaborationService) SendMessage(message *models.ChatMessage) error {
	if err := s.messageRepo.Create(message); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventChatMessageSent,
		"collaboration-service",
		map[string]interface{}{
			"message_id": message.ID.String(),
			"thread_id":  message.ThreadID.String(),
			"sender_id":  message.SenderID.String(),
		},
	)

	return s.eventBus.Publish(nil, event)
}

func (s *CollaborationService) GetThreadMessages(threadID uuid.UUID, limit, offset int) ([]models.ChatMessage, error) {
	return s.messageRepo.GetByThread(threadID, limit, offset)
}

func (s *CollaborationService) CreateDispute(dispute *models.Dispute) error {
	return s.disputeRepo.Create(dispute)
}

func (s *CollaborationService) GetDispute(id uuid.UUID) (*models.Dispute, error) {
	return s.disputeRepo.GetByID(id)
}

func (s *CollaborationService) ListDisputes(tenantID uuid.UUID, status string) ([]models.Dispute, error) {
	return s.disputeRepo.List(tenantID, status)
}

func (s *CollaborationService) ResolveDispute(disputeID, resolvedBy uuid.UUID, resolution string) error {
	return s.disputeRepo.Resolve(disputeID, resolvedBy, resolution)
}

func (s *CollaborationService) CreateRating(rating *models.Rating) error {
	return s.ratingRepo.Create(rating)
}

func (s *CollaborationService) GetRatings(entityType string, entityID uuid.UUID) ([]models.Rating, error) {
	return s.ratingRepo.GetByEntity(entityType, entityID)
}

func (s *CollaborationService) GetAverageRating(entityType string, entityID uuid.UUID) (float64, error) {
	return s.ratingRepo.GetAverageRating(entityType, entityID)
}

func (s *CollaborationService) ModerateRating(ratingID, moderatedBy uuid.UUID) error {
	return s.ratingRepo.Moderate(ratingID, moderatedBy)
}
