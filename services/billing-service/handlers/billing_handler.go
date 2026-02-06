package handlers

import (
	"net/http"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BillingHandler struct {
	service *service.BillingService
}

func NewBillingHandler(service *service.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

// Plan endpoints
func (h *BillingHandler) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *BillingHandler) GetPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	plan, err := h.service.GetPlan(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *BillingHandler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// Subscription endpoints
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	var subscription models.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := auth.GetTenantID(c)
	subscription.TenantID = tenantID

	if err := h.service.CreateSubscription(&subscription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *BillingHandler) GetSubscription(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)

	subscription, err := h.service.GetTenantSubscription(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.CancelSubscription(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription cancelled"})
}

func (h *BillingHandler) CheckEntitlement(c *gin.Context) {
	tenantID, _ := auth.GetTenantID(c)
	feature := c.Query("feature")

	if feature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feature parameter required"})
		return
	}

	hasAccess, limit, err := h.service.CheckEntitlement(tenantID, feature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_access": hasAccess,
		"limit":      limit,
		"feature":    feature,
	})
}
