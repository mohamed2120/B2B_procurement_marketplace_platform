package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Code        string    `gorm:"type:varchar(100);not null;unique" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Currency    string    `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	BillingCycle string   `gorm:"type:varchar(50);not null" json:"billing_cycle"` // monthly, yearly
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Entitlements []Entitlement `gorm:"foreignKey:PlanID" json:"entitlements,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:PlanID" json:"subscriptions,omitempty"`
}

func (Plan) TableName() string {
	return "billing.plans"
}

type Entitlement struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID      uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	Feature     string    `gorm:"type:varchar(100);not null" json:"feature"`
	Limit       int       `json:"limit"` // -1 for unlimited
	Unit        string    `gorm:"type:varchar(50)" json:"unit"` // users, orders, storage_gb, etc.
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Plan Plan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

func (Entitlement) TableName() string {
	return "billing.entitlements"
}

type Subscription struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PlanID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"plan_id"`
	Status      string         `gorm:"type:varchar(50);default:'active';index" json:"status"` // active, cancelled, expired, suspended
	StartedAt   time.Time      `gorm:"not null" json:"started_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CancelledAt *time.Time     `json:"cancelled_at,omitempty"`
	AutoRenew   bool           `gorm:"default:true" json:"auto_renew"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Plan Plan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

func (Subscription) TableName() string {
	return "billing.subscriptions"
}
