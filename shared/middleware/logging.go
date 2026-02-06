package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/b2b-platform/shared/observability"
)

func RequestLogging(logger *observability.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or generate request ID
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		// Create request logger
		reqLogger := logger.WithRequest(requestID, c.FullPath(), c.Request.Method)
		
		// Extract tenant and user from context if available
		if tenantID, exists := c.Get("tenant_id"); exists {
			if tid, ok := tenantID.(string); ok {
				reqLogger.SetTenantID(tid)
			}
		}
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(string); ok {
				reqLogger.SetUserID(uid)
			}
		}

		start := time.Now()
		reqLogger.LogStart()

		// Process request
		c.Next()

		// Log completion
		statusCode := c.Writer.Status()
		if statusCode >= 500 {
			// Extract error code if available
			errorCode := ""
			if errCode, exists := c.Get("error_code"); exists {
				if ec, ok := errCode.(string); ok {
					errorCode = ec
				}
			}
			errorMessage := ""
			if errMsg, exists := c.Get("error_message"); exists {
				if em, ok := errMsg.(string); ok {
					errorMessage = em
				}
			}
			reqLogger.LogError(statusCode, errorCode, errorMessage)
		} else {
			reqLogger.LogEnd(statusCode)
		}
	}
}
