package repository

import (
	"github.com/b2b-platform/catalog-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartRepository struct {
	db *gorm.DB
}

func NewPartRepository(db *gorm.DB) *PartRepository {
	return &PartRepository{db: db}
}

func (r *PartRepository) Create(part *models.LibraryPart) error {
	return r.db.Create(part).Error
}

func (r *PartRepository) GetByID(id uuid.UUID) (*models.LibraryPart, error) {
	var part models.LibraryPart
	err := r.db.Preload("Manufacturer").Preload("Category").Preload("PartAttributes.Attribute").
		Where("id = ?", id).First(&part).Error
	return &part, err
}

func (r *PartRepository) List(limit, offset int, status string) ([]models.LibraryPart, error) {
	var parts []models.LibraryPart
	query := r.db.Preload("Manufacturer").Preload("Category")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	err := query.Find(&parts).Error
	return parts, err
}

func (r *PartRepository) Update(part *models.LibraryPart) error {
	return r.db.Save(part).Error
}

func (r *PartRepository) FindDuplicate(partNumber string, manufacturerID uuid.UUID) (*models.LibraryPart, error) {
	var part models.LibraryPart
	err := r.db.Where("part_number = ? AND manufacturer_id = ? AND is_duplicate = ?", 
		partNumber, manufacturerID, false).First(&part).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &part, err
}

func (r *PartRepository) MarkAsDuplicate(partID, duplicateOf uuid.UUID) error {
	return r.db.Model(&models.LibraryPart{}).
		Where("id = ?", partID).
		Updates(map[string]interface{}{
			"is_duplicate": true,
			"duplicate_of": duplicateOf,
		}).Error
}

func (r *PartRepository) GetPendingApproval() ([]models.LibraryPart, error) {
	var parts []models.LibraryPart
	err := r.db.Preload("Manufacturer").Preload("Category").
		Where("status = ?", "pending").Find(&parts).Error
	return parts, err
}
