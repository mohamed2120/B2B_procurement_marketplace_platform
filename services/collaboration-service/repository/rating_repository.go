package repository

import (
	"github.com/b2b-platform/collaboration-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Create(rating *models.Rating) error {
	return r.db.Create(rating).Error
}

func (r *RatingRepository) GetByID(id uuid.UUID) (*models.Rating, error) {
	var rating models.Rating
	err := r.db.Where("id = ?", id).First(&rating).Error
	return &rating, err
}

func (r *RatingRepository) GetByEntity(entityType string, entityID uuid.UUID) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Where("rated_entity_type = ? AND rated_entity_id = ?", entityType, entityID).
		Order("created_at DESC").Find(&ratings).Error
	return ratings, err
}

func (r *RatingRepository) GetAverageRating(entityType string, entityID uuid.UUID) (float64, error) {
	var result struct {
		Average float64
		Count   int64
	}
	err := r.db.Model(&models.Rating{}).
		Select("AVG(rating) as average, COUNT(*) as count").
		Where("rated_entity_type = ? AND rated_entity_id = ?", entityType, entityID).
		Scan(&result).Error
	return result.Average, err
}

func (r *RatingRepository) Moderate(ratingID, moderatedBy uuid.UUID) error {
	return r.db.Model(&models.Rating{}).
		Where("id = ?", ratingID).
		Updates(map[string]interface{}{
			"is_moderated": true,
			"moderated_by": moderatedBy,
		}).Error
}
