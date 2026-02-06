package repository

import (
	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmergencyRepository struct {
	db *gorm.DB
}

func NewEmergencyRepository(db *gorm.DB) *EmergencyRepository {
	return &EmergencyRepository{db: db}
}

func (r *EmergencyRepository) Create(sourcing *models.EmergencySourcing) error {
	return r.db.Create(sourcing).Error
}

func (r *EmergencyRepository) GetByID(id uuid.UUID) (*models.EmergencySourcing, error) {
	var sourcing models.EmergencySourcing
	err := r.db.Where("id = ?", id).First(&sourcing).Error
	return &sourcing, err
}

func (r *EmergencyRepository) List(tenantID uuid.UUID, status string) ([]models.EmergencySourcing, error) {
	var sourcing []models.EmergencySourcing
	query := r.db.Where("tenant_id = ?", tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("priority DESC, created_at DESC").Find(&sourcing).Error
	return sourcing, err
}

func (r *EmergencyRepository) Update(sourcing *models.EmergencySourcing) error {
	return r.db.Save(sourcing).Error
}

func (r *EmergencyRepository) Fulfill(sourcingID, fulfilledBy uuid.UUID) error {
	return r.db.Model(&models.EmergencySourcing{}).
		Where("id = ?", sourcingID).
		Updates(map[string]interface{}{
			"status":      "fulfilled",
			"fulfilled_by": fulfilledBy,
		}).Error
}
