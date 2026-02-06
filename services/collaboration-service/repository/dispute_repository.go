package repository

import (
	"github.com/b2b-platform/collaboration-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DisputeRepository struct {
	db *gorm.DB
}

func NewDisputeRepository(db *gorm.DB) *DisputeRepository {
	return &DisputeRepository{db: db}
}

func (r *DisputeRepository) Create(dispute *models.Dispute) error {
	return r.db.Create(dispute).Error
}

func (r *DisputeRepository) GetByID(id uuid.UUID) (*models.Dispute, error) {
	var dispute models.Dispute
	err := r.db.Where("id = ?", id).First(&dispute).Error
	return &dispute, err
}

func (r *DisputeRepository) List(tenantID uuid.UUID, status string) ([]models.Dispute, error) {
	var disputes []models.Dispute
	query := r.db.Where("tenant_id = ?", tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&disputes).Error
	return disputes, err
}

func (r *DisputeRepository) Update(dispute *models.Dispute) error {
	return r.db.Save(dispute).Error
}

func (r *DisputeRepository) Resolve(disputeID, resolvedBy uuid.UUID, resolution string) error {
	return r.db.Model(&models.Dispute{}).
		Where("id = ?", disputeID).
		Updates(map[string]interface{}{
			"status":      "resolved",
			"resolved_by": resolvedBy,
			"resolution":  resolution,
		}).Error
}
