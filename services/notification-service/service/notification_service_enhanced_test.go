package service

import (
	"testing"

	"github.com/b2b-platform/notification-service/models"
	"github.com/google/uuid"
)

func TestNotificationService_CreateTemplate(t *testing.T) {
	mockTemplateRepo := &MockTemplateRepository{
		templates: make(map[string]*models.NotificationTemplate),
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: make(map[string]bool),
	}
	mockNotifRepo := &MockNotificationRepository{}

	service := NewNotificationService(mockTemplateRepo, mockPrefRepo, mockNotifRepo)

	template := &models.NotificationTemplate{
		Code:      "test_template",
		Subject:   "Test Subject",
		Body:      "Test Body",
		Channel:   "in_app",
		EventType: "test.event.v1",
	}

	err := service.CreateTemplate(template)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify template was created
	retrieved, _ := service.GetTemplateByCode("test_template")
	if retrieved == nil {
		t.Errorf("expected template to be created")
	}
	if retrieved.Subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got %s", retrieved.Subject)
	}
}

func TestNotificationService_SendNotification_WithPreferenceDisabled(t *testing.T) {
	mockTemplateRepo := &MockTemplateRepository{
		templates: make(map[string]*models.NotificationTemplate),
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: map[string]bool{
			"in_app:test.event.v1": false, // Disabled
		},
	}
	mockNotifRepo := &MockNotificationRepository{}

	service := NewNotificationService(mockTemplateRepo, mockPrefRepo, mockNotifRepo)

	userID := uuid.New()
	tenantID := uuid.New()
	notification := &models.Notification{
		ID:       uuid.New(),
		UserID:   userID,
		TenantID: tenantID,
		Type:     "test.event.v1",
		Channel:  "in_app",
		Title:    "Test Notification",
		Status:   "pending",
	}

	err := service.SendNotification(notification)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify notification was created with skipped status
	created, _ := mockNotifRepo.GetByID(notification.ID)
	if created == nil {
		t.Errorf("expected notification to be created")
	}
	if created.Status != "skipped" {
		t.Errorf("expected status 'skipped', got %s", created.Status)
	}
}

func TestNotificationService_SendNotification_WithPreferenceEnabled(t *testing.T) {
	mockTemplateRepo := &MockTemplateRepository{
		templates: make(map[string]*models.NotificationTemplate),
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: map[string]bool{
			"in_app:test.event.v1": true, // Enabled
		},
	}
	mockNotifRepo := &MockNotificationRepository{}

	service := NewNotificationService(mockTemplateRepo, mockPrefRepo, mockNotifRepo)

	userID := uuid.New()
	tenantID := uuid.New()
	notification := &models.Notification{
		ID:       uuid.New(),
		UserID:   userID,
		TenantID: tenantID,
		Type:     "test.event.v1",
		Channel:  "in_app",
		Title:    "Test Notification",
		Status:   "pending",
	}

	err := service.SendNotification(notification)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify notification was created with sent status
	created, _ := mockNotifRepo.GetByID(notification.ID)
	if created == nil {
		t.Errorf("expected notification to be created")
	}
	if created.Status != "sent" {
		t.Errorf("expected status 'sent', got %s", created.Status)
	}
	if created.SentAt == nil {
		t.Errorf("expected SentAt to be set")
	}
}

func TestNotificationService_GetUserNotifications(t *testing.T) {
	mockTemplateRepo := &MockTemplateRepository{
		templates: make(map[string]*models.NotificationTemplate),
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: make(map[string]bool),
	}
	mockNotifRepo := &MockNotificationRepository{}

	service := NewNotificationService(mockTemplateRepo, mockPrefRepo, mockNotifRepo)

	userID := uuid.New()
	tenantID := uuid.New()

	// Create some notifications
	notif1 := &models.Notification{
		ID:       uuid.New(),
		UserID:   userID,
		TenantID: tenantID,
		Title:    "Notification 1",
		IsRead:   false,
	}
	notif2 := &models.Notification{
		ID:       uuid.New(),
		UserID:   userID,
		TenantID: tenantID,
		Title:    "Notification 2",
		IsRead:   true,
	}
	mockNotifRepo.Create(notif1)
	mockNotifRepo.Create(notif2)

	// Test getting all notifications
	notifications, err := service.GetUserNotifications(userID, tenantID, 10, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(notifications) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(notifications))
	}

	// Test getting only unread notifications
	unread, err := service.GetUserNotifications(userID, tenantID, 10, 0, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(unread) != 1 {
		t.Errorf("expected 1 unread notification, got %d", len(unread))
	}
	if unread[0].ID != notif1.ID {
		t.Errorf("expected unread notification ID %s, got %s", notif1.ID, unread[0].ID)
	}
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	mockTemplateRepo := &MockTemplateRepository{
		templates: make(map[string]*models.NotificationTemplate),
	}
	mockPrefRepo := &MockPreferenceRepository{
		preferences: make(map[string]bool),
	}
	mockNotifRepo := &MockNotificationRepository{}

	service := NewNotificationService(mockTemplateRepo, mockPrefRepo, mockNotifRepo)

	notificationID := uuid.New()
	notification := &models.Notification{
		ID:     notificationID,
		IsRead: false,
	}
	mockNotifRepo.Create(notification)

	err := service.MarkAsRead(notificationID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify notification was marked as read
	updated, _ := mockNotifRepo.GetByID(notificationID)
	if !updated.IsRead {
		t.Errorf("expected notification to be marked as read")
	}
	if updated.ReadAt == nil {
		t.Errorf("expected ReadAt to be set")
	}
}
