package service

import (
	"time"

	"github.com/b2b-platform/catalog-service/models"
	"github.com/b2b-platform/catalog-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type CatalogService struct {
	manufacturerRepo *repository.ManufacturerRepository
	categoryRepo     *repository.CategoryRepository
	partRepo         *repository.PartRepository
	attributeRepo    *repository.AttributeRepository
	eventBus         events.EventBus
}

func NewCatalogService(
	manufacturerRepo *repository.ManufacturerRepository,
	categoryRepo *repository.CategoryRepository,
	partRepo *repository.PartRepository,
	attributeRepo *repository.AttributeRepository,
	eventBus events.EventBus,
) *CatalogService {
	return &CatalogService{
		manufacturerRepo: manufacturerRepo,
		categoryRepo:     categoryRepo,
		partRepo:         partRepo,
		attributeRepo:    attributeRepo,
		eventBus:         eventBus,
	}
}

func (s *CatalogService) CreateManufacturer(manufacturer *models.Manufacturer) error {
	return s.manufacturerRepo.Create(manufacturer)
}

func (s *CatalogService) GetManufacturer(id uuid.UUID) (*models.Manufacturer, error) {
	return s.manufacturerRepo.GetByID(id)
}

func (s *CatalogService) ListManufacturers() ([]models.Manufacturer, error) {
	return s.manufacturerRepo.List()
}

func (s *CatalogService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CatalogService) GetCategory(id uuid.UUID) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CatalogService) ListCategories() ([]models.Category, error) {
	return s.categoryRepo.List()
}

func (s *CatalogService) CreateAttribute(attribute *models.Attribute) error {
	return s.attributeRepo.Create(attribute)
}

func (s *CatalogService) ListAttributes() ([]models.Attribute, error) {
	return s.attributeRepo.List()
}

func (s *CatalogService) CreatePart(part *models.LibraryPart) error {
	// Check for duplicates
	existing, err := s.partRepo.FindDuplicate(part.PartNumber, part.ManufacturerID)
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if existing != nil {
		// Mark new part as duplicate
		if err := s.partRepo.Create(part); err != nil {
			return err
		}
		return s.partRepo.MarkAsDuplicate(part.ID, existing.ID)
	}

	return s.partRepo.Create(part)
}

func (s *CatalogService) GetPart(id uuid.UUID) (*models.LibraryPart, error) {
	return s.partRepo.GetByID(id)
}

func (s *CatalogService) ListParts(limit, offset int, status string) ([]models.LibraryPart, error) {
	return s.partRepo.List(limit, offset, status)
}

func (s *CatalogService) ApprovePart(partID, approvedBy uuid.UUID) error {
	part, err := s.partRepo.GetByID(partID)
	if err != nil {
		return err
	}

	now := time.Now()
	part.Status = "approved"
	part.ApprovedAt = &now
	part.ApprovedBy = &approvedBy

	if err := s.partRepo.Update(part); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventCatalogPartApproved,
		"catalog-service",
		map[string]interface{}{
			"part_id":     part.ID.String(),
			"part_number": part.PartNumber,
			"name":        part.Name,
			"manufacturer_id": part.ManufacturerID.String(),
		},
	)

	return s.eventBus.Publish(nil, event)
}

func (s *CatalogService) RejectPart(partID uuid.UUID, reason string) error {
	part, err := s.partRepo.GetByID(partID)
	if err != nil {
		return err
	}

	part.Status = "rejected"
	part.RejectedReason = reason

	return s.partRepo.Update(part)
}

func (s *CatalogService) GetPendingParts() ([]models.LibraryPart, error) {
	return s.partRepo.GetPendingApproval()
}

func (s *CatalogService) AddPartAttribute(partAttribute *models.PartAttribute) error {
	return s.attributeRepo.AddPartAttribute(partAttribute)
}
