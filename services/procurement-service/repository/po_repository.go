package repository

import (
	"github.com/b2b-platform/procurement-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PORepository struct {
	db *gorm.DB
}

func NewPORepository(db *gorm.DB) *PORepository {
	return &PORepository{db: db}
}

func (r *PORepository) Create(po *models.PurchaseOrder) error {
	return r.db.Create(po).Error
}

func (r *PORepository) GetByID(id uuid.UUID) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	err := r.db.Preload("Items").Preload("PR").Preload("Quote").Where("id = ?", id).First(&po).Error
	return &po, err
}

func (r *PORepository) List(tenantID uuid.UUID) ([]models.PurchaseOrder, error) {
	var pos []models.PurchaseOrder
	err := r.db.Where("tenant_id = ?", tenantID).Find(&pos).Error
	return pos, err
}

func (r *PORepository) Update(po *models.PurchaseOrder) error {
	return r.db.Save(po).Error
}
