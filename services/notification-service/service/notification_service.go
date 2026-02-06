package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/b2b-platform/notification-service/models"
	"github.com/b2b-platform/notification-service/repository"
	"github.com/google/uuid"
)

type NotificationService struct {
	templateRepo    *repository.TemplateRepository
	preferenceRepo *repository.PreferenceRepository
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(
	templateRepo *repository.TemplateRepository,
	preferenceRepo *repository.PreferenceRepository,
	notificationRepo *repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		templateRepo:     templateRepo,
		preferenceRepo:   preferenceRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) CreateTemplate(template *models.NotificationTemplate) error {
	return s.templateRepo.Create(template)
}

func (s *NotificationService) GetTemplate(id uuid.UUID) (*models.NotificationTemplate, error) {
	return s.templateRepo.GetByID(id)
}

func (s *NotificationService) GetTemplateByCode(code string) (*models.NotificationTemplate, error) {
	return s.templateRepo.GetByCode(code)
}

func (s *NotificationService) ListTemplates() ([]models.NotificationTemplate, error) {
	return s.templateRepo.List()
}

func (s *NotificationService) CreatePreference(preference *models.NotificationPreference) error {
	return s.preferenceRepo.Create(preference)
}

func (s *NotificationService) GetUserPreferences(userID, tenantID uuid.UUID) ([]models.NotificationPreference, error) {
	return s.preferenceRepo.GetUserPreferences(userID, tenantID)
}

func (s *NotificationService) UpdatePreference(preference *models.NotificationPreference) error {
	return s.preferenceRepo.Update(preference)
}

func (s *NotificationService) SendNotification(notification *models.Notification) error {
	// Check user preference
	enabled, err := s.preferenceRepo.IsEnabled(notification.UserID, notification.TenantID, notification.Channel, notification.Type)
	if err != nil {
		return err
	}

	if !enabled {
		notification.Status = "skipped"
		return s.notificationRepo.Create(notification)
	}

	// Create notification
	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send notification (mock for local development)
	if err := s.send(notification); err != nil {
		now := time.Now()
		notification.Status = "failed"
		notification.FailedAt = &now
		notification.Error = err.Error()
		s.notificationRepo.UpdateStatus(notification.ID, "failed")
		return err
	}

	now := time.Now()
	notification.Status = "sent"
	notification.SentAt = &now
	return s.notificationRepo.UpdateStatus(notification.ID, "sent")
}

func (s *NotificationService) send(notification *models.Notification) error {
	// Mock implementation for local development
	// In production, this would send via email service, push notification service, etc.
	fmt.Printf("[MOCK] Sending %s notification to user %s: %s\n", 
		notification.Channel, notification.UserID, notification.Title)
	return nil
}

func (s *NotificationService) GetUserNotifications(userID, tenantID uuid.UUID, limit, offset int, unreadOnly bool) ([]models.Notification, error) {
	return s.notificationRepo.GetUserNotifications(userID, tenantID, limit, offset, unreadOnly)
}

func (s *NotificationService) MarkAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID, tenantID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID, tenantID)
}

func (s *NotificationService) CreateFromTemplate(templateCode string, userID, tenantID uuid.UUID, data map[string]interface{}) error {
	template, err := s.templateRepo.GetByCode(templateCode)
	if err != nil {
		return err
	}

	// Simple template variable replacement (in production, use a proper templating engine)
	message := template.Body
	title := template.Subject

	dataJSON, _ := json.Marshal(data)

	notification := &models.Notification{
		TenantID:   tenantID,
		UserID:     userID,
		TemplateID: &template.ID,
		Channel:    template.Channel,
		Type:       template.EventType,
		Title:      title,
		Message:    message,
		Data:       string(dataJSON),
		Status:     "pending",
	}

	return s.SendNotification(notification)
}
