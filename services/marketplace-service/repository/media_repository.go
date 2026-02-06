package repository

import (
	"github.com/b2b-platform/marketplace-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(media *models.ListingMedia) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) GetByID(id uuid.UUID) (*models.ListingMedia, error) {
	var media models.ListingMedia
	err := r.db.Preload("Listing").Where("id = ?", id).First(&media).Error
	return &media, err
}

func (r *MediaRepository) GetByListing(listingID uuid.UUID) ([]models.ListingMedia, error) {
	var media []models.ListingMedia
	err := r.db.Where("listing_id = ?", listingID).Order("sort_order ASC").Find(&media).Error
	return media, err
}

func (r *MediaRepository) Update(media *models.ListingMedia) error {
	return r.db.Save(media).Error
}

func (r *MediaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ListingMedia{}, id).Error
}

func (r *MediaRepository) SetPrimary(listingID, mediaID uuid.UUID) error {
	// Unset all primary flags for this listing
	if err := r.db.Model(&models.ListingMedia{}).
		Where("listing_id = ?", listingID).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	// Set the specified media as primary
	return r.db.Model(&models.ListingMedia{}).
		Where("id = ?", mediaID).
		Update("is_primary", true).Error
}
