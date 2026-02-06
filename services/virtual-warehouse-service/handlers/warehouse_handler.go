package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/virtual-warehouse-service/models"
	"github.com/b2b-platform/virtual-warehouse-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(service *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// Inventory endpoints
func (h *WarehouseHandler) CreateInventory(c *gin.Context) {
	var inventory models.SharedInventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	inventory.TenantID = tenantID

	if err := h.service.CreateInventory(&inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

func (h *WarehouseHandler) ListInventory(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	inventory, err := h.service.ListInventory(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *WarehouseHandler) GetAvailable(c *gin.Context) {
	partID, err := uuid.Parse(c.Query("part_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid part_id"})
		return
	}

	quantity, _ := strconv.ParseFloat(c.DefaultQuery("quantity", "1"), 64)

	inventory, err := h.service.GetAvailable(partID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// Equipment group endpoints
func (h *WarehouseHandler) CreateGroup(c *gin.Context) {
	var group models.EquipmentGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	group.TenantID = tenantID

	if err := h.service.CreateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *WarehouseHandler) GetGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.service.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *WarehouseHandler) ListGroups(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	groups, err := h.service.ListGroups(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// Transfer endpoints
func (h *WarehouseHandler) CreateTransfer(c *gin.Context) {
	var transfer models.InterCompanyTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	transfer.FromTenantID = tenantID
	transfer.RequestedBy = userID
	transfer.Status = "pending"

	if err := h.service.CreateTransfer(&transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (h *WarehouseHandler) ListTransfers(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	transfers, err := h.service.ListTransfers(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

func (h *WarehouseHandler) ApproveTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.ApproveTransfer(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer approved"})
}

// Emergency sourcing endpoints
func (h *WarehouseHandler) CreateEmergencySourcing(c *gin.Context) {
	var sourcing models.EmergencySourcing
	if err := c.ShouldBindJSON(&sourcing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	sourcing.TenantID = tenantID
	sourcing.RequestedBy = userID
	sourcing.Status = "open"

	if err := h.service.CreateEmergencySourcing(&sourcing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sourcing)
}

func (h *WarehouseHandler) ListEmergencySourcing(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	status := c.Query("status")

	sourcing, err := h.service.ListEmergencySourcing(tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sourcing)
}
