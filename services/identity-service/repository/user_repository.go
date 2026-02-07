package repository

import (
	"github.com/b2b-platform/identity-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserRoles.Role").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(tenantID uuid.UUID, email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserRoles.Role").Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) AssignRole(userID, roleID, tenantID uuid.UUID) error {
	userRole := models.UserRole{
		UserID:   userID,
		RoleID:   roleID,
		TenantID: tenantID,
	}
	return r.db.Create(&userRole).Error
}

func (r *UserRepository) GetUserRoles(userID, tenantID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Table("identity.roles").
		Joins("INNER JOIN identity.user_roles ON identity.roles.id = identity.user_roles.role_id").
		Where("identity.user_roles.user_id = ? AND identity.user_roles.tenant_id = ?", userID, tenantID).
		Find(&roles).Error
	return roles, err
}

// List all users (admin only - can filter by tenant)
func (r *UserRepository) List(tenantID *uuid.UUID) ([]models.User, error) {
	var users []models.User
	query := r.db.Preload("UserRoles.Role")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Find(&users).Error
	return users, err
}

// ToggleActive toggles user active status
func (r *UserRepository) ToggleActive(userID uuid.UUID, isActive bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive).Error
}
