package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(50);default:'active';index" json:"status"` // active, suspended, closed
	IsVerified  bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Listings []Listing `gorm:"foreignKey:StoreID" json:"listings,omitempty"`
}

func (Store) TableName() string {
	return "marketplace.stores"
}

type Listing struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StoreID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	ListingType string         `gorm:"type:varchar(50);not null;index" json:"listing_type"` // product, service, surplus
	PartID      *uuid.UUID     `gorm:"type:uuid;index" json:"part_id,omitempty"` // Reference to catalog part
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	SKU         string         `gorm:"type:varchar(100)" json:"sku"`
	Status      string         `gorm:"type:varchar(50);default:'draft';index" json:"status"` // draft, active, sold_out, inactive
	Price       float64        `gorm:"not null" json:"price"`
	Currency    string         `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	StockQuantity float64      `gorm:"default:0" json:"stock_quantity"`
	MinOrderQuantity float64   `gorm:"default:1" json:"min_order_quantity"`
	LeadTimeDays int          `gorm:"default:0" json:"lead_time_days"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Store      Store         `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Media      []ListingMedia `gorm:"foreignKey:ListingID" json:"media,omitempty"`
}

func (Listing) TableName() string {
	return "marketplace.listings"
}

type ListingMedia struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ListingID   uuid.UUID `gorm:"type:uuid;not null;index" json:"listing_id"`
	MediaType   string    `gorm:"type:varchar(50);not null" json:"media_type"` // image, video, document
	URL         string    `gorm:"type:text;not null" json:"url"`
	ThumbnailURL string   `gorm:"type:text" json:"thumbnail_url,omitempty"`
	FileName    string    `gorm:"type:varchar(255)" json:"file_name"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `gorm:"type:varchar(100)" json:"mime_type"`
	IsPrimary  bool       `gorm:"default:false" json:"is_primary"`
	SortOrder  int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Listing Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

func (ListingMedia) TableName() string {
	return "marketplace.listing_media"
}
