package repository

import (
	"github.com/b2b-platform/equipment-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(equipment *models.Equipment) error {
	return r.db.Create(equipment).Error
}

func (r *EquipmentRepository) GetByID(id uuid.UUID) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.Preload("BOMNodes").Preload("CompatibilityMappings").
		Where("id = ?", id).First(&equipment).Error
	return &equipment, err
}

func (r *EquipmentRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Equipment, error) {
	var equipment []models.Equipment
	query := r.db.Where("tenant_id = ?", tenantID)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&equipment).Error
	return equipment, err
}

func (r *EquipmentRepository) Update(equipment *models.Equipment) error {
	return r.db.Save(equipment).Error
}

func (r *EquipmentRepository) GetByEquipmentNumber(tenantID uuid.UUID, equipmentNumber string) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.Where("tenant_id = ? AND equipment_number = ?", tenantID, equipmentNumber).First(&equipment).Error
	return &equipment, err
}
