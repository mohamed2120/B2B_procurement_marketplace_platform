package repository

import (
	"github.com/b2b-platform/company-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) GetByID(id uuid.UUID) (*models.Company, error) {
	var company models.Company
	err := r.db.Preload("Documents").Where("id = ?", id).First(&company).Error
	return &company, err
}

func (r *CompanyRepository) List(limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := r.db.Model(&models.Company{})
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) GetBySubdomain(subdomain string) (*models.Company, error) {
	var company models.Company
	err := r.db.Where("subdomain = ?", subdomain).First(&company).Error
	return &company, err
}

func (r *CompanyRepository) AddDocument(doc *models.CompanyDocument) error {
	return r.db.Create(doc).Error
}

func (r *CompanyRepository) CreateSubdomainRequest(req *models.SubdomainRequest) error {
	return r.db.Create(req).Error
}
