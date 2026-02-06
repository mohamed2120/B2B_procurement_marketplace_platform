package repository

import (
	"github.com/b2b-platform/collaboration-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) GetByID(id uuid.UUID) (*models.ChatMessage, error) {
	var message models.ChatMessage
	err := r.db.Preload("Files").Where("id = ?", id).First(&message).Error
	return &message, err
}

func (r *MessageRepository) GetByThread(threadID uuid.UUID, limit, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := r.db.Preload("Files").Where("thread_id = ?", threadID)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) Update(message *models.ChatMessage) error {
	return r.db.Save(message).Error
}

func (r *MessageRepository) MarkAsRead(messageID uuid.UUID) error {
	return r.db.Model(&models.ChatMessage{}).
		Where("id = ?", messageID).
		Update("is_read", true).Error
}

func (r *MessageRepository) AddFile(file *models.MessageFile) error {
	return r.db.Create(file).Error
}
