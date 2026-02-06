package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type ServiceHeartbeat struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ServiceName string    `gorm:"type:varchar(100);not null" json:"service_name"`
	InstanceID  string    `gorm:"type:varchar(255);not null" json:"instance_id"`
	LastSeenAt  time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"last_seen_at"`
	Status      string    `gorm:"type:varchar(20);not null;default:'healthy'" json:"status"`
	Version     string    `gorm:"type:varchar(50)" json:"version"`
	Env         string    `gorm:"type:varchar(20);not null" json:"env"`
	CreatedAt   time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

func (ServiceHeartbeat) TableName() string {
	return "diagnostics.service_heartbeats"
}

type Incident struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Severity        string         `gorm:"type:varchar(20);not null" json:"severity"`
	ServiceName     string         `gorm:"type:varchar(100);not null" json:"service_name"`
	Category        string         `gorm:"type:varchar(50);not null" json:"category"`
	ErrorCode       *string        `gorm:"type:varchar(50)" json:"error_code,omitempty"`
	Title           string         `gorm:"type:varchar(255);not null" json:"title"`
	DetailsJSON     JSONB          `gorm:"type:jsonb" json:"details_json,omitempty"`
	RequestID       *string        `gorm:"type:varchar(255)" json:"request_id,omitempty"`
	TenantID        *string        `gorm:"type:uuid" json:"tenant_id,omitempty"`
	UserID          *string        `gorm:"type:uuid" json:"user_id,omitempty"`
	OccurredAt      time.Time      `gorm:"type:timestamp with time zone;not null;default:now()" json:"occurred_at"`
	ResolvedAt      *time.Time     `gorm:"type:timestamp with time zone" json:"resolved_at,omitempty"`
	ResolutionNotes *string        `gorm:"type:text" json:"resolution_notes,omitempty"`
	CreatedAt       time.Time      `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

func (Incident) TableName() string {
	return "diagnostics.incidents"
}

type EventFailure struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventName     string    `gorm:"type:varchar(100);not null" json:"event_name"`
	Direction     string    `gorm:"type:varchar(20);not null" json:"direction"`
	ServiceName   string    `gorm:"type:varchar(100);not null" json:"service_name"`
	PayloadJSON   JSONB     `gorm:"type:jsonb" json:"payload_json,omitempty"`
	ErrorMessage  *string   `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount    int       `gorm:"not null;default:0" json:"retry_count"`
	LastAttemptAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"last_attempt_at"`
	Status        string    `gorm:"type:varchar(20);not null;default:'FAILED'" json:"status"`
	CreatedAt     time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

func (EventFailure) TableName() string {
	return "diagnostics.event_failures"
}

type Job struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobName      string    `gorm:"type:varchar(100);not null" json:"job_name"`
	ServiceName  string    `gorm:"type:varchar(100);not null" json:"service_name"`
	Status       string    `gorm:"type:varchar(20);not null" json:"status"`
	StartedAt    time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"started_at"`
	FinishedAt   *time.Time `gorm:"type:timestamp with time zone" json:"finished_at,omitempty"`
	ErrorMessage *string   `gorm:"type:text" json:"error_message,omitempty"`
	MetadataJSON JSONB     `gorm:"type:jsonb" json:"metadata_json,omitempty"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

func (Job) TableName() string {
	return "diagnostics.jobs"
}

type APIMetricMinute struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MinuteTS   time.Time `gorm:"type:timestamp with time zone;not null" json:"minute_ts"`
	ServiceName string   `gorm:"type:varchar(100);not null" json:"service_name"`
	Route      string    `gorm:"type:varchar(255);not null" json:"route"`
	Method     string    `gorm:"type:varchar(10);not null" json:"method"`
	CountTotal int       `gorm:"not null;default:0" json:"count_total"`
	Count2xx   int       `gorm:"not null;default:0" json:"count_2xx"`
	Count4xx   int       `gorm:"not null;default:0" json:"count_4xx"`
	Count5xx   int       `gorm:"not null;default:0" json:"count_5xx"`
	P95MS      *int      `gorm:"type:integer" json:"p95_ms,omitempty"`
	AvgMS      *int      `gorm:"type:integer" json:"avg_ms,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
}

func (APIMetricMinute) TableName() string {
	return "diagnostics.api_metrics_minute"
}
