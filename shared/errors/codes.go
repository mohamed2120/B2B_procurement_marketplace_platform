package errors

// Error codes follow pattern: CATEGORY_NUMBER
// Categories:
// - AUTH: Authentication/Authorization
// - TENANT: Tenant-related
// - DB: Database
// - EVENT: Event bus
// - FILE: File storage
// - SEARCH: Search/indexing
// - API: API/HTTP
// - BILLING: Billing/payments
// - LOGISTICS: Logistics/shipments
// - VALIDATION: Validation errors
// - INTERNAL: Internal server errors

const (
	// AUTH errors (001-099)
	ErrorCodeAuthInvalidToken     = "AUTH_001"
	ErrorCodeAuthTokenExpired     = "AUTH_002"
	ErrorCodeAuthMissingToken     = "AUTH_003"
	ErrorCodeAuthInvalidCredentials = "AUTH_004"
	ErrorCodeAuthUnauthorized     = "AUTH_005"
	ErrorCodeAuthForbidden        = "AUTH_006"

	// TENANT errors (001-099)
	ErrorCodeTenantNotFound       = "TENANT_001"
	ErrorCodeTenantMismatch       = "TENANT_002"
	ErrorCodeTenantMissing        = "TENANT_003"
	ErrorCodeTenantInvalid        = "TENANT_004"

	// DB errors (001-099)
	ErrorCodeDBConnectionFailed    = "DB_001"
	ErrorCodeDBQueryFailed        = "DB_002"
	ErrorCodeDBTransactionFailed   = "DB_003"
	ErrorCodeDBUniqueViolation     = "DB_004"
	ErrorCodeDBForeignKeyViolation = "DB_005"
	ErrorCodeDBRecordNotFound      = "DB_006"

	// EVENT errors (001-099)
	ErrorCodeEventPublishFailed   = "EVENT_001"
	ErrorCodeEventSubscribeFailed = "EVENT_002"
	ErrorCodeEventDeserializeFailed = "EVENT_003"
	ErrorCodeEventSerializeFailed = "EVENT_004"

	// FILE errors (001-099)
	ErrorCodeFileUploadFailed     = "FILE_001"
	ErrorCodeFileDownloadFailed   = "FILE_002"
	ErrorCodeFileDeleteFailed     = "FILE_003"
	ErrorCodeFileNotFound         = "FILE_004"
	ErrorCodeFileTooLarge         = "FILE_005"
	ErrorCodeFileInvalidType       = "FILE_006"

	// SEARCH errors (001-099)
	ErrorCodeSearchIndexFailed    = "SEARCH_001"
	ErrorCodeSearchQueryFailed    = "SEARCH_002"
	ErrorCodeSearchConnectionFailed = "SEARCH_003"

	// API errors (001-099)
	ErrorCodeAPIBadRequest        = "API_001"
	ErrorCodeAPINotFound          = "API_002"
	ErrorCodeAPIMethodNotAllowed  = "API_003"
	ErrorCodeAPIInternalError     = "API_004"
	ErrorCodeAPIServiceUnavailable = "API_005"
	ErrorCodeAPITimeout           = "API_006"

	// BILLING errors (001-099)
	ErrorCodeBillingPaymentFailed = "BILLING_001"
	ErrorCodeBillingSubscriptionNotFound = "BILLING_002"
	ErrorCodeBillingPlanNotFound  = "BILLING_003"

	// LOGISTICS errors (001-099)
	ErrorCodeLogisticsShipmentNotFound = "LOGISTICS_001"
	ErrorCodeLogisticsTrackingFailed   = "LOGISTICS_002"

	// VALIDATION errors (001-099)
	ErrorCodeValidationFailed     = "VALIDATION_001"
	ErrorCodeValidationRequired    = "VALIDATION_002"
	ErrorCodeValidationInvalid     = "VALIDATION_003"

	// INTERNAL errors (001-099)
	ErrorCodeInternalError         = "INTERNAL_001"
	ErrorCodeInternalPanic         = "INTERNAL_002"
)

var ErrorMessages = map[string]string{
	ErrorCodeAuthInvalidToken:     "Invalid authentication token",
	ErrorCodeAuthTokenExpired:     "Authentication token has expired",
	ErrorCodeAuthMissingToken:     "Authentication token is required",
	ErrorCodeAuthInvalidCredentials: "Invalid credentials",
	ErrorCodeAuthUnauthorized:     "Unauthorized access",
	ErrorCodeAuthForbidden:        "Forbidden: insufficient permissions",

	ErrorCodeTenantNotFound:       "Tenant not found",
	ErrorCodeTenantMismatch:       "Tenant mismatch",
	ErrorCodeTenantMissing:        "Tenant ID is required",
	ErrorCodeTenantInvalid:        "Invalid tenant ID",

	ErrorCodeDBConnectionFailed:    "Database connection failed",
	ErrorCodeDBQueryFailed:          "Database query failed",
	ErrorCodeDBTransactionFailed:   "Database transaction failed",
	ErrorCodeDBUniqueViolation:      "Unique constraint violation",
	ErrorCodeDBForeignKeyViolation: "Foreign key constraint violation",
	ErrorCodeDBRecordNotFound:       "Record not found",

	ErrorCodeEventPublishFailed:     "Failed to publish event",
	ErrorCodeEventSubscribeFailed:   "Failed to subscribe to event",
	ErrorCodeEventDeserializeFailed: "Failed to deserialize event",
	ErrorCodeEventSerializeFailed:   "Failed to serialize event",

	ErrorCodeFileUploadFailed:   "File upload failed",
	ErrorCodeFileDownloadFailed: "File download failed",
	ErrorCodeFileDeleteFailed:   "File deletion failed",
	ErrorCodeFileNotFound:       "File not found",
	ErrorCodeFileTooLarge:       "File size exceeds limit",
	ErrorCodeFileInvalidType:     "Invalid file type",

	ErrorCodeSearchIndexFailed:      "Search indexing failed",
	ErrorCodeSearchQueryFailed:      "Search query failed",
	ErrorCodeSearchConnectionFailed: "Search service connection failed",

	ErrorCodeAPIBadRequest:         "Bad request",
	ErrorCodeAPINotFound:           "Resource not found",
	ErrorCodeAPIMethodNotAllowed:   "Method not allowed",
	ErrorCodeAPIInternalError:      "Internal server error",
	ErrorCodeAPIServiceUnavailable: "Service unavailable",
	ErrorCodeAPITimeout:            "Request timeout",

	ErrorCodeBillingPaymentFailed:        "Payment processing failed",
	ErrorCodeBillingSubscriptionNotFound: "Subscription not found",
	ErrorCodeBillingPlanNotFound:        "Plan not found",

	ErrorCodeLogisticsShipmentNotFound: "Shipment not found",
	ErrorCodeLogisticsTrackingFailed:   "Tracking information unavailable",

	ErrorCodeValidationFailed:  "Validation failed",
	ErrorCodeValidationRequired: "Required field missing",
	ErrorCodeValidationInvalid:  "Invalid value",

	ErrorCodeInternalError: "Internal error occurred",
	ErrorCodeInternalPanic: "Internal panic occurred",
}

func GetErrorMessage(code string) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}
