package service

import (
	"time"

	"github.com/b2b-platform/equipment-service/models"
	"github.com/b2b-platform/equipment-service/repository"
	"github.com/google/uuid"
)

type EquipmentService struct {
	equipmentRepo      *repository.EquipmentRepository
	bomRepo            *repository.BOMRepository
	compatibilityRepo  *repository.CompatibilityRepository
}

func NewEquipmentService(
	equipmentRepo *repository.EquipmentRepository,
	bomRepo *repository.BOMRepository,
	compatibilityRepo *repository.CompatibilityRepository,
) *EquipmentService {
	return &EquipmentService{
		equipmentRepo:     equipmentRepo,
		bomRepo:           bomRepo,
		compatibilityRepo:  compatibilityRepo,
	}
}

func (s *EquipmentService) Create(equipment *models.Equipment) error {
	return s.equipmentRepo.Create(equipment)
}

func (s *EquipmentService) GetByID(id uuid.UUID) (*models.Equipment, error) {
	return s.equipmentRepo.GetByID(id)
}

func (s *EquipmentService) List(tenantID uuid.UUID, limit, offset int) ([]models.Equipment, error) {
	return s.equipmentRepo.List(tenantID, limit, offset)
}

func (s *EquipmentService) Update(equipment *models.Equipment) error {
	return s.equipmentRepo.Update(equipment)
}

func (s *EquipmentService) AddBOMNode(node *models.BOMNode) error {
	return s.bomRepo.Create(node)
}

func (s *EquipmentService) GetBOM(equipmentID uuid.UUID) ([]models.BOMNode, error) {
	return s.bomRepo.GetByEquipment(equipmentID)
}

func (s *EquipmentService) CreateCompatibilityMapping(mapping *models.CompatibilityMapping) error {
	return s.compatibilityRepo.Create(mapping)
}

func (s *EquipmentService) VerifyCompatibility(mappingID, verifiedBy uuid.UUID) error {
	now := time.Now()
	mapping, err := s.compatibilityRepo.GetByID(mappingID)
	if err != nil {
		return err
	}

	mapping.IsCompatible = true
	mapping.VerifiedBy = &verifiedBy
	mapping.VerifiedAt = &now

	return s.compatibilityRepo.VerifyCompatibility(mappingID, verifiedBy)
}

func (s *EquipmentService) CheckCompatibility(equipmentID, partID uuid.UUID) (*models.CompatibilityMapping, error) {
	return s.compatibilityRepo.CheckCompatibility(equipmentID, partID)
}

func (s *EquipmentService) GetCompatibilityMappings(equipmentID uuid.UUID) ([]models.CompatibilityMapping, error) {
	return s.compatibilityRepo.GetByEquipment(equipmentID)
}
