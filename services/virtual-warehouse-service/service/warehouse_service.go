package service

import (
	"time"

	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/b2b-platform/virtual-warehouse-service/repository"
	"github.com/google/uuid"
)

type WarehouseService struct {
	inventoryRepo  *repository.InventoryRepository
	groupRepo      *repository.EquipmentGroupRepository
	transferRepo   *repository.TransferRepository
	emergencyRepo  *repository.EmergencyRepository
}

func NewWarehouseService(
	inventoryRepo *repository.InventoryRepository,
	groupRepo *repository.EquipmentGroupRepository,
	transferRepo *repository.TransferRepository,
	emergencyRepo *repository.EmergencyRepository,
) *WarehouseService {
	return &WarehouseService{
		inventoryRepo: inventoryRepo,
		groupRepo:     groupRepo,
		transferRepo:  transferRepo,
		emergencyRepo: emergencyRepo,
	}
}

func (s *WarehouseService) CreateInventory(inventory *models.SharedInventory) error {
	return s.inventoryRepo.Create(inventory)
}

func (s *WarehouseService) ListInventory(tenantID uuid.UUID) ([]models.SharedInventory, error) {
	return s.inventoryRepo.List(tenantID)
}

func (s *WarehouseService) GetAvailable(partID uuid.UUID, quantity float64) ([]models.SharedInventory, error) {
	return s.inventoryRepo.GetAvailable(partID, quantity)
}

func (s *WarehouseService) Reserve(inventoryID uuid.UUID, quantity float64) error {
	return s.inventoryRepo.Reserve(inventoryID, quantity)
}

func (s *WarehouseService) CreateGroup(group *models.EquipmentGroup) error {
	return s.groupRepo.Create(group)
}

func (s *WarehouseService) GetGroup(id uuid.UUID) (*models.EquipmentGroup, error) {
	return s.groupRepo.GetByID(id)
}

func (s *WarehouseService) ListGroups(tenantID uuid.UUID) ([]models.EquipmentGroup, error) {
	return s.groupRepo.List(tenantID)
}

func (s *WarehouseService) AddGroupMember(member *models.EquipmentGroupMember) error {
	return s.groupRepo.AddMember(member)
}

func (s *WarehouseService) CreateTransfer(transfer *models.InterCompanyTransfer) error {
	return s.transferRepo.Create(transfer)
}

func (s *WarehouseService) GetTransfer(id uuid.UUID) (*models.InterCompanyTransfer, error) {
	return s.transferRepo.GetByID(id)
}

func (s *WarehouseService) ListTransfers(tenantID uuid.UUID) ([]models.InterCompanyTransfer, error) {
	return s.transferRepo.List(tenantID)
}

func (s *WarehouseService) ApproveTransfer(transferID, approvedBy uuid.UUID) error {
	return s.transferRepo.Approve(transferID, approvedBy)
}

func (s *WarehouseService) RejectTransfer(transferID uuid.UUID, reason string) error {
	return s.transferRepo.Reject(transferID, reason)
}

func (s *WarehouseService) CreateEmergencySourcing(sourcing *models.EmergencySourcing) error {
	return s.emergencyRepo.Create(sourcing)
}

func (s *WarehouseService) ListEmergencySourcing(tenantID uuid.UUID, status string) ([]models.EmergencySourcing, error) {
	return s.emergencyRepo.List(tenantID, status)
}

func (s *WarehouseService) FulfillEmergencySourcing(sourcingID, fulfilledBy uuid.UUID) error {
	now := time.Now()
	sourcing, err := s.emergencyRepo.GetByID(sourcingID)
	if err != nil {
		return err
	}

	sourcing.Status = "fulfilled"
	sourcing.FulfilledBy = &fulfilledBy
	sourcing.FulfilledAt = &now

	return s.emergencyRepo.Update(sourcing)
}
