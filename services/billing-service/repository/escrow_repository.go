package repository

import (
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EscrowRepository struct {
	db *gorm.DB
}

func NewEscrowRepository(db *gorm.DB) *EscrowRepository {
	return &EscrowRepository{db: db}
}

func (r *EscrowRepository) Create(hold *models.EscrowHold) error {
	return r.db.Create(hold).Error
}

func (r *EscrowRepository) GetByID(id uuid.UUID) (*models.EscrowHold, error) {
	var hold models.EscrowHold
	if err := r.db.Preload("Payment").First(&hold, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &hold, nil
}

func (r *EscrowRepository) GetByPaymentID(paymentID uuid.UUID) (*models.EscrowHold, error) {
	var hold models.EscrowHold
	if err := r.db.Preload("Payment").First(&hold, "payment_id = ?", paymentID).Error; err != nil {
		return nil, err
	}
	return &hold, nil
}

func (r *EscrowRepository) GetByOrderID(orderID uuid.UUID) (*models.EscrowHold, error) {
	var hold models.EscrowHold
	if err := r.db.Preload("Payment").First(&hold, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &hold, nil
}

func (r *EscrowRepository) Update(hold *models.EscrowHold) error {
	return r.db.Save(hold).Error
}

func (r *EscrowRepository) ListPendingRelease(tenantID uuid.UUID) ([]models.EscrowHold, error) {
	var holds []models.EscrowHold
	now := time.Now()
	err := r.db.Where("tenant_id = ? AND status = ? AND blocked_by_dispute = ? AND (auto_release_date IS NULL OR auto_release_date <= ?)",
		tenantID, "held", false, now).
		Find(&holds).Error
	return holds, err
}

func (r *EscrowRepository) ListBySupplier(supplierID uuid.UUID, limit, offset int) ([]models.EscrowHold, error) {
	var holds []models.EscrowHold
	err := r.db.Where("supplier_id = ?", supplierID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&holds).Error
	return holds, err
}
