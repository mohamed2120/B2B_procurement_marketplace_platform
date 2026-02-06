package service

import (
	"testing"

	"github.com/b2b-platform/identity-service/models"
	"github.com/google/uuid"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[uuid.UUID]*models.User
	roles map[uuid.UUID][]models.Role
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.users == nil {
		m.users = make(map[uuid.UUID]*models.User)
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(tenantID uuid.UUID, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email && user.TenantID == tenantID {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) AssignRole(userID, roleID, tenantID uuid.UUID) error {
	if m.roles == nil {
		m.roles = make(map[uuid.UUID][]models.Role)
	}
	role := models.Role{ID: roleID, Name: "test-role"}
	m.roles[userID] = append(m.roles[userID], role)
	return nil
}

func (m *MockUserRepository) GetUserRoles(userID, tenantID uuid.UUID) ([]models.Role, error) {
	if roles, ok := m.roles[userID]; ok {
		return roles, nil
	}
	return []models.Role{}, nil
}

// MockRoleRepository for testing
type MockRoleRepository struct {
	roles map[uuid.UUID]*models.Role
}

func (m *MockRoleRepository) Create(role *models.Role) error {
	if m.roles == nil {
		m.roles = make(map[uuid.UUID]*models.Role)
	}
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	if role, ok := m.roles[id]; ok {
		return role, nil
	}
	return nil, nil
}

func (m *MockRoleRepository) GetByName(name string) (*models.Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}

func TestUserService_Create(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	service := NewUserService(mockUserRepo, mockRoleRepo)

	tenantID := uuid.New()
	user := &models.User{
		ID:       uuid.New(),
		TenantID: tenantID,
		Email:    "test@example.com",
		Name:     "Test User",
	}

	err := service.Create(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify user was created
	created, _ := mockUserRepo.GetByID(user.ID)
	if created == nil {
		t.Errorf("expected user to be created")
	}
	if created.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", created.Email)
	}
}

func TestUserService_GetByID(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	service := NewUserService(mockUserRepo, mockRoleRepo)

	userID := uuid.New()
	tenantID := uuid.New()
	user := &models.User{
		ID:       userID,
		TenantID: tenantID,
		Email:    "test@example.com",
	}
	mockUserRepo.Create(user)

	result, err := service.GetByID(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected user to be found")
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", result.Email)
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	service := NewUserService(mockUserRepo, mockRoleRepo)

	tenantID := uuid.New()
	user := &models.User{
		ID:       uuid.New(),
		TenantID: tenantID,
		Email:    "test@example.com",
	}
	mockUserRepo.Create(user)

	result, err := service.GetByEmail(tenantID, "test@example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected user to be found")
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", result.Email)
	}
}

func TestUserService_AssignRole(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	service := NewUserService(mockUserRepo, mockRoleRepo)

	userID := uuid.New()
	roleID := uuid.New()
	tenantID := uuid.New()

	err := service.AssignRole(userID, roleID, tenantID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify role was assigned
	roles, _ := service.GetUserRoles(userID, tenantID)
	if len(roles) == 0 {
		t.Errorf("expected role to be assigned")
	}
	if roles[0].ID != roleID {
		t.Errorf("expected role ID %s, got %s", roleID, roles[0].ID)
	}
}

func TestUserService_GetUserRoles(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	service := NewUserService(mockUserRepo, mockRoleRepo)

	userID := uuid.New()
	roleID := uuid.New()
	tenantID := uuid.New()

	// Assign role first
	service.AssignRole(userID, roleID, tenantID)

	roles, err := service.GetUserRoles(userID, tenantID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}
	if roles[0].ID != roleID {
		t.Errorf("expected role ID %s, got %s", roleID, roles[0].ID)
	}
}
