package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/collaboration-service/models"
	"github.com/b2b-platform/collaboration-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollaborationHandler struct {
	service *service.CollaborationService
}

func NewCollaborationHandler(service *service.CollaborationService) *CollaborationHandler {
	return &CollaborationHandler{service: service}
}

// Thread endpoints
func (h *CollaborationHandler) CreateThread(c *gin.Context) {
	var thread models.ChatThread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	thread.TenantID = tenantID
	thread.CreatedBy = userID

	if err := h.service.CreateThread(&thread); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, thread)
}

func (h *CollaborationHandler) GetThread(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	thread, err := h.service.GetThread(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (h *CollaborationHandler) ListThreads(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	threads, err := h.service.ListThreads(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threads)
}

func (h *CollaborationHandler) GetUserThreads(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	threads, err := h.service.GetUserThreads(userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threads)
}

// Message endpoints
func (h *CollaborationHandler) SendMessage(c *gin.Context) {
	threadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	var message models.ChatMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := auth.GetUserID(c)
	message.ThreadID = threadID
	message.SenderID = userID
	message.MessageType = "text"

	if err := h.service.SendMessage(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *CollaborationHandler) GetThreadMessages(c *gin.Context) {
	threadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.service.GetThreadMessages(threadID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// Dispute endpoints
func (h *CollaborationHandler) CreateDispute(c *gin.Context) {
	var dispute models.Dispute
	if err := c.ShouldBindJSON(&dispute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	dispute.TenantID = tenantID
	dispute.RaisedBy = userID
	dispute.Status = "open"

	if err := h.service.CreateDispute(&dispute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dispute)
}

func (h *CollaborationHandler) GetDispute(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dispute, err := h.service.GetDispute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, dispute)
}

func (h *CollaborationHandler) ListDisputes(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	status := c.Query("status")

	disputes, err := h.service.ListDisputes(tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, disputes)
}

func (h *CollaborationHandler) ResolveDispute(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Resolution string `json:"resolution" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.ResolveDispute(id, userID, req.Resolution); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dispute resolved"})
}

// Rating endpoints
func (h *CollaborationHandler) CreateRating(c *gin.Context) {
	var rating models.Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	rating.TenantID = tenantID
	rating.RatedBy = userID

	if err := h.service.CreateRating(&rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rating)
}

func (h *CollaborationHandler) GetRatings(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID, err := uuid.Parse(c.Query("entity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}

	ratings, err := h.service.GetRatings(entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func (h *CollaborationHandler) GetAverageRating(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID, err := uuid.Parse(c.Query("entity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}

	average, err := h.service.GetAverageRating(entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"average_rating": average})
}
