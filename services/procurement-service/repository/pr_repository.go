package repository

import (
	"github.com/b2b-platform/procurement-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PRRepository struct {
	db *gorm.DB
}

func NewPRRepository(db *gorm.DB) *PRRepository {
	return &PRRepository{db: db}
}

func (r *PRRepository) Create(pr *models.PurchaseRequest) error {
	return r.db.Create(pr).Error
}

func (r *PRRepository) GetByID(id uuid.UUID) (*models.PurchaseRequest, error) {
	var pr models.PurchaseRequest
	err := r.db.Preload("Items").Preload("Approvals").Where("id = ?", id).First(&pr).Error
	return &pr, err
}

func (r *PRRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.PurchaseRequest, error) {
	var prs []models.PurchaseRequest
	query := r.db.Where("tenant_id = ?", tenantID)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&prs).Error
	return prs, err
}

func (r *PRRepository) Update(pr *models.PurchaseRequest) error {
	return r.db.Save(pr).Error
}

func (r *PRRepository) AddApproval(approval *models.PRApproval) error {
	return r.db.Create(approval).Error
}
