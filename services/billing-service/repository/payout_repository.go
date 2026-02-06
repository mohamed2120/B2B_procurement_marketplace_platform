package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository(db *gorm.DB) *PayoutRepository {
	return &PayoutRepository{db: db}
}

func (r *PayoutRepository) Create(account *models.PayoutAccount) error {
	return r.db.Create(account).Error
}

func (r *PayoutRepository) GetByID(id uuid.UUID) (*models.PayoutAccount, error) {
	var account models.PayoutAccount
	if err := r.db.First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *PayoutRepository) GetBySupplierID(supplierID uuid.UUID) ([]models.PayoutAccount, error) {
	var accounts []models.PayoutAccount
	err := r.db.Where("supplier_id = ?", supplierID).Find(&accounts).Error
	return accounts, err
}

func (r *PayoutRepository) GetDefaultBySupplierID(supplierID uuid.UUID) (*models.PayoutAccount, error) {
	var account models.PayoutAccount
	if err := r.db.First(&account, "supplier_id = ? AND is_default = ?", supplierID, true).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *PayoutRepository) Update(account *models.PayoutAccount) error {
	return r.db.Save(account).Error
}

func (r *PayoutRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PayoutAccount{}, "id = ?", id).Error
}

func (r *PayoutRepository) SetDefault(supplierID uuid.UUID, accountID uuid.UUID) error {
	// Unset all defaults for supplier
	if err := r.db.Model(&models.PayoutAccount{}).
		Where("supplier_id = ?", supplierID).
		Update("is_default", false).Error; err != nil {
		return err
	}

	// Set new default
	return r.db.Model(&models.PayoutAccount{}).
		Where("id = ?", accountID).
		Update("is_default", true).Error
}
