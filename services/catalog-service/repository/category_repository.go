package repository

import (
	"github.com/b2b-platform/catalog-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *CategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("is_active = ? AND parent_id IS NULL", true).Preload("Children").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByCode(code string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("code = ?", code).First(&category).Error
	return &category, err
}
