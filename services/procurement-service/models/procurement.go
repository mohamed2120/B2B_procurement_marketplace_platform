package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseRequest struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PRNumber    string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_pr" json:"pr_number"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(50);default:'draft';index" json:"status"` // draft, submitted, approved, rejected, cancelled
	Priority    string         `gorm:"type:varchar(50);default:'normal'" json:"priority"` // low, normal, high, urgent
	RequestedBy uuid.UUID      `gorm:"type:uuid;not null" json:"requested_by"`
	Department  string         `gorm:"type:varchar(100)" json:"department"`
	Budget      float64        `json:"budget"`
	Currency    string         `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	RequiredDate *time.Time    `json:"required_date,omitempty"`
	ApprovedAt  *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy  *uuid.UUID     `json:"approved_by,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Items    []PRItem        `gorm:"foreignKey:PRID" json:"items,omitempty"`
	Approvals []PRApproval   `gorm:"foreignKey:PRID" json:"approvals,omitempty"`
	RFQs     []RFQ           `gorm:"foreignKey:PRID" json:"rfqs,omitempty"`
}

func (PurchaseRequest) TableName() string {
	return "procurement.purchase_requests"
}

type PRItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PRID        uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_id"`
	PartID      *uuid.UUID `gorm:"type:uuid;index" json:"part_id,omitempty"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Unit        string    `gorm:"type:varchar(50)" json:"unit"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	Specifications string `gorm:"type:text" json:"specifications"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PR PurchaseRequest `gorm:"foreignKey:PRID" json:"pr,omitempty"`
}

func (PRItem) TableName() string {
	return "procurement.pr_items"
}

type PRApproval struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PRID      uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_id"`
	ApproverID uuid.UUID `gorm:"type:uuid;not null" json:"approver_id"`
	Status    string    `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, approved, rejected
	Comments  string    `gorm:"type:text" json:"comments"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	PR PurchaseRequest `gorm:"foreignKey:PRID" json:"pr,omitempty"`
}

func (PRApproval) TableName() string {
	return "procurement.pr_approvals"
}

type RFQ struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PRID        uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_id"`
	RFQNumber   string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_rfq" json:"rfq_number"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);default:'draft';index" json:"status"` // draft, sent, closed, cancelled
	DueDate     time.Time `gorm:"not null" json:"due_date"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PR     PurchaseRequest `gorm:"foreignKey:PRID" json:"pr,omitempty"`
	Quotes []Quote         `gorm:"foreignKey:RFQID" json:"quotes,omitempty"`
}

func (RFQ) TableName() string {
	return "procurement.rfqs"
}

type Quote struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	RFQID       uuid.UUID `gorm:"type:uuid;not null;index" json:"rfq_id"`
	SupplierID  uuid.UUID `gorm:"type:uuid;not null;index" json:"supplier_id"`
	QuoteNumber string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_quote" json:"quote_number"`
	Status      string    `gorm:"type:varchar(50);default:'submitted';index" json:"status"` // submitted, accepted, rejected, expired
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	ValidUntil  time.Time `json:"valid_until"`
	Notes       string    `gorm:"type:text" json:"notes"`
	SubmittedAt time.Time `json:"submitted_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	RFQ    RFQ      `gorm:"foreignKey:RFQID" json:"rfq,omitempty"`
	Items  []QuoteItem `gorm:"foreignKey:QuoteID" json:"items,omitempty"`
}

func (Quote) TableName() string {
	return "procurement.quotes"
}

type QuoteItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	QuoteID     uuid.UUID `gorm:"type:uuid;not null;index" json:"quote_id"`
	PRItemID    uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_item_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	LeadTime    int       `json:"lead_time"` // days
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Quote Quote `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
}

func (QuoteItem) TableName() string {
	return "procurement.quote_items"
}

type PurchaseOrder struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PRID        uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_id"`
	RFQID       uuid.UUID `gorm:"type:uuid;index" json:"rfq_id"`
	QuoteID     uuid.UUID `gorm:"type:uuid;not null;index" json:"quote_id"`
	PONumber    string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_po" json:"po_number"`
	Status      string    `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, confirmed, shipped, delivered, cancelled
	TotalAmount    float64    `json:"total_amount"`
	Currency       string      `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	PaymentMode    string      `gorm:"type:varchar(50);default:'DIRECT';index" json:"payment_mode"` // DIRECT, ESCROW
	PaymentStatus  string      `gorm:"type:varchar(50);default:'pending';index" json:"payment_status"` // pending, processing, succeeded, failed
	PaymentID      *uuid.UUID  `gorm:"type:uuid;index" json:"payment_id,omitempty"`
	SupplierID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"supplier_id"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PR    PurchaseRequest `gorm:"foreignKey:PRID" json:"pr,omitempty"`
	RFQ   RFQ             `gorm:"foreignKey:RFQID" json:"rfq,omitempty"`
	Quote Quote           `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	Items []POItem       `gorm:"foreignKey:POID" json:"items,omitempty"`
}

func (PurchaseOrder) TableName() string {
	return "procurement.purchase_orders"
}

type POItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	POID        uuid.UUID `gorm:"type:uuid;not null;index" json:"po_id"`
	PRItemID    uuid.UUID `gorm:"type:uuid;not null;index" json:"pr_item_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PO PurchaseOrder `gorm:"foreignKey:POID" json:"po,omitempty"`
}

func (POItem) TableName() string {
	return "procurement.po_items"
}
