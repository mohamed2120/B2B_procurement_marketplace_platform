package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/b2b-platform/shared/diagnostics"
	"github.com/b2b-platform/shared/errors"
)

func ErrorHandler(reporter *diagnostics.Reporter, serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Extract API error if available
			var apiErr *errors.APIError
			if ae, ok := err.Err.(*errors.APIError); ok {
				apiErr = ae
			} else {
				// Create API error from generic error
				requestID, _ := c.Get("request_id")
				reqID := ""
				if rid, ok := requestID.(string); ok {
					reqID = rid
				}
				apiErr = errors.NewAPIError(errors.ErrorCodeAPIInternalError, reqID, map[string]interface{}{
					"error": err.Error(),
				})
			}

			// Set error code and message in context
			c.Set("error_code", apiErr.Code)
			c.Set("error_message", apiErr.Message)

			// Report incident for 5xx errors
			statusCode := c.Writer.Status()
			if statusCode >= 500 || statusCode == 0 {
				statusCode = http.StatusInternalServerError
				c.JSON(statusCode, apiErr)
			}

			// Extract tenant and user
			tenantID := ""
			userID := ""
			if tid, exists := c.Get("tenant_id"); exists {
				if t, ok := tid.(string); ok {
					tenantID = t
				}
			}
			if uid, exists := c.Get("user_id"); exists {
				if u, ok := uid.(string); ok {
					userID = u
				}
			}

			// Report incident
			if reporter != nil && (statusCode >= 500 || apiErr.Code != "") {
				reporter.ReportAPIError(serviceName, statusCode, apiErr, tenantID, userID)
			}
		}
	}
}
