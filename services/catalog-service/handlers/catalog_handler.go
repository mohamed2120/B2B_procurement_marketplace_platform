package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/catalog-service/models"
	"github.com/b2b-platform/catalog-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatalogHandler struct {
	service *service.CatalogService
}

func NewCatalogHandler(service *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

// Manufacturer endpoints
func (h *CatalogHandler) CreateManufacturer(c *gin.Context) {
	var manufacturer models.Manufacturer
	if err := c.ShouldBindJSON(&manufacturer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateManufacturer(&manufacturer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, manufacturer)
}

func (h *CatalogHandler) GetManufacturer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	manufacturer, err := h.service.GetManufacturer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, manufacturer)
}

func (h *CatalogHandler) ListManufacturers(c *gin.Context) {
	manufacturers, err := h.service.ListManufacturers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, manufacturers)
}

// Category endpoints
func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CatalogHandler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	category, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CatalogHandler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Part endpoints
func (h *CatalogHandler) CreatePart(c *gin.Context) {
	var part models.LibraryPart
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := auth.GetUserID(c)
	part.CreatedBy = userID
	part.Status = "pending"

	if err := h.service.CreatePart(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, part)
}

func (h *CatalogHandler) GetPart(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	part, err := h.service.GetPart(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, part)
}

func (h *CatalogHandler) ListParts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")

	parts, err := h.service.ListParts(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, parts)
}

func (h *CatalogHandler) ApprovePart(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.ApprovePart(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "part approved"})
}

func (h *CatalogHandler) RejectPart(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RejectPart(id, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "part rejected"})
}

func (h *CatalogHandler) GetPendingParts(c *gin.Context) {
	parts, err := h.service.GetPendingParts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, parts)
}

// Attribute endpoints
func (h *CatalogHandler) CreateAttribute(c *gin.Context) {
	var attribute models.Attribute
	if err := c.ShouldBindJSON(&attribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateAttribute(&attribute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attribute)
}

func (h *CatalogHandler) ListAttributes(c *gin.Context) {
	attributes, err := h.service.ListAttributes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attributes)
}
