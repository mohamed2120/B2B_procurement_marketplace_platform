package handlers

import (
	"net/http"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	payoutService  *service.PayoutService
}

func NewPaymentHandler(paymentService *service.PaymentService, payoutService *service.PayoutService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		payoutService:  payoutService,
	}
}

// CreatePaymentIntent creates a payment intent for an order
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req service.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)

	response, err := h.paymentService.CreatePaymentIntent(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// HandleWebhook processes payment provider webhooks
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// Read raw body for signature verification
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read payload"})
		return
	}

	signature := c.GetHeader("X-Signature")
	if signature == "" {
		signature = c.GetHeader("Stripe-Signature") // For Stripe compatibility
	}

	if err := h.paymentService.HandleWebhook(c.Request.Context(), payload, signature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// ReleaseEscrow releases held escrow funds to supplier
func (h *PaymentHandler) ReleaseEscrow(c *gin.Context) {
	var req struct {
		EscrowHoldID uuid.UUID `json:"escrow_hold_id" binding:"required"`
		Reason       string    `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := auth.GetUserID(c)

	if err := h.paymentService.ReleaseEscrow(c.Request.Context(), req.EscrowHoldID, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "escrow released successfully"})
}

// CreateRefund issues a refund
func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	var req struct {
		PaymentID uuid.UUID `json:"payment_id" binding:"required"`
		Amount    float64   `json:"amount" binding:"required"`
		Reason    string    `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	userID, _ := auth.GetUserID(c)

	refund, err := h.paymentService.CreateRefund(c.Request.Context(), tenantID, req.PaymentID, req.Amount, req.Reason, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, refund)
}

// GetPayment gets a payment by ID
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	payment, err := h.paymentService.GetPayment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ListPayments lists payments for tenant
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	limit := 100
	offset := 0

	payments, err := h.paymentService.ListPayments(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetEscrowHold gets an escrow hold by ID
func (h *PaymentHandler) GetEscrowHold(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	hold, err := h.paymentService.GetEscrowHold(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, hold)
}

// ListEscrowHolds lists escrow holds for supplier
func (h *PaymentHandler) ListEscrowHolds(c *gin.Context) {
	supplierIDStr := c.Query("supplier_id")
	if supplierIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id query parameter required"})
		return
	}

	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier_id"})
		return
	}

	limit := 100
	offset := 0

	holds, err := h.paymentService.ListEscrowHolds(supplierID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, holds)
}

// PayoutAccount CRUD handlers
func (h *PaymentHandler) CreatePayoutAccount(c *gin.Context) {
	var account models.PayoutAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	account.TenantID = tenantID

	if err := h.payoutService.CreatePayoutAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *PaymentHandler) GetPayoutAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.payoutService.GetPayoutAccount(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *PaymentHandler) ListPayoutAccounts(c *gin.Context) {
	supplierIDStr := c.Query("supplier_id")
	if supplierIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id query parameter required"})
		return
	}

	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier_id"})
		return
	}

	accounts, err := h.payoutService.ListPayoutAccounts(supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *PaymentHandler) UpdatePayoutAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var account models.PayoutAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.ID = id

	if err := h.payoutService.UpdatePayoutAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *PaymentHandler) DeletePayoutAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.payoutService.DeletePayoutAccount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payout account deleted"})
}

func (h *PaymentHandler) SetDefaultPayoutAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get supplier_id from query or account
	supplierIDStr := c.Query("supplier_id")
	if supplierIDStr == "" {
		account, err := h.payoutService.GetPayoutAccount(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id required"})
			return
		}
		supplierIDStr = account.SupplierID.String()
	}

	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier_id"})
		return
	}

	if err := h.payoutService.SetDefaultPayoutAccount(supplierID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "default payout account set"})
}
