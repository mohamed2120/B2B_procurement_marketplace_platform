package repository

import (
	"github.com/b2b-platform/procurement-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RFQRepository struct {
	db *gorm.DB
}

func NewRFQRepository(db *gorm.DB) *RFQRepository {
	return &RFQRepository{db: db}
}

func (r *RFQRepository) Create(rfq *models.RFQ) error {
	return r.db.Create(rfq).Error
}

func (r *RFQRepository) GetByID(id uuid.UUID) (*models.RFQ, error) {
	var rfq models.RFQ
	err := r.db.Preload("PR").Preload("Quotes").Where("id = ?", id).First(&rfq).Error
	return &rfq, err
}

func (r *RFQRepository) List(tenantID uuid.UUID) ([]models.RFQ, error) {
	var rfqs []models.RFQ
	err := r.db.Where("tenant_id = ?", tenantID).Find(&rfqs).Error
	return rfqs, err
}
