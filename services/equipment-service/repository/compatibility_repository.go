package repository

import (
	"github.com/b2b-platform/equipment-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompatibilityRepository struct {
	db *gorm.DB
}

func NewCompatibilityRepository(db *gorm.DB) *CompatibilityRepository {
	return &CompatibilityRepository{db: db}
}

func (r *CompatibilityRepository) Create(mapping *models.CompatibilityMapping) error {
	return r.db.Create(mapping).Error
}

func (r *CompatibilityRepository) GetByID(id uuid.UUID) (*models.CompatibilityMapping, error) {
	var mapping models.CompatibilityMapping
	err := r.db.Preload("Equipment").Where("id = ?", id).First(&mapping).Error
	return &mapping, err
}

func (r *CompatibilityRepository) GetByEquipment(equipmentID uuid.UUID) ([]models.CompatibilityMapping, error) {
	var mappings []models.CompatibilityMapping
	err := r.db.Preload("Equipment").Where("equipment_id = ?", equipmentID).Find(&mappings).Error
	return mappings, err
}

func (r *CompatibilityRepository) GetByPart(partID uuid.UUID) ([]models.CompatibilityMapping, error) {
	var mappings []models.CompatibilityMapping
	err := r.db.Preload("Equipment").Where("part_id = ?", partID).Find(&mappings).Error
	return mappings, err
}

func (r *CompatibilityRepository) VerifyCompatibility(mappingID, verifiedBy uuid.UUID) error {
	return r.db.Model(&models.CompatibilityMapping{}).
		Where("id = ?", mappingID).
		Updates(map[string]interface{}{
			"is_compatible": true,
			"verified_by":   verifiedBy,
		}).Error
}

func (r *CompatibilityRepository) CheckCompatibility(equipmentID, partID uuid.UUID) (*models.CompatibilityMapping, error) {
	var mapping models.CompatibilityMapping
	err := r.db.Where("equipment_id = ? AND part_id = ?", equipmentID, partID).First(&mapping).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mapping, err
}
