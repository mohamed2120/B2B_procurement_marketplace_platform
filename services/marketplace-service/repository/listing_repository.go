package repository

import (
	"github.com/b2b-platform/marketplace-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListingRepository struct {
	db *gorm.DB
}

func NewListingRepository(db *gorm.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) Create(listing *models.Listing) error {
	return r.db.Create(listing).Error
}

func (r *ListingRepository) GetByID(id uuid.UUID) (*models.Listing, error) {
	var listing models.Listing
	err := r.db.Preload("Store").Preload("Media").
		Where("id = ?", id).First(&listing).Error
	return &listing, err
}

func (r *ListingRepository) List(tenantID uuid.UUID, limit, offset int, listingType, status string) ([]models.Listing, error) {
	var listings []models.Listing
	query := r.db.Preload("Store").Preload("Media").Where("tenant_id = ?", tenantID)
	
	if listingType != "" {
		query = query.Where("listing_type = ?", listingType)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	err := query.Find(&listings).Error
	return listings, err
}

func (r *ListingRepository) Update(listing *models.Listing) error {
	return r.db.Save(listing).Error
}

func (r *ListingRepository) GetByStore(storeID uuid.UUID) ([]models.Listing, error) {
	var listings []models.Listing
	err := r.db.Preload("Media").Where("store_id = ?", storeID).Find(&listings).Error
	return listings, err
}

func (r *ListingRepository) UpdateStock(listingID uuid.UUID, quantity float64) error {
	return r.db.Model(&models.Listing{}).
		Where("id = ?", listingID).
		Update("stock_quantity", quantity).Error
}
