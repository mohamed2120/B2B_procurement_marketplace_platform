package handlers

import (
	"net/http"
	"time"

	"github.com/b2b-platform/identity-service/models"
	"github.com/b2b-platform/identity-service/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *service.UserService
	jwtService  *service.JWTService
}

func NewAuthHandler(userService *service.UserService, jwtService *service.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TenantID  string `json:"tenant_id" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      interface{} `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	user, err := h.userService.GetByEmail(tenantID, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get user roles
	roles, err := h.userService.GetUserRoles(user.ID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user roles"})
		return
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate token
	token, err := h.jwtService.GenerateToken(user.ID, tenantID, user.Email, roleNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	h.userService.Update(user)

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &models.User{
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		IsVerified:   false,
	}

	if err := h.userService.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// Token validation is handled by middleware
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	email, _ := c.Get("email")
	roles, _ := c.Get("roles")

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"tenant_id": tenantID,
		"email":     email,
		"roles":     roles,
	})
}
