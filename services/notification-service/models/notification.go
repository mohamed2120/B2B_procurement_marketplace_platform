package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string    `gorm:"type:varchar(100);not null;unique" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Subject     string    `gorm:"type:varchar(255)" json:"subject"`
	Body        string    `gorm:"type:text;not null" json:"body"`
	BodyHTML    string    `gorm:"type:text" json:"body_html"`
	Channel     string    `gorm:"type:varchar(50);not null" json:"channel"` // email, in_app, sms, push
	EventType   string    `gorm:"type:varchar(100)" json:"event_type"` // Event that triggers this template
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Notifications []Notification `gorm:"foreignKey:TemplateID" json:"notifications,omitempty"`
}

func (NotificationTemplate) TableName() string {
	return "notification.templates"
}

type NotificationPreference struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Channel     string    `gorm:"type:varchar(50);not null;index" json:"channel"` // email, in_app, sms, push
	EventType   string    `gorm:"type:varchar(100);not null;index" json:"event_type"`
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
}

func (NotificationPreference) TableName() string {
	return "notification.preferences"
}

type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TemplateID  *uuid.UUID     `gorm:"type:uuid;index" json:"template_id,omitempty"`
	Channel     string         `gorm:"type:varchar(50);not null;index" json:"channel"` // email, in_app, sms, push
	Type        string         `gorm:"type:varchar(100);not null;index" json:"type"` // order_placed, shipment_late, etc.
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Data        string         `gorm:"type:jsonb" json:"data"` // Additional JSON data
	Status      string         `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, sent, failed, read
	IsRead      bool           `gorm:"default:false;index" json:"is_read"`
	ReadAt      *time.Time     `json:"read_at,omitempty"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	FailedAt    *time.Time     `json:"failed_at,omitempty"`
	Error       string         `gorm:"type:text" json:"error,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Template NotificationTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

func (Notification) TableName() string {
	return "notification.notifications"
}
