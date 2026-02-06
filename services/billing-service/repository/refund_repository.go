package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundRepository struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) *RefundRepository {
	return &RefundRepository{db: db}
}

func (r *RefundRepository) Create(refund *models.Refund) error {
	return r.db.Create(refund).Error
}

func (r *RefundRepository) GetByID(id uuid.UUID) (*models.Refund, error) {
	var refund models.Refund
	if err := r.db.Preload("Payment").First(&refund, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepository) GetByRefundNumber(tenantID uuid.UUID, refundNumber string) (*models.Refund, error) {
	var refund models.Refund
	if err := r.db.Preload("Payment").First(&refund, "tenant_id = ? AND refund_number = ?", tenantID, refundNumber).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepository) Update(refund *models.Refund) error {
	return r.db.Save(refund).Error
}

func (r *RefundRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Refund, error) {
	var refunds []models.Refund
	err := r.db.Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&refunds).Error
	return refunds, err
}
