package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/equipment-service/models"
	"github.com/b2b-platform/equipment-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EquipmentHandler struct {
	service *service.EquipmentService
}

func NewEquipmentHandler(service *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

func (h *EquipmentHandler) Create(c *gin.Context) {
	var equipment models.Equipment
	if err := c.ShouldBindJSON(&equipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	equipment.TenantID = tenantID

	if err := h.service.Create(&equipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

func (h *EquipmentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	equipment, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) List(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	equipment, err := h.service.List(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	equipment, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := c.ShouldBindJSON(&equipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(equipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) AddBOMNode(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var node models.BOMNode
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	node.TenantID = tenantID
	node.EquipmentID = id

	if err := h.service.AddBOMNode(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *EquipmentHandler) GetBOM(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	nodes, err := h.service.GetBOM(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodes)
}

func (h *EquipmentHandler) CreateCompatibilityMapping(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mapping models.CompatibilityMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	mapping.TenantID = tenantID
	mapping.EquipmentID = id

	if err := h.service.CreateCompatibilityMapping(&mapping); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapping)
}

func (h *EquipmentHandler) VerifyCompatibility(c *gin.Context) {
	mappingID, err := uuid.Parse(c.Param("mapping_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping_id"})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.VerifyCompatibility(mappingID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "compatibility verified"})
}

func (h *EquipmentHandler) CheckCompatibility(c *gin.Context) {
	equipmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment id"})
		return
	}

	partID, err := uuid.Parse(c.Query("part_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid part_id"})
		return
	}

	mapping, err := h.service.CheckCompatibility(equipmentID, partID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if mapping == nil {
		c.JSON(http.StatusOK, gin.H{"is_compatible": false, "message": "no mapping found"})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

func (h *EquipmentHandler) GetCompatibilityMappings(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mappings, err := h.service.GetCompatibilityMappings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappings)
}
