package repository

import (
	"github.com/b2b-platform/catalog-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ManufacturerRepository struct {
	db *gorm.DB
}

func NewManufacturerRepository(db *gorm.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

func (r *ManufacturerRepository) Create(manufacturer *models.Manufacturer) error {
	return r.db.Create(manufacturer).Error
}

func (r *ManufacturerRepository) GetByID(id uuid.UUID) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	err := r.db.Where("id = ?", id).First(&manufacturer).Error
	return &manufacturer, err
}

func (r *ManufacturerRepository) List() ([]models.Manufacturer, error) {
	var manufacturers []models.Manufacturer
	err := r.db.Where("is_active = ?", true).Find(&manufacturers).Error
	return manufacturers, err
}

func (r *ManufacturerRepository) GetByCode(code string) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	err := r.db.Where("code = ?", code).First(&manufacturer).Error
	return &manufacturer, err
}
