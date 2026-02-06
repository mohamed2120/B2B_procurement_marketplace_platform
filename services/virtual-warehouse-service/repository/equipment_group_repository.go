package repository

import (
	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentGroupRepository struct {
	db *gorm.DB
}

func NewEquipmentGroupRepository(db *gorm.DB) *EquipmentGroupRepository {
	return &EquipmentGroupRepository{db: db}
}

func (r *EquipmentGroupRepository) Create(group *models.EquipmentGroup) error {
	return r.db.Create(group).Error
}

func (r *EquipmentGroupRepository) GetByID(id uuid.UUID) (*models.EquipmentGroup, error) {
	var group models.EquipmentGroup
	err := r.db.Preload("Members").Where("id = ?", id).First(&group).Error
	return &group, err
}

func (r *EquipmentGroupRepository) List(tenantID uuid.UUID) ([]models.EquipmentGroup, error) {
	var groups []models.EquipmentGroup
	err := r.db.Preload("Members").Where("tenant_id = ?", tenantID).Find(&groups).Error
	return groups, err
}

func (r *EquipmentGroupRepository) AddMember(member *models.EquipmentGroupMember) error {
	return r.db.Create(member).Error
}

func (r *EquipmentGroupRepository) RemoveMember(groupID, equipmentID uuid.UUID) error {
	return r.db.Where("group_id = ? AND equipment_id = ?", groupID, equipmentID).
		Delete(&models.EquipmentGroupMember{}).Error
}
