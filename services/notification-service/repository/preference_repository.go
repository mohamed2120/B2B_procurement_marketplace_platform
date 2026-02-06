package repository

import (
	"github.com/b2b-platform/notification-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) Create(preference *models.NotificationPreference) error {
	return r.db.Create(preference).Error
}

func (r *PreferenceRepository) GetByID(id uuid.UUID) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference
	err := r.db.Where("id = ?", id).First(&preference).Error
	return &preference, err
}

func (r *PreferenceRepository) GetUserPreferences(userID, tenantID uuid.UUID) ([]models.NotificationPreference, error) {
	var preferences []models.NotificationPreference
	err := r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&preferences).Error
	return preferences, err
}

func (r *PreferenceRepository) Update(preference *models.NotificationPreference) error {
	return r.db.Save(preference).Error
}

func (r *PreferenceRepository) IsEnabled(userID, tenantID uuid.UUID, channel, eventType string) (bool, error) {
	var preference models.NotificationPreference
	err := r.db.Where("user_id = ? AND tenant_id = ? AND channel = ? AND event_type = ?", 
		userID, tenantID, channel, eventType).First(&preference).Error
	
	if err == gorm.ErrRecordNotFound {
		// Default to enabled if preference doesn't exist
		return true, nil
	}
	
	return preference.IsEnabled, err
}
