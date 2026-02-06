package service

import (
	"github.com/b2b-platform/identity-service/models"
	"github.com/b2b-platform/identity-service/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetByEmail(tenantID uuid.UUID, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(tenantID, email)
}

func (s *UserService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) AssignRole(userID, roleID, tenantID uuid.UUID) error {
	return s.userRepo.AssignRole(userID, roleID, tenantID)
}

func (s *UserService) GetUserRoles(userID, tenantID uuid.UUID) ([]models.Role, error) {
	return s.userRepo.GetUserRoles(userID, tenantID)
}
