package handlers

import (
	"net/http"

	"github.com/b2b-platform/identity-service/service"
	"github.com/b2b-platform/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListUsers lists all users (admin only)
// Query params: tenant_id (optional filter)
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Check if admin role
	roles, exists := c.Get(auth.RolesKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rolesList, ok := roles.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isAdmin := false
	for _, role := range rolesList {
		if role == "admin" || role == "super_admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	// Optional tenant filter
	tenantIDStr := c.Query("tenant_id")
	var tenantID *uuid.UUID
	if tenantIDStr != "" {
		id, err := uuid.Parse(tenantIDStr)
		if err == nil {
			tenantID = &id
		}
	}

	users, err := h.userService.List(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	// Check if admin or same user
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Get current user ID from context
	currentUserID, exists := c.Get(auth.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if admin or same user
	roles, _ := c.Get(auth.RolesKey)
	rolesList, ok := roles.([]string)
	isAdmin := false
	if ok {
		for _, role := range rolesList {
			if role == "admin" || role == "super_admin" {
				isAdmin = true
				break
			}
		}
	}

	if !isAdmin && currentUserID.(uuid.UUID) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ToggleActive toggles user active status (admin only)
func (h *UserHandler) ToggleActive(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Check admin role
	roles, exists := c.Get(auth.RolesKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rolesList, ok := roles.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isAdmin := false
	for _, role := range rolesList {
		if role == "admin" || role == "super_admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ToggleActive(userID, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
