package rbac

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type Permission struct {
	Resource string
	Action   string
}

func (p Permission) String() string {
	return fmt.Sprintf("%s:%s", p.Resource, p.Action)
}

type RBACService struct {
	redisClient *redis.Client
}

func NewRBACService(redisClient *redis.Client) *RBACService {
	return &RBACService{
		redisClient: redisClient,
	}
}

// CheckPermission verifies if a user has a specific permission
func (r *RBACService) CheckPermission(ctx context.Context, userID, tenantID string, permission Permission) (bool, error) {
	// Check Redis cache first
	cacheKey := fmt.Sprintf("permission:%s:%s:%s", userID, tenantID, permission.String())
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cached == "true" {
		return true, nil
	}

	// In a real implementation, check database
	// For now, we'll use a simple role-based check
	// This should be replaced with actual DB lookup

	return false, nil
}

// RequirePermission middleware that checks if user has required permission
func (r *RBACService) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		_, exists = c.Get("tenant_id")
		if !exists {
			c.JSON(400, gin.H{"error": "tenant_id required"})
			c.Abort()
			return
		}

		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(403, gin.H{"error": "no roles assigned"})
			c.Abort()
			return
		}

		// Check if user has permission through roles
		hasPermission := r.checkRolePermission(roles, resource, action)
		if !hasPermission {
			c.JSON(403, gin.H{"error": ErrPermissionDenied.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (r *RBACService) checkRolePermission(roles interface{}, resource, action string) bool {
	roleList, ok := roles.([]string)
	if !ok {
		return false
	}

	// Admin has all permissions
	for _, role := range roleList {
		if role == "admin" || role == "super_admin" {
			return true
		}
	}

	// Define role-based permissions
	permissions := map[string]map[string][]string{
		"procurement": {
			"create":  {"requester", "procurement_manager", "buyer", "admin"},
			"read":    {"requester", "procurement_manager", "buyer", "approver", "admin"},
			"update":  {"requester", "procurement_manager", "buyer", "admin"},
			"approve": {"procurement_manager", "approver", "admin"},
			"award":   {"procurement_manager", "admin"},
			"delete":  {"procurement_manager", "admin"},
		},
		"catalog": {
			"create": {"catalog_admin", "admin"},
			"read":   {"*"},
			"update": {"catalog_admin", "admin"},
			"delete": {"catalog_admin", "admin"},
		},
		"equipment": {
			"create": {"equipment_manager", "admin"},
			"read":   {"*"},
			"update": {"equipment_manager", "admin"},
			"delete": {"equipment_manager", "admin"},
		},
		"company": {
			"create": {"company_admin", "admin"},
			"read":   {"*"},
			"update": {"company_admin", "admin"},
			"delete": {"company_admin", "admin"},
		},
	}

	resourcePerms, exists := permissions[resource]
	if !exists {
		return false
	}

	allowedRoles, exists := resourcePerms[action]
	if !exists {
		return false
	}

	// Check if any role matches
	for _, role := range roleList {
		for _, allowedRole := range allowedRoles {
			if allowedRole == "*" || strings.EqualFold(role, allowedRole) {
				return true
			}
		}
	}

	return false
}
