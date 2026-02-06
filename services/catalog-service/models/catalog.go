package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manufacturer struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Code        string    `gorm:"type:varchar(100);unique" json:"code"`
	Website     string    `gorm:"type:varchar(255)" json:"website"`
	Country     string    `gorm:"type:varchar(100)" json:"country"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Parts []LibraryPart `gorm:"foreignKey:ManufacturerID" json:"parts,omitempty"`
}

func (Manufacturer) TableName() string {
	return "catalog.manufacturers"
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Code        string    `gorm:"type:varchar(100);unique" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Parent   *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Parts    []LibraryPart  `gorm:"foreignKey:CategoryID" json:"parts,omitempty"`
}

func (Category) TableName() string {
	return "catalog.categories"
}

type Attribute struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Code        string    `gorm:"type:varchar(100);unique" json:"code"`
	DataType    string    `gorm:"type:varchar(50);not null" json:"data_type"` // string, number, boolean, date
	Unit        string    `gorm:"type:varchar(50)" json:"unit"`
	IsRequired  bool      `gorm:"default:false" json:"is_required"`
	IsSearchable bool     `gorm:"default:true" json:"is_searchable"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PartAttributes []PartAttribute `gorm:"foreignKey:AttributeID" json:"part_attributes,omitempty"`
}

func (Attribute) TableName() string {
	return "catalog.attributes"
}

type LibraryPart struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PartNumber      string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_part_number_mfr" json:"part_number"`
	ManufacturerID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"manufacturer_id"`
	CategoryID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Status          string         `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, approved, rejected
	ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy      *uuid.UUID     `gorm:"type:uuid" json:"approved_by,omitempty"`
	RejectedReason  string         `gorm:"type:text" json:"rejected_reason,omitempty"`
	IsDuplicate     bool           `gorm:"default:false;index" json:"is_duplicate"`
	DuplicateOf     *uuid.UUID     `gorm:"type:uuid;index" json:"duplicate_of,omitempty"`
	CreatedBy       uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Manufacturer    Manufacturer    `gorm:"foreignKey:ManufacturerID" json:"manufacturer,omitempty"`
	Category        Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	PartAttributes  []PartAttribute `gorm:"foreignKey:PartID" json:"part_attributes,omitempty"`
}

func (LibraryPart) TableName() string {
	return "catalog.library_parts"
}

type PartAttribute struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PartID      uuid.UUID `gorm:"type:uuid;not null;index" json:"part_id"`
	AttributeID uuid.UUID `gorm:"type:uuid;not null;index" json:"attribute_id"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Part      LibraryPart `gorm:"foreignKey:PartID" json:"part,omitempty"`
	Attribute Attribute   `gorm:"foreignKey:AttributeID" json:"attribute,omitempty"`
}

func (PartAttribute) TableName() string {
	return "catalog.part_attributes"
}
