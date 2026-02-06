package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"` // PO ID
	PaymentIntentID string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"payment_intent_id"`
	Provider        string     `gorm:"type:varchar(50);not null" json:"provider"` // stripe, mock, etc.
	Amount          float64    `gorm:"not null" json:"amount"`
	Currency        string     `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	Status          string     `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, succeeded, failed, cancelled
	PaymentMode     string     `gorm:"type:varchar(50);not null;index" json:"payment_mode"` // DIRECT, ESCROW
	Metadata        string     `gorm:"type:jsonb" json:"metadata,omitempty"` // Provider-specific data
	FailedReason    string     `gorm:"type:text" json:"failed_reason,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	EscrowHold *EscrowHold `gorm:"foreignKey:PaymentID" json:"escrow_hold,omitempty"`
	Refunds    []Refund    `gorm:"foreignKey:PaymentID" json:"refunds,omitempty"`
}

func (Payment) TableName() string {
	return "billing.payments"
}

type EscrowHold struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PaymentID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"payment_id"`
	OrderID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"` // PO ID
	SupplierID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"supplier_id"`
	Amount            float64    `gorm:"not null" json:"amount"`
	Currency          string     `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	Status            string     `gorm:"type:varchar(50);default:'held';index" json:"status"` // held, released, refunded
	AutoReleaseDays   int        `gorm:"default:30" json:"auto_release_days"`
	AutoReleaseDate   *time.Time `json:"auto_release_date,omitempty"`
	ReleasedAt        *time.Time `json:"released_at,omitempty"`
	ReleasedBy        *uuid.UUID `json:"released_by,omitempty"`
	ReleaseReason     string     `gorm:"type:text" json:"release_reason,omitempty"`
	BlockedByDispute  bool       `gorm:"default:false;index" json:"blocked_by_dispute"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relationships
	Payment Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

func (EscrowHold) TableName() string {
	return "billing.escrow_holds"
}

type Settlement struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EscrowHoldID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"escrow_hold_id"`
	SupplierID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"supplier_id"`
	PayoutAccountID uuid.UUID  `gorm:"type:uuid;index" json:"payout_account_id,omitempty"`
	Amount          float64    `gorm:"not null" json:"amount"`
	Currency        string     `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	Status          string     `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, processing, completed, failed
	ProviderPayoutID string    `gorm:"type:varchar(255)" json:"provider_payout_id,omitempty"`
	FailedReason    string     `gorm:"type:text" json:"failed_reason,omitempty"`
	SettledAt       *time.Time `json:"settled_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationships
	EscrowHold   EscrowHold   `gorm:"foreignKey:EscrowHoldID" json:"escrow_hold,omitempty"`
	PayoutAccount *PayoutAccount `gorm:"foreignKey:PayoutAccountID" json:"payout_account,omitempty"`
}

func (Settlement) TableName() string {
	return "billing.settlements"
}

type Refund struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PaymentID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"payment_id"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	RefundNumber    string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_refund" json:"refund_number"`
	Amount          float64    `gorm:"not null" json:"amount"`
	Currency        string     `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	Reason          string     `gorm:"type:text" json:"reason"`
	Status          string     `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, processing, completed, failed
	ProviderRefundID string    `gorm:"type:varchar(255)" json:"provider_refund_id,omitempty"`
	FailedReason    string     `gorm:"type:text" json:"failed_reason,omitempty"`
	RefundedAt      *time.Time `json:"refunded_at,omitempty"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationships
	Payment Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

func (Refund) TableName() string {
	return "billing.refunds"
}

type PayoutAccount struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SupplierID      uuid.UUID `gorm:"type:uuid;not null;index" json:"supplier_id"`
	AccountType     string    `gorm:"type:varchar(50);not null" json:"account_type"` // bank_account, stripe_account, etc.
	Provider        string    `gorm:"type:varchar(50);not null" json:"provider"` // stripe, bank, etc.
	AccountNumber   string    `gorm:"type:varchar(255)" json:"account_number,omitempty"` // Last 4 digits or masked
	RoutingNumber   string    `gorm:"type:varchar(50)" json:"routing_number,omitempty"`
	AccountHolderName string  `gorm:"type:varchar(255)" json:"account_holder_name"`
	BankName        string    `gorm:"type:varchar(255)" json:"bank_name,omitempty"`
	ProviderAccountID string  `gorm:"type:varchar(255);index" json:"provider_account_id,omitempty"`
	IsDefault       bool     `gorm:"default:false;index" json:"is_default"`
	IsVerified      bool     `gorm:"default:false" json:"is_verified"`
	Metadata        string   `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Settlements []Settlement `gorm:"foreignKey:PayoutAccountID" json:"settlements,omitempty"`
}

func (PayoutAccount) TableName() string {
	return "billing.payout_accounts"
}
