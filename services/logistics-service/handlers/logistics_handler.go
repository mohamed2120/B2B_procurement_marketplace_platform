package handlers

import (
	"net/http"

	"github.com/b2b-platform/logistics-service/models"
	"github.com/b2b-platform/logistics-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogisticsHandler struct {
	service *service.LogisticsService
}

func NewLogisticsHandler(service *service.LogisticsService) *LogisticsHandler {
	return &LogisticsHandler{service: service}
}

func (h *LogisticsHandler) Create(c *gin.Context) {
	var shipment models.Shipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	shipment.TenantID = tenantID

	if err := h.service.Create(&shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

func (h *LogisticsHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	shipment, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *LogisticsHandler) List(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	shipments, err := h.service.List(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shipments)
}

func (h *LogisticsHandler) UpdateTracking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var event models.TrackingEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ShipmentID = id

	if err := h.service.UpdateTracking(id, &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tracking updated"})
}
