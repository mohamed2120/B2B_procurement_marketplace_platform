package service

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/repository"
	"github.com/google/uuid"
)

type PayoutService struct {
	payoutRepo *repository.PayoutRepository
}

func NewPayoutService(payoutRepo *repository.PayoutRepository) *PayoutService {
	return &PayoutService{
		payoutRepo: payoutRepo,
	}
}

func (s *PayoutService) CreatePayoutAccount(account *models.PayoutAccount) error {
	// If this is the first account for supplier, set as default
	existing, _ := s.payoutRepo.GetBySupplierID(account.SupplierID)
	if len(existing) == 0 {
		account.IsDefault = true
	}

	return s.payoutRepo.Create(account)
}

func (s *PayoutService) GetPayoutAccount(id uuid.UUID) (*models.PayoutAccount, error) {
	return s.payoutRepo.GetByID(id)
}

func (s *PayoutService) ListPayoutAccounts(supplierID uuid.UUID) ([]models.PayoutAccount, error) {
	return s.payoutRepo.GetBySupplierID(supplierID)
}

func (s *PayoutService) UpdatePayoutAccount(account *models.PayoutAccount) error {
	return s.payoutRepo.Update(account)
}

func (s *PayoutService) DeletePayoutAccount(id uuid.UUID) error {
	return s.payoutRepo.Delete(id)
}

func (s *PayoutService) SetDefaultPayoutAccount(supplierID uuid.UUID, accountID uuid.UUID) error {
	return s.payoutRepo.SetDefault(supplierID, accountID)
}
