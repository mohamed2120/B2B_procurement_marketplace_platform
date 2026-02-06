package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	LegalName       string         `gorm:"type:varchar(255)" json:"legal_name"`
	TaxID           string         `gorm:"type:varchar(100)" json:"tax_id"`
	Subdomain       string         `gorm:"type:varchar(100);unique" json:"subdomain"`
	Status          string         `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, approved, rejected, suspended
	VerificationStatus string       `gorm:"type:varchar(50);default:'pending'" json:"verification_status"`
	Address         string         `gorm:"type:text" json:"address"`
	City            string         `gorm:"type:varchar(100)" json:"city"`
	State           string         `gorm:"type:varchar(100)" json:"state"`
	Country         string         `gorm:"type:varchar(100)" json:"country"`
	PostalCode      string         `gorm:"type:varchar(20)" json:"postal_code"`
	Phone           string         `gorm:"type:varchar(50)" json:"phone"`
	Email           string         `gorm:"type:varchar(255)" json:"email"`
	Website         string         `gorm:"type:varchar(255)" json:"website"`
	Industry        string         `gorm:"type:varchar(100)" json:"industry"`
	CompanyType     string         `gorm:"type:varchar(50)" json:"company_type"` // buyer, supplier, both
	ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy      *uuid.UUID     `json:"approved_by,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Documents []CompanyDocument `gorm:"foreignKey:CompanyID" json:"documents,omitempty"`
}

func (Company) TableName() string {
	return "company.companies"
}

type CompanyDocument struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`
	DocumentType string   `gorm:"type:varchar(100);not null" json:"document_type"` // registration, tax_cert, license, etc.
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL     string    `gorm:"type:text;not null" json:"file_url"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `gorm:"type:varchar(100)" json:"mime_type"`
	UploadedBy  uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	Status      string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, verified, rejected
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	VerifiedBy  *uuid.UUID `json:"verified_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (CompanyDocument) TableName() string {
	return "company.documents"
}

type SubdomainRequest struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`
	Subdomain   string    `gorm:"type:varchar(100);not null;unique" json:"subdomain"`
	Status      string    `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, approved, rejected
	RequestedBy uuid.UUID `gorm:"type:uuid;not null" json:"requested_by"`
	ReviewedBy  *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	Reason      string    `gorm:"type:text" json:"reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (SubdomainRequest) TableName() string {
	return "company.subdomain_requests"
}
