package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettlementRepository struct {
	db *gorm.DB
}

func NewSettlementRepository(db *gorm.DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

func (r *SettlementRepository) Create(settlement *models.Settlement) error {
	return r.db.Create(settlement).Error
}

func (r *SettlementRepository) GetByID(id uuid.UUID) (*models.Settlement, error) {
	var settlement models.Settlement
	if err := r.db.Preload("EscrowHold").Preload("PayoutAccount").First(&settlement, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &settlement, nil
}

func (r *SettlementRepository) GetByEscrowHoldID(escrowHoldID uuid.UUID) (*models.Settlement, error) {
	var settlement models.Settlement
	if err := r.db.Preload("EscrowHold").Preload("PayoutAccount").First(&settlement, "escrow_hold_id = ?", escrowHoldID).Error; err != nil {
		return nil, err
	}
	return &settlement, nil
}

func (r *SettlementRepository) Update(settlement *models.Settlement) error {
	return r.db.Save(settlement).Error
}

func (r *SettlementRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	var settlements []models.Settlement
	err := r.db.Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&settlements).Error
	return settlements, err
}
