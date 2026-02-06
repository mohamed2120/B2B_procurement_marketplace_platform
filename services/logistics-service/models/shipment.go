package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shipment struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	POID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"po_id"`
	TrackingNumber string        `gorm:"type:varchar(100);unique" json:"tracking_number"`
	Status        string         `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, in_transit, delivered, delayed
	Carrier       string         `gorm:"type:varchar(100)" json:"carrier"`
	ETA           time.Time      `gorm:"not null;index" json:"eta"`
	ActualDeliveryDate *time.Time `json:"actual_delivery_date,omitempty"`
	Origin        string         `gorm:"type:varchar(255)" json:"origin"`
	Destination   string         `gorm:"type:varchar(255)" json:"destination"`
	IsLate        bool           `gorm:"default:false;index" json:"is_late"`
	LateAlertSent bool           `gorm:"default:false" json:"late_alert_sent"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	TrackingEvents []TrackingEvent `gorm:"foreignKey:ShipmentID" json:"tracking_events,omitempty"`
	POD            *ProofOfDelivery `gorm:"foreignKey:ShipmentID" json:"pod,omitempty"`
}

func (Shipment) TableName() string {
	return "logistics.shipments"
}

type TrackingEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ShipmentID  uuid.UUID `gorm:"type:uuid;not null;index" json:"shipment_id"`
	EventType   string    `gorm:"type:varchar(50);not null" json:"event_type"` // picked_up, in_transit, out_for_delivery, delivered
	Location    string    `gorm:"type:varchar(255)" json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	Timestamp   time.Time `gorm:"not null" json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Shipment Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
}

func (TrackingEvent) TableName() string {
	return "logistics.tracking_events"
}

type ProofOfDelivery struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ShipmentID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"shipment_id"`
	SignedBy    string    `gorm:"type:varchar(255)" json:"signed_by"`
	SignatureURL string   `gorm:"type:text" json:"signature_url"`
	DeliveredAt time.Time `gorm:"not null" json:"delivered_at"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Shipment Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
}

func (ProofOfDelivery) TableName() string {
	return "logistics.proof_of_delivery"
}
