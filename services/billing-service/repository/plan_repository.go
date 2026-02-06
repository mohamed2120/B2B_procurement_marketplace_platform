package repository

import (
	"github.com/b2b-platform/billing-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(plan *models.Plan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepository) GetByID(id uuid.UUID) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Preload("Entitlements").Where("id = ?", id).First(&plan).Error
	return &plan, err
}

func (r *PlanRepository) GetByCode(code string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Preload("Entitlements").Where("code = ?", code).First(&plan).Error
	return &plan, err
}

func (r *PlanRepository) List() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Preload("Entitlements").Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}
