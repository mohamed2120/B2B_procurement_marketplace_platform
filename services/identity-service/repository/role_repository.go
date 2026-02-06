package repository

import (
	"github.com/b2b-platform/identity-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("RolePermissions.Permission").Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) AssignPermission(roleID, permissionID uuid.UUID) error {
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(&rolePermission).Error
}
