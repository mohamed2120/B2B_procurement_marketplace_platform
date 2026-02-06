package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/notification-service/models"
	"github.com/b2b-platform/notification-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Template endpoints
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var template models.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *NotificationHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	template, err := h.service.GetTemplate(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *NotificationHandler) ListTemplates(c *gin.Context) {
	templates, err := h.service.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// Preference endpoints
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	preferences, err := h.service.GetUserPreferences(userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func (h *NotificationHandler) UpdatePreference(c *gin.Context) {
	var preference models.NotificationPreference
	if err := c.ShouldBindJSON(&preference); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	preference.TenantID = tenantID
	preference.UserID = userID

	if err := h.service.UpdatePreference(&preference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preference)
}

// Notification endpoints
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	notification.TenantID = tenantID

	if err := h.service.SendNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	unreadOnly := c.Query("unread_only") == "true"

	notifications, err := h.service.GetUserNotifications(userID, tenantID, limit, offset, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	if err := h.service.MarkAllAsRead(userID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}
