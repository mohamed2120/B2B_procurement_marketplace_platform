package errors

import (
	"encoding/json"
	"time"
)

type APIError struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func NewAPIError(code, requestID string, details map[string]interface{}) *APIError {
	return &APIError{
		Code:      code,
		Message:   GetErrorMessage(code),
		RequestID: requestID,
		Timestamp: time.Now().UTC(),
		Details:   details,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Helper to extract error code from database errors
func ExtractDBErrorCode(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()
	
	// PostgreSQL error patterns
	if contains(errStr, "duplicate key") || contains(errStr, "unique constraint") {
		return ErrorCodeDBUniqueViolation
	}
	if contains(errStr, "foreign key constraint") || contains(errStr, "violates foreign key") {
		return ErrorCodeDBForeignKeyViolation
	}
	if contains(errStr, "connection") || contains(errStr, "connect") {
		return ErrorCodeDBConnectionFailed
	}
	if contains(errStr, "no rows") || contains(errStr, "not found") {
		return ErrorCodeDBRecordNotFound
	}

	return ErrorCodeDBQueryFailed
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
