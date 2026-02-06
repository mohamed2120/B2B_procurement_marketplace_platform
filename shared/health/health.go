package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type HealthChecker struct {
	db          *gorm.DB
	redisClient *redis.Client
	serviceName string
}

func NewHealthChecker(serviceName string, db *gorm.DB, redisClient *redis.Client) *HealthChecker {
	return &HealthChecker{
		db:          db,
		redisClient: redisClient,
		serviceName: serviceName,
	}
}

// Health endpoint - just checks if process is alive
func (h *HealthChecker) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": h.serviceName,
	})
}

// Ready endpoint - checks all dependencies
func (h *HealthChecker) Ready(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	// Check database
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = "error: " + err.Error()
			allHealthy = false
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			
			if err := sqlDB.PingContext(ctx); err != nil {
				checks["database"] = "unhealthy: " + err.Error()
				allHealthy = false
			} else {
				checks["database"] = "healthy"
			}
		}
	}

	// Check Redis
	if h.redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		if err := h.redisClient.Ping(ctx).Err(); err != nil {
			checks["redis"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["redis"] = "healthy"
		}
	}

	if allHealthy {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": h.serviceName,
			"checks":  checks,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"service": h.serviceName,
			"checks":  checks,
		})
	}
}
