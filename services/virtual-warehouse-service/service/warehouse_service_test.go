package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/google/uuid"
)

// MockInventoryRepository for testing
type MockInventoryRepository struct {
	inventory map[uuid.UUID]*models.SharedInventory
}

func (m *MockInventoryRepository) Create(inventory *models.SharedInventory) error {
	if m.inventory == nil {
		m.inventory = make(map[uuid.UUID]*models.SharedInventory)
	}
	m.inventory[inventory.ID] = inventory
	return nil
}

func (m *MockInventoryRepository) List(tenantID uuid.UUID) ([]models.SharedInventory, error) {
	var result []models.SharedInventory
	for _, inv := range m.inventory {
		if inv.TenantID == tenantID {
			result = append(result, *inv)
		}
	}
	return result, nil
}

func (m *MockInventoryRepository) GetAvailable(partID uuid.UUID, quantity float64) ([]models.SharedInventory, error) {
	var result []models.SharedInventory
	for _, inv := range m.inventory {
		if inv.PartID == partID && inv.AvailableQuantity >= quantity {
			result = append(result, *inv)
		}
	}
	return result, nil
}

func (m *MockInventoryRepository) Reserve(inventoryID uuid.UUID, quantity float64) error {
	if inv, ok := m.inventory[inventoryID]; ok {
		if inv.AvailableQuantity < quantity {
			return nil // Would return error in real implementation
		}
		inv.AvailableQuantity -= quantity
		inv.ReservedQuantity += quantity
	}
	return nil
}

// MockEquipmentGroupRepository for testing
type MockEquipmentGroupRepository struct {
	groups  map[uuid.UUID]*models.EquipmentGroup
	members []models.EquipmentGroupMember
}

func (m *MockEquipmentGroupRepository) Create(group *models.EquipmentGroup) error {
	if m.groups == nil {
		m.groups = make(map[uuid.UUID]*models.EquipmentGroup)
	}
	m.groups[group.ID] = group
	return nil
}

func (m *MockEquipmentGroupRepository) GetByID(id uuid.UUID) (*models.EquipmentGroup, error) {
	if group, ok := m.groups[id]; ok {
		return group, nil
	}
	return nil, nil
}

func (m *MockEquipmentGroupRepository) List(tenantID uuid.UUID) ([]models.EquipmentGroup, error) {
	var result []models.EquipmentGroup
	for _, group := range m.groups {
		if group.TenantID == tenantID {
			result = append(result, *group)
		}
	}
	return result, nil
}

func (m *MockEquipmentGroupRepository) AddMember(member *models.EquipmentGroupMember) error {
	if m.members == nil {
		m.members = make([]models.EquipmentGroupMember, 0)
	}
	m.members = append(m.members, *member)
	return nil
}

// MockTransferRepository for testing
type MockTransferRepository struct {
	transfers map[uuid.UUID]*models.InterCompanyTransfer
}

func (m *MockTransferRepository) Create(transfer *models.InterCompanyTransfer) error {
	if m.transfers == nil {
		m.transfers = make(map[uuid.UUID]*models.InterCompanyTransfer)
	}
	m.transfers[transfer.ID] = transfer
	return nil
}

func (m *MockTransferRepository) GetByID(id uuid.UUID) (*models.InterCompanyTransfer, error) {
	if transfer, ok := m.transfers[id]; ok {
		return transfer, nil
	}
	return nil, nil
}

func (m *MockTransferRepository) List(tenantID uuid.UUID) ([]models.InterCompanyTransfer, error) {
	var result []models.InterCompanyTransfer
	for _, transfer := range m.transfers {
		if transfer.FromTenantID == tenantID || transfer.ToTenantID == tenantID {
			result = append(result, *transfer)
		}
	}
	return result, nil
}

func (m *MockTransferRepository) Approve(transferID, approvedBy uuid.UUID) error {
	if transfer, ok := m.transfers[transferID]; ok {
		now := time.Now()
		transfer.Status = "approved"
		transfer.ApprovedBy = &approvedBy
		transfer.ApprovedAt = &now
	}
	return nil
}

func (m *MockTransferRepository) Reject(transferID uuid.UUID, reason string) error {
	if transfer, ok := m.transfers[transferID]; ok {
		transfer.Status = "rejected"
		transfer.RejectionReason = reason
	}
	return nil
}

// MockEmergencyRepository for testing
type MockEmergencyRepository struct {
	sourcing map[uuid.UUID]*models.EmergencySourcing
}

func (m *MockEmergencyRepository) Create(sourcing *models.EmergencySourcing) error {
	if m.sourcing == nil {
		m.sourcing = make(map[uuid.UUID]*models.EmergencySourcing)
	}
	m.sourcing[sourcing.ID] = sourcing
	return nil
}

func (m *MockEmergencyRepository) GetByID(id uuid.UUID) (*models.EmergencySourcing, error) {
	if s, ok := m.sourcing[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *MockEmergencyRepository) List(tenantID uuid.UUID, status string) ([]models.EmergencySourcing, error) {
	var result []models.EmergencySourcing
	for _, s := range m.sourcing {
		if s.TenantID == tenantID {
			if status == "" || s.Status == status {
				result = append(result, *s)
			}
		}
	}
	return result, nil
}

func (m *MockEmergencyRepository) Update(sourcing *models.EmergencySourcing) error {
	m.sourcing[sourcing.ID] = sourcing
	return nil
}

func TestWarehouseService_CreateInventory(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	tenantID := uuid.New()
	inventory := &models.SharedInventory{
		ID:               uuid.New(),
		TenantID:        tenantID,
		PartID:          uuid.New(),
		TotalQuantity:   100.0,
		AvailableQuantity: 100.0,
	}

	err := service.CreateInventory(inventory)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify inventory was created
	list, _ := service.ListInventory(tenantID)
	if len(list) != 1 {
		t.Errorf("expected 1 inventory item, got %d", len(list))
	}
	if list[0].TotalQuantity != 100.0 {
		t.Errorf("expected total quantity 100.0, got %.2f", list[0].TotalQuantity)
	}
}

func TestWarehouseService_Reserve(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	inventoryID := uuid.New()
	inventory := &models.SharedInventory{
		ID:               inventoryID,
		AvailableQuantity: 100.0,
		ReservedQuantity: 0.0,
	}
	mockInvRepo.Create(inventory)

	err := service.Reserve(inventoryID, 30.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify inventory was reserved
	updated, _ := mockInvRepo.GetByID(inventoryID)
	if updated.AvailableQuantity != 70.0 {
		t.Errorf("expected available quantity 70.0, got %.2f", updated.AvailableQuantity)
	}
	if updated.ReservedQuantity != 30.0 {
		t.Errorf("expected reserved quantity 30.0, got %.2f", updated.ReservedQuantity)
	}
}

func TestWarehouseService_CreateTransfer(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	fromTenantID := uuid.New()
	toTenantID := uuid.New()
	transfer := &models.InterCompanyTransfer{
		ID:           uuid.New(),
		FromTenantID: fromTenantID,
		ToTenantID:   toTenantID,
		PartID:       uuid.New(),
		Quantity:     50.0,
		Status:       "pending",
	}

	err := service.CreateTransfer(transfer)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify transfer was created
	created, _ := mockTransferRepo.GetByID(transfer.ID)
	if created == nil {
		t.Errorf("expected transfer to be created")
	}
	if created.Quantity != 50.0 {
		t.Errorf("expected quantity 50.0, got %.2f", created.Quantity)
	}
}

func TestWarehouseService_ApproveTransfer(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	transferID := uuid.New()
	approvedBy := uuid.New()
	transfer := &models.InterCompanyTransfer{
		ID:     transferID,
		Status: "pending",
	}
	mockTransferRepo.Create(transfer)

	err := service.ApproveTransfer(transferID, approvedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify transfer was approved
	updated, _ := mockTransferRepo.GetByID(transferID)
	if updated.Status != "approved" {
		t.Errorf("expected status 'approved', got %s", updated.Status)
	}
	if updated.ApprovedBy == nil || *updated.ApprovedBy != approvedBy {
		t.Errorf("expected ApprovedBy to be set")
	}
	if updated.ApprovedAt == nil {
		t.Errorf("expected ApprovedAt to be set")
	}
}

func TestWarehouseService_RejectTransfer(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	transferID := uuid.New()
	transfer := &models.InterCompanyTransfer{
		ID:     transferID,
		Status: "pending",
	}
	mockTransferRepo.Create(transfer)

	reason := "Insufficient inventory"
	err := service.RejectTransfer(transferID, reason)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify transfer was rejected
	updated, _ := mockTransferRepo.GetByID(transferID)
	if updated.Status != "rejected" {
		t.Errorf("expected status 'rejected', got %s", updated.Status)
	}
	if updated.RejectionReason != reason {
		t.Errorf("expected rejection reason '%s', got %s", reason, updated.RejectionReason)
	}
}

func TestWarehouseService_FulfillEmergencySourcing(t *testing.T) {
	mockInvRepo := &MockInventoryRepository{}
	mockGroupRepo := &MockEquipmentGroupRepository{}
	mockTransferRepo := &MockTransferRepository{}
	mockEmergencyRepo := &MockEmergencyRepository{}

	service := NewWarehouseService(mockInvRepo, mockGroupRepo, mockTransferRepo, mockEmergencyRepo)

	sourcingID := uuid.New()
	fulfilledBy := uuid.New()
	sourcing := &models.EmergencySourcing{
		ID:     sourcingID,
		Status: "pending",
	}
	mockEmergencyRepo.Create(sourcing)

	err := service.FulfillEmergencySourcing(sourcingID, fulfilledBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify sourcing was fulfilled
	updated, _ := mockEmergencyRepo.GetByID(sourcingID)
	if updated.Status != "fulfilled" {
		t.Errorf("expected status 'fulfilled', got %s", updated.Status)
	}
	if updated.FulfilledBy == nil || *updated.FulfilledBy != fulfilledBy {
		t.Errorf("expected FulfilledBy to be set")
	}
	if updated.FulfilledAt == nil {
		t.Errorf("expected FulfilledAt to be set")
	}
}
