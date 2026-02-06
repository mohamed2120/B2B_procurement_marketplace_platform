package handlers

import (
	"net/http"

	"github.com/b2b-platform/company-service/models"
	"github.com/b2b-platform/company-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CompanyHandler struct {
	service *service.CompanyService
}

func NewCompanyHandler(service *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		LegalName   string `json:"legal_name"`
		TaxID       string `json:"tax_id"`
		Address     string `json:"address"`
		City        string `json:"city"`
		State       string `json:"state"`
		Country     string `json:"country"`
		PostalCode  string `json:"postal_code"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Website     string `json:"website"`
		Industry    string `json:"industry"`
		CompanyType string `json:"company_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := &models.Company{
		Name:        req.Name,
		LegalName:   req.LegalName,
		TaxID:       req.TaxID,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		Industry:    req.Industry,
		CompanyType: req.CompanyType,
		Status:      "pending",
	}

	if err := h.service.Create(company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *CompanyHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	company, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) List(c *gin.Context) {
	companies, err := h.service.List(100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	company, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.Approve(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company approved"})
}

func (h *CompanyHandler) RequestSubdomain(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Subdomain string `json:"subdomain" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.service.RequestSubdomain(id, req.Subdomain, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subdomain request created"})
}
