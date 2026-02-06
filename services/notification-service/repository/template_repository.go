package repository

import (
	"github.com/b2b-platform/notification-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(template *models.NotificationTemplate) error {
	return r.db.Create(template).Error
}

func (r *TemplateRepository) GetByID(id uuid.UUID) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.Where("id = ?", id).First(&template).Error
	return &template, err
}

func (r *TemplateRepository) GetByCode(code string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&template).Error
	return &template, err
}

func (r *TemplateRepository) List() ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	err := r.db.Where("is_active = ?", true).Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) GetByEventType(eventType string) ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	err := r.db.Where("event_type = ? AND is_active = ?", eventType, true).Find(&templates).Error
	return templates, err
}
