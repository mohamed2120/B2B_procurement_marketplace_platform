package service

import (
	"time"

	"github.com/b2b-platform/company-service/models"
	"github.com/b2b-platform/company-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type CompanyService struct {
	repo     *repository.CompanyRepository
	eventBus events.EventBus
}

func NewCompanyService(repo *repository.CompanyRepository, eventBus events.EventBus) *CompanyService {
	return &CompanyService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *CompanyService) Create(company *models.Company) error {
	return s.repo.Create(company)
}

func (s *CompanyService) GetByID(id uuid.UUID) (*models.Company, error) {
	return s.repo.GetByID(id)
}

func (s *CompanyService) List(limit, offset int) ([]models.Company, error) {
	return s.repo.List(limit, offset)
}

func (s *CompanyService) Update(company *models.Company) error {
	return s.repo.Update(company)
}

func (s *CompanyService) Approve(companyID, approvedBy uuid.UUID) error {
	company, err := s.repo.GetByID(companyID)
	if err != nil {
		return err
	}

	now := time.Now()
	company.Status = "approved"
	company.VerificationStatus = "verified"
	company.ApprovedAt = &now
	company.ApprovedBy = &approvedBy

	if err := s.repo.Update(company); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventCompanyApproved,
		"company-service",
		map[string]interface{}{
			"company_id": company.ID.String(),
			"name":       company.Name,
			"subdomain":  company.Subdomain,
		},
	).WithTenantID(companyID)

	return s.eventBus.Publish(nil, event)
}

func (s *CompanyService) RequestSubdomain(companyID uuid.UUID, subdomain string, requestedBy uuid.UUID) error {
	req := &models.SubdomainRequest{
		CompanyID:   companyID,
		Subdomain:   subdomain,
		Status:      "pending",
		RequestedBy: requestedBy,
	}
	return s.repo.CreateSubdomainRequest(req)
}
