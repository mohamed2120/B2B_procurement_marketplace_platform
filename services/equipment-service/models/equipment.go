package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Equipment struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EquipmentNumber string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_eq_num" json:"equipment_number"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Type            string         `gorm:"type:varchar(100)" json:"type"` // excavator, loader, crane, etc.
	Manufacturer    string         `gorm:"type:varchar(255)" json:"manufacturer"`
	Model           string         `gorm:"type:varchar(255)" json:"model"`
	SerialNumber    string         `gorm:"type:varchar(255)" json:"serial_number"`
	Year            int            `json:"year"`
	Status          string         `gorm:"type:varchar(50);default:'active';index" json:"status"` // active, maintenance, retired
	Location        string         `gorm:"type:varchar(255)" json:"location"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	BOMNodes        []BOMNode            `gorm:"foreignKey:EquipmentID" json:"bom_nodes,omitempty"`
	CompatibilityMappings []CompatibilityMapping `gorm:"foreignKey:EquipmentID" json:"compatibility_mappings,omitempty"`
}

func (Equipment) TableName() string {
	return "equipment.equipment"
}

type BOMNode struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"equipment_id"`
	PartID      *uuid.UUID `gorm:"type:uuid;index" json:"part_id,omitempty"` // Reference to catalog library part
	PartNumber  string    `gorm:"type:varchar(255)" json:"part_number"` // If not in library
	PartName    string    `gorm:"type:varchar(255);not null" json:"part_name"`
	Description string    `gorm:"type:text" json:"description"`
	Quantity    float64   `gorm:"not null;default:1" json:"quantity"`
	Unit        string    `gorm:"type:varchar(50)" json:"unit"`
	Position    string    `gorm:"type:varchar(100)" json:"position"` // Location in equipment
	ParentNodeID *uuid.UUID `gorm:"type:uuid;index" json:"parent_node_id,omitempty"`
	Level       int       `gorm:"default:0" json:"level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Equipment  Equipment  `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	ParentNode *BOMNode   `gorm:"foreignKey:ParentNodeID" json:"parent_node,omitempty"`
	ChildNodes []BOMNode  `gorm:"foreignKey:ParentNodeID" json:"child_nodes,omitempty"`
}

func (BOMNode) TableName() string {
	return "equipment.bom_nodes"
}

type CompatibilityMapping struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"equipment_id"`
	PartID      uuid.UUID `gorm:"type:uuid;not null;index" json:"part_id"` // Catalog library part
	IsCompatible bool     `gorm:"default:true" json:"is_compatible"`
	Notes       string    `gorm:"type:text" json:"notes"`
	VerifiedBy  *uuid.UUID `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Equipment Equipment `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
}

func (CompatibilityMapping) TableName() string {
	return "equipment.compatibility_mappings"
}
