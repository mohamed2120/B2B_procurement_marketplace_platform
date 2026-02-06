package repository

import (
	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *models.SharedInventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetByID(id uuid.UUID) (*models.SharedInventory, error) {
	var inventory models.SharedInventory
	err := r.db.Where("id = ?", id).First(&inventory).Error
	return &inventory, err
}

func (r *InventoryRepository) List(tenantID uuid.UUID) ([]models.SharedInventory, error) {
	var inventory []models.SharedInventory
	err := r.db.Where("tenant_id = ?", tenantID).Find(&inventory).Error
	return inventory, err
}

func (r *InventoryRepository) GetAvailable(partID uuid.UUID, quantity float64) ([]models.SharedInventory, error) {
	var inventory []models.SharedInventory
	err := r.db.Where("part_id = ? AND is_available = ? AND (quantity - reserved_qty) >= ?", 
		partID, true, quantity).Find(&inventory).Error
	return inventory, err
}

func (r *InventoryRepository) Reserve(inventoryID uuid.UUID, quantity float64) error {
	return r.db.Model(&models.SharedInventory{}).
		Where("id = ?", inventoryID).
		UpdateColumn("reserved_qty", gorm.Expr("reserved_qty + ?", quantity)).Error
}

func (r *InventoryRepository) Release(inventoryID uuid.UUID, quantity float64) error {
	return r.db.Model(&models.SharedInventory{}).
		Where("id = ?", inventoryID).
		UpdateColumn("reserved_qty", gorm.Expr("GREATEST(0, reserved_qty - ?)", quantity)).Error
}
