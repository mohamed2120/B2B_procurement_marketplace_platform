package repository

import (
	"github.com/b2b-platform/notification-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Preload("Template").Where("id = ?", id).First(&notification).Error
	return &notification, err
}

func (r *NotificationRepository) GetUserNotifications(userID, tenantID uuid.UUID, limit, offset int, unreadOnly bool) ([]models.Notification, error) {
	var notifications []models.Notification
	query := r.db.Preload("Template").Where("user_id = ? AND tenant_id = ?", userID, tenantID)
	
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": true,
		}).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID, tenantID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND tenant_id = ? AND is_read = ?", userID, tenantID, false).
		Update("is_read", true).Error
}

func (r *NotificationRepository) GetPending() ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("status = ?", "pending").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) UpdateStatus(notificationID uuid.UUID, status string) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("status", status).Error
}
