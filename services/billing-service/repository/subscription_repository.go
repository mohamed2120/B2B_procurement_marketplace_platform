package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("Plan").Where("id = ?", id).First(&subscription).Error
	return &subscription, err
}

func (r *SubscriptionRepository) GetByTenant(tenantID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("Plan").Preload("Plan.Entitlements").
		Where("tenant_id = ? AND status = ?", tenantID, "active").First(&subscription).Error
	return &subscription, err
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *SubscriptionRepository) Cancel(subscriptionID uuid.UUID) error {
	return r.db.Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"cancelled_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}
