package diagnostics

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/b2b-platform/shared/errors"
	"gorm.io/gorm"
)

type Reporter struct {
	db *gorm.DB
}

func NewReporter(db *gorm.DB) *Reporter {
	return &Reporter{db: db}
}

type IncidentSeverity string
type IncidentCategory string

const (
	SeverityInfo     IncidentSeverity = "INFO"
	SeverityWarn     IncidentSeverity = "WARN"
	SeverityError    IncidentSeverity = "ERROR"
	SeverityCritical IncidentSeverity = "CRITICAL"

	CategoryDB        IncidentCategory = "DB"
	CategoryAuth      IncidentCategory = "AUTH"
	CategoryEvent     IncidentCategory = "EVENT"
	CategoryAPI       IncidentCategory = "API"
	CategoryFile      IncidentCategory = "FILE"
	CategorySearch    IncidentCategory = "SEARCH"
	CategoryBilling   IncidentCategory = "BILLING"
	CategoryLogistics IncidentCategory = "LOGISTICS"
	CategoryOther     IncidentCategory = "OTHER"
)

type IncidentData struct {
	Severity     IncidentSeverity
	ServiceName  string
	Category     IncidentCategory
	ErrorCode    string
	Title        string
	Details      map[string]interface{}
	RequestID    string
	TenantID     string
	UserID       string
}

func (r *Reporter) ReportIncident(data IncidentData) error {
	detailsJSON, _ := json.Marshal(data.Details)
	
	var tenantID, userID, requestID *string
	if data.TenantID != "" {
		tenantID = &data.TenantID
	}
	if data.UserID != "" {
		userID = &data.UserID
	}
	if data.RequestID != "" {
		requestID = &data.RequestID
	}

	var errorCode *string
	if data.ErrorCode != "" {
		errorCode = &data.ErrorCode
	}

	query := `
		INSERT INTO diagnostics.incidents 
		(severity, service_name, category, error_code, title, details_json, request_id, tenant_id, user_id, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id
	`

	var id string
	err := r.db.Raw(query,
		string(data.Severity),
		data.ServiceName,
		string(data.Category),
		errorCode,
		data.Title,
		string(detailsJSON),
		requestID,
		tenantID,
		userID,
	).Scan(&id).Error

	return err
}

type EventFailureData struct {
	EventName    string
	Direction    string // "PUBLISH" or "CONSUME"
	ServiceName  string
	Payload      map[string]interface{}
	ErrorMessage string
}

func (r *Reporter) ReportEventFailure(data EventFailureData) error {
	payloadJSON, _ := json.Marshal(data.Payload)
	
	var errorMsg *string
	if data.ErrorMessage != "" {
		errorMsg = &data.ErrorMessage
	}

	query := `
		INSERT INTO diagnostics.event_failures 
		(event_name, direction, service_name, payload_json, error_message, last_attempt_at, status)
		VALUES ($1, $2, $3, $4, $5, NOW(), 'FAILED')
		RETURNING id
	`

	var id string
	err := r.db.Raw(query,
		data.EventName,
		data.Direction,
		data.ServiceName,
		string(payloadJSON),
		errorMsg,
	).Scan(&id).Error

	return err
}

func (r *Reporter) UpdateHeartbeat(serviceName, instanceID, status, version, env string) error {
	query := `
		INSERT INTO diagnostics.service_heartbeats 
		(service_name, instance_id, last_seen_at, status, version, env)
		VALUES ($1, $2, NOW(), $3, $4, $5)
		ON CONFLICT (service_name, instance_id) 
		DO UPDATE SET 
			last_seen_at = NOW(),
			status = $3,
			version = $4,
			updated_at = NOW()
	`

	return r.db.Exec(query, serviceName, instanceID, status, version, env).Error
}

func (r *Reporter) RecordJob(jobName, serviceName, status string, errorMsg *string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	
	var finishedAt *time.Time
	now := time.Now()
	if status == "SUCCESS" || status == "FAILED" {
		finishedAt = &now
	}

	query := `
		INSERT INTO diagnostics.jobs 
		(job_name, service_name, status, started_at, finished_at, error_message, metadata_json)
		VALUES ($1, $2, $3, NOW(), $4, $5, $6)
		RETURNING id
	`

	var id string
	err := r.db.Raw(query,
		jobName,
		serviceName,
		status,
		finishedAt,
		errorMsg,
		string(metadataJSON),
	).Scan(&id).Error

	return err
}

// Helper to automatically report incidents from API errors
func (r *Reporter) ReportAPIError(serviceName string, statusCode int, apiErr *errors.APIError, tenantID, userID string) {
	severity := SeverityError
	if statusCode >= 500 {
		severity = SeverityCritical
	} else if statusCode >= 400 {
		severity = SeverityWarn
	}

	category := CategoryAPI
	if apiErr.Code != "" {
		if apiErr.Code[:4] == "AUTH" {
			category = CategoryAuth
		} else if apiErr.Code[:2] == "DB" {
			category = CategoryDB
		} else if apiErr.Code[:5] == "EVENT" {
			category = CategoryEvent
		} else if apiErr.Code[:4] == "FILE" {
			category = CategoryFile
		} else if apiErr.Code[:6] == "SEARCH" {
			category = CategorySearch
		}
	}

	details := map[string]interface{}{
		"status_code": statusCode,
	}
	if apiErr.Details != nil {
		for k, v := range apiErr.Details {
			details[k] = v
		}
	}

	r.ReportIncident(IncidentData{
		Severity:    severity,
		ServiceName: serviceName,
		Category:    category,
		ErrorCode:   apiErr.Code,
		Title:       apiErr.Message,
		Details:     details,
		RequestID:   apiErr.RequestID,
		TenantID:    tenantID,
		UserID:      userID,
	})
}
