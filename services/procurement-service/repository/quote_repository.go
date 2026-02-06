package repository

import (
	"github.com/b2b-platform/procurement-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (r *QuoteRepository) Create(quote *models.Quote) error {
	return r.db.Create(quote).Error
}

func (r *QuoteRepository) GetByID(id uuid.UUID) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.Preload("Items").Preload("RFQ").Where("id = ?", id).First(&quote).Error
	return &quote, err
}

func (r *QuoteRepository) GetByRFQ(rfqID uuid.UUID) ([]models.Quote, error) {
	var quotes []models.Quote
	err := r.db.Preload("Items").Where("rfq_id = ?", rfqID).Find(&quotes).Error
	return quotes, err
}

func (r *QuoteRepository) Update(quote *models.Quote) error {
	return r.db.Save(quote).Error
}
