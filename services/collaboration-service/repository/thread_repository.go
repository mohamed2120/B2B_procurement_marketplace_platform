package repository

import (
	"github.com/b2b-platform/collaboration-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadRepository struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(thread *models.ChatThread) error {
	return r.db.Create(thread).Error
}

func (r *ThreadRepository) GetByID(id uuid.UUID) (*models.ChatThread, error) {
	var thread models.ChatThread
	err := r.db.Preload("Participants").Preload("Messages").
		Where("id = ?", id).First(&thread).Error
	return &thread, err
}

func (r *ThreadRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.ChatThread, error) {
	var threads []models.ChatThread
	query := r.db.Preload("Participants").Where("tenant_id = ?", tenantID)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Order("updated_at DESC").Find(&threads).Error
	return threads, err
}

func (r *ThreadRepository) AddParticipant(participant *models.ThreadParticipant) error {
	return r.db.Create(participant).Error
}

func (r *ThreadRepository) GetUserThreads(userID, tenantID uuid.UUID) ([]models.ChatThread, error) {
	var threads []models.ChatThread
	err := r.db.Table("collaboration.chat_threads").
		Joins("INNER JOIN collaboration.thread_participants ON collaboration.chat_threads.id = collaboration.thread_participants.thread_id").
		Where("collaboration.thread_participants.user_id = ? AND collaboration.chat_threads.tenant_id = ?", userID, tenantID).
		Preload("Participants").
		Order("collaboration.chat_threads.updated_at DESC").
		Find(&threads).Error
	return threads, err
}
