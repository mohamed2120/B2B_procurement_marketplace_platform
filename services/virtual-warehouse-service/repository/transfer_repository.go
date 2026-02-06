package repository

import (
	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Create(transfer *models.InterCompanyTransfer) error {
	return r.db.Create(transfer).Error
}

func (r *TransferRepository) GetByID(id uuid.UUID) (*models.InterCompanyTransfer, error) {
	var transfer models.InterCompanyTransfer
	err := r.db.Where("id = ?", id).First(&transfer).Error
	return &transfer, err
}

func (r *TransferRepository) List(tenantID uuid.UUID) ([]models.InterCompanyTransfer, error) {
	var transfers []models.InterCompanyTransfer
	err := r.db.Where("from_tenant_id = ? OR to_tenant_id = ?", tenantID, tenantID).
		Order("created_at DESC").Find(&transfers).Error
	return transfers, err
}

func (r *TransferRepository) Update(transfer *models.InterCompanyTransfer) error {
	return r.db.Save(transfer).Error
}

func (r *TransferRepository) Approve(transferID, approvedBy uuid.UUID) error {
	return r.db.Model(&models.InterCompanyTransfer{}).
		Where("id = ?", transferID).
		Updates(map[string]interface{}{
			"status":     "approved",
			"approved_by": approvedBy,
		}).Error
}

func (r *TransferRepository) Reject(transferID uuid.UUID, reason string) error {
	return r.db.Model(&models.InterCompanyTransfer{}).
		Where("id = ?", transferID).
		Updates(map[string]interface{}{
			"status":          "rejected",
			"rejection_reason": reason,
		}).Error
}
