package repository

import (
	"github.com/b2b-platform/marketplace-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) GetByID(id uuid.UUID) (*models.Store, error) {
	var store models.Store
	err := r.db.Preload("Listings").Where("id = ?", id).First(&store).Error
	return &store, err
}

func (r *StoreRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Store, error) {
	var stores []models.Store
	query := r.db.Where("tenant_id = ?", tenantID)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&stores).Error
	return stores, err
}

func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) GetByCompany(companyID uuid.UUID) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Where("company_id = ?", companyID).Find(&stores).Error
	return stores, err
}
