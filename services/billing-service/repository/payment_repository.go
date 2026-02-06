package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByPaymentIntentID(paymentIntentID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "payment_intent_id = ?", paymentIntentID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}
