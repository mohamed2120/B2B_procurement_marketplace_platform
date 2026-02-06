package handlers

import (
	"net/http"
	"strconv"

	"github.com/b2b-platform/procurement-service/models"
	"github.com/b2b-platform/procurement-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProcurementHandler struct {
	service *service.ProcurementService
}

func NewProcurementHandler(service *service.ProcurementService) *ProcurementHandler {
	return &ProcurementHandler{service: service}
}

func (h *ProcurementHandler) CreatePR(c *gin.Context) {
	var pr models.PurchaseRequest
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	pr.TenantID = tenantID
	pr.RequestedBy = userID
	pr.Status = "draft"

	if err := h.service.CreatePR(&pr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pr)
}

func (h *ProcurementHandler) GetPR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	pr, err := h.service.GetPR(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, pr)
}

func (h *ProcurementHandler) ListPRs(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	prs, err := h.service.ListPRs(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prs)
}

func (h *ProcurementHandler) UpdatePR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	pr, err := h.service.GetPR(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update logic would go here
	c.JSON(http.StatusOK, pr)
}

func (h *ProcurementHandler) ApprovePR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.ApprovePR(id, uuid.New()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PR approved"})
}

func (h *ProcurementHandler) CreateRFQ(c *gin.Context) {
	var rfq models.RFQ
	if err := c.ShouldBindJSON(&rfq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	rfq.TenantID = tenantID
	rfq.CreatedBy = userID

	if err := h.service.CreateRFQ(&rfq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rfq)
}

func (h *ProcurementHandler) GetRFQ(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rfq, err := h.service.GetRFQ(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, rfq)
}

func (h *ProcurementHandler) SubmitQuote(c *gin.Context) {
	var quote models.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	quote.TenantID = tenantID

	if err := h.service.SubmitQuote(&quote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quote)
}

func (h *ProcurementHandler) GetQuote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	quote, err := h.service.GetQuote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (h *ProcurementHandler) CreatePO(c *gin.Context) {
	var po models.PurchaseOrder
	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	po.TenantID = tenantID
	po.CreatedBy = userID

	// Get auth token from header
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		// Try to get from context if available
		if token, exists := c.Get("token"); exists {
			authToken = token.(string)
		}
	}

	if err := h.service.CreatePO(&po, authToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, po)
}

func (h *ProcurementHandler) GetPO(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	po, err := h.service.GetPO(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, po)
}

func (h *ProcurementHandler) UpdatePOPaymentStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		PaymentStatus string `json:"payment_status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	po, err := h.service.GetPO(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	po.PaymentStatus = req.PaymentStatus
	if err := h.service.UpdatePO(po); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, po)
}
