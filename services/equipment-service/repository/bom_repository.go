package repository

import (
	"github.com/b2b-platform/equipment-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BOMRepository struct {
	db *gorm.DB
}

func NewBOMRepository(db *gorm.DB) *BOMRepository {
	return &BOMRepository{db: db}
}

func (r *BOMRepository) Create(node *models.BOMNode) error {
	return r.db.Create(node).Error
}

func (r *BOMRepository) GetByID(id uuid.UUID) (*models.BOMNode, error) {
	var node models.BOMNode
	err := r.db.Preload("ParentNode").Preload("ChildNodes").Where("id = ?", id).First(&node).Error
	return &node, err
}

func (r *BOMRepository) GetByEquipment(equipmentID uuid.UUID) ([]models.BOMNode, error) {
	var nodes []models.BOMNode
	err := r.db.Preload("ParentNode").Preload("ChildNodes").
		Where("equipment_id = ? AND parent_node_id IS NULL", equipmentID).
		Find(&nodes).Error
	return nodes, err
}

func (r *BOMRepository) Update(node *models.BOMNode) error {
	return r.db.Save(node).Error
}

func (r *BOMRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BOMNode{}, id).Error
}
