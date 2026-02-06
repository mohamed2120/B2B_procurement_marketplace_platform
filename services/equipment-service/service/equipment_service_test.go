package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/equipment-service/models"
	"github.com/google/uuid"
)

// MockEquipmentRepository for testing
type MockEquipmentRepository struct {
	equipment map[uuid.UUID]*models.Equipment
}

func (m *MockEquipmentRepository) Create(equipment *models.Equipment) error {
	if m.equipment == nil {
		m.equipment = make(map[uuid.UUID]*models.Equipment)
	}
	m.equipment[equipment.ID] = equipment
	return nil
}

func (m *MockEquipmentRepository) GetByID(id uuid.UUID) (*models.Equipment, error) {
	if eq, ok := m.equipment[id]; ok {
		return eq, nil
	}
	return nil, nil
}

func (m *MockEquipmentRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Equipment, error) {
	var result []models.Equipment
	for _, eq := range m.equipment {
		if eq.TenantID == tenantID {
			result = append(result, *eq)
		}
	}
	return result, nil
}

func (m *MockEquipmentRepository) Update(equipment *models.Equipment) error {
	m.equipment[equipment.ID] = equipment
	return nil
}

// MockBOMRepository for testing
type MockBOMRepository struct {
	nodes map[uuid.UUID]*models.BOMNode
}

func (m *MockBOMRepository) Create(node *models.BOMNode) error {
	if m.nodes == nil {
		m.nodes = make(map[uuid.UUID]*models.BOMNode)
	}
	m.nodes[node.ID] = node
	return nil
}

func (m *MockBOMRepository) GetByEquipment(equipmentID uuid.UUID) ([]models.BOMNode, error) {
	var result []models.BOMNode
	for _, node := range m.nodes {
		if node.EquipmentID == equipmentID {
			result = append(result, *node)
		}
	}
	return result, nil
}

// MockCompatibilityRepository for testing
type MockCompatibilityRepository struct {
	mappings map[uuid.UUID]*models.CompatibilityMapping
}

func (m *MockCompatibilityRepository) Create(mapping *models.CompatibilityMapping) error {
	if m.mappings == nil {
		m.mappings = make(map[uuid.UUID]*models.CompatibilityMapping)
	}
	m.mappings[mapping.ID] = mapping
	return nil
}

func (m *MockCompatibilityRepository) GetByID(id uuid.UUID) (*models.CompatibilityMapping, error) {
	if mapping, ok := m.mappings[id]; ok {
		return mapping, nil
	}
	return nil, nil
}

func (m *MockCompatibilityRepository) CheckCompatibility(equipmentID, partID uuid.UUID) (*models.CompatibilityMapping, error) {
	for _, mapping := range m.mappings {
		if mapping.EquipmentID == equipmentID && mapping.PartID == partID {
			return mapping, nil
		}
	}
	return nil, nil
}

func (m *MockCompatibilityRepository) GetByEquipment(equipmentID uuid.UUID) ([]models.CompatibilityMapping, error) {
	var result []models.CompatibilityMapping
	for _, mapping := range m.mappings {
		if mapping.EquipmentID == equipmentID {
			result = append(result, *mapping)
		}
	}
	return result, nil
}

func (m *MockCompatibilityRepository) VerifyCompatibility(mappingID, verifiedBy uuid.UUID) error {
	if mapping, ok := m.mappings[mappingID]; ok {
		now := time.Now()
		mapping.IsCompatible = true
		mapping.VerifiedBy = &verifiedBy
		mapping.VerifiedAt = &now
	}
	return nil
}

func TestEquipmentService_Create(t *testing.T) {
	mockEqRepo := &MockEquipmentRepository{}
	mockBOMRepo := &MockBOMRepository{}
	mockCompatRepo := &MockCompatibilityRepository{}

	service := NewEquipmentService(mockEqRepo, mockBOMRepo, mockCompatRepo)

	tenantID := uuid.New()
	equipment := &models.Equipment{
		ID:       uuid.New(),
		TenantID: tenantID,
		Name:     "Test Equipment",
		Model:    "MODEL-001",
	}

	err := service.Create(equipment)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify equipment was created
	created, _ := mockEqRepo.GetByID(equipment.ID)
	if created == nil {
		t.Errorf("expected equipment to be created")
	}
	if created.Name != "Test Equipment" {
		t.Errorf("expected name 'Test Equipment', got %s", created.Name)
	}
}

func TestEquipmentService_AddBOMNode(t *testing.T) {
	mockEqRepo := &MockEquipmentRepository{}
	mockBOMRepo := &MockBOMRepository{}
	mockCompatRepo := &MockCompatibilityRepository{}

	service := NewEquipmentService(mockEqRepo, mockBOMRepo, mockCompatRepo)

	equipmentID := uuid.New()
	node := &models.BOMNode{
		ID:          uuid.New(),
		EquipmentID: equipmentID,
		PartID:      uuid.New(),
		Quantity:    2.0,
	}

	err := service.AddBOMNode(node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify BOM node was created
	bom, _ := service.GetBOM(equipmentID)
	if len(bom) != 1 {
		t.Errorf("expected 1 BOM node, got %d", len(bom))
	}
	if bom[0].ID != node.ID {
		t.Errorf("expected BOM node ID %s, got %s", node.ID, bom[0].ID)
	}
}

func TestEquipmentService_VerifyCompatibility(t *testing.T) {
	mockEqRepo := &MockEquipmentRepository{}
	mockBOMRepo := &MockBOMRepository{}
	mockCompatRepo := &MockCompatibilityRepository{}

	service := NewEquipmentService(mockEqRepo, mockBOMRepo, mockCompatRepo)

	mappingID := uuid.New()
	verifiedBy := uuid.New()
	mapping := &models.CompatibilityMapping{
		ID:          mappingID,
		EquipmentID: uuid.New(),
		PartID:      uuid.New(),
		IsCompatible: false,
	}
	mockCompatRepo.Create(mapping)

	err := service.VerifyCompatibility(mappingID, verifiedBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify mapping was updated
	updated, _ := mockCompatRepo.GetByID(mappingID)
	if !updated.IsCompatible {
		t.Errorf("expected IsCompatible to be true")
	}
	if updated.VerifiedBy == nil || *updated.VerifiedBy != verifiedBy {
		t.Errorf("expected VerifiedBy to be set")
	}
	if updated.VerifiedAt == nil {
		t.Errorf("expected VerifiedAt to be set")
	}
}

func TestEquipmentService_CheckCompatibility(t *testing.T) {
	mockEqRepo := &MockEquipmentRepository{}
	mockBOMRepo := &MockBOMRepository{}
	mockCompatRepo := &MockCompatibilityRepository{}

	service := NewEquipmentService(mockEqRepo, mockBOMRepo, mockCompatRepo)

	equipmentID := uuid.New()
	partID := uuid.New()
	mapping := &models.CompatibilityMapping{
		ID:          uuid.New(),
		EquipmentID: equipmentID,
		PartID:      partID,
		IsCompatible: true,
	}
	mockCompatRepo.Create(mapping)

	result, err := service.CheckCompatibility(equipmentID, partID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected compatibility mapping to be found")
	}
	if !result.IsCompatible {
		t.Errorf("expected IsCompatible to be true")
	}
}

func TestEquipmentService_GetCompatibilityMappings(t *testing.T) {
	mockEqRepo := &MockEquipmentRepository{}
	mockBOMRepo := &MockBOMRepository{}
	mockCompatRepo := &MockCompatibilityRepository{}

	service := NewEquipmentService(mockEqRepo, mockBOMRepo, mockCompatRepo)

	equipmentID := uuid.New()
	mapping1 := &models.CompatibilityMapping{
		ID:          uuid.New(),
		EquipmentID: equipmentID,
		PartID:      uuid.New(),
		IsCompatible: true,
	}
	mapping2 := &models.CompatibilityMapping{
		ID:          uuid.New(),
		EquipmentID: equipmentID,
		PartID:      uuid.New(),
		IsCompatible: false,
	}
	mockCompatRepo.Create(mapping1)
	mockCompatRepo.Create(mapping2)

	mappings, err := service.GetCompatibilityMappings(equipmentID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mappings) != 2 {
		t.Errorf("expected 2 compatibility mappings, got %d", len(mappings))
	}
}
