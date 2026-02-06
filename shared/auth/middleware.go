package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey   = "user_id"
	TenantIDKey = "tenant_id"
	EmailKey    = "email"
	RolesKey    = "roles"
)

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Set context values
		c.Set(UserIDKey, claims.UserID)
		c.Set(TenantIDKey, claims.TenantID)
		c.Set(EmailKey, claims.Email)
		c.Set(RolesKey, claims.Roles)

		c.Next()
	}
}

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get tenant from subdomain
		host := c.GetHeader("Host")
		if host != "" {
			parts := strings.Split(host, ".")
			if len(parts) > 0 {
				_ = parts[0] // subdomain - reserved for future use
				// In production, resolve subdomain to tenant_id
				// For now, we'll use the tenant_id from JWT
			}
		}

		// Fallback to JWT tenant_id (set by AuthMiddleware)
		tenantID, exists := c.Get(TenantIDKey)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id required"})
			c.Abort()
			return
		}

		c.Set(TenantIDKey, tenantID)
		c.Next()
	}
}

func GetTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return uuid.Nil, ErrInvalidToken
	}

	id, ok := tenantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return id, nil
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, ErrInvalidToken
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return id, nil
}
