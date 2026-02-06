package repository

import (
	"github.com/b2b-platform/catalog-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) *AttributeRepository {
	return &AttributeRepository{db: db}
}

func (r *AttributeRepository) Create(attribute *models.Attribute) error {
	return r.db.Create(attribute).Error
}

func (r *AttributeRepository) GetByID(id uuid.UUID) (*models.Attribute, error) {
	var attribute models.Attribute
	err := r.db.Where("id = ?", id).First(&attribute).Error
	return &attribute, err
}

func (r *AttributeRepository) List() ([]models.Attribute, error) {
	var attributes []models.Attribute
	err := r.db.Find(&attributes).Error
	return attributes, err
}

func (r *AttributeRepository) AddPartAttribute(partAttribute *models.PartAttribute) error {
	return r.db.Create(partAttribute).Error
}

func (r *AttributeRepository) GetPartAttributes(partID uuid.UUID) ([]models.PartAttribute, error) {
	var partAttributes []models.PartAttribute
	err := r.db.Preload("Attribute").Where("part_id = ?", partID).Find(&partAttributes).Error
	return partAttributes, err
}
