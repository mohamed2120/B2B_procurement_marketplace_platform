package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SharedInventory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PartID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"part_id"` // Catalog part
	EquipmentID *uuid.UUID     `gorm:"type:uuid;index" json:"equipment_id,omitempty"`
	Quantity    float64        `gorm:"not null" json:"quantity"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`
	IsAvailable bool           `gorm:"default:true;index" json:"is_available"`
	ReservedQty float64        `gorm:"default:0" json:"reserved_qty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
}

func (SharedInventory) TableName() string {
	return "virtual_warehouse.shared_inventory"
}

type EquipmentGroup struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Members []EquipmentGroupMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

func (EquipmentGroup) TableName() string {
	return "virtual_warehouse.equipment_groups"
}

type EquipmentGroupMember struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupID     uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"equipment_id"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Group EquipmentGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

func (EquipmentGroupMember) TableName() string {
	return "virtual_warehouse.equipment_group_members"
}

type InterCompanyTransfer struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromTenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"from_tenant_id"`
	ToTenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"to_tenant_id"`
	PartID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"part_id"`
	Quantity        float64        `gorm:"not null" json:"quantity"`
	Status          string         `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, approved, rejected, completed
	RequestedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy      *uuid.UUID     `gorm:"type:uuid" json:"approved_by,omitempty"`
	RejectionReason string         `gorm:"type:text" json:"rejection_reason,omitempty"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
}

func (InterCompanyTransfer) TableName() string {
	return "virtual_warehouse.inter_company_transfers"
}

type EmergencySourcing struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	PartID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"part_id"`
	EquipmentID *uuid.UUID     `gorm:"type:uuid;index" json:"equipment_id,omitempty"`
	Quantity    float64        `gorm:"not null" json:"quantity"`
	Priority    string         `gorm:"type:varchar(50);default:'high'" json:"priority"` // low, normal, high, urgent
	Status      string         `gorm:"type:varchar(50);default:'open';index" json:"status"` // open, sourcing, fulfilled, cancelled
	RequestedBy uuid.UUID      `gorm:"type:uuid;not null" json:"requested_by"`
	FulfilledBy *uuid.UUID     `gorm:"type:uuid" json:"fulfilled_by,omitempty"`
	FulfilledAt *time.Time     `json:"fulfilled_at,omitempty"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
}

func (EmergencySourcing) TableName() string {
	return "virtual_warehouse.emergency_sourcing"
}
