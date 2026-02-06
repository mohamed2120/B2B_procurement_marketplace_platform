package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatThread struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Title       string         `gorm:"type:varchar(255)" json:"title"`
	ThreadType  string         `gorm:"type:varchar(50);not null;index" json:"thread_type"` // order, rfq, quote, dispute, general
	ReferenceID *uuid.UUID     `gorm:"type:uuid;index" json:"reference_id,omitempty"` // ID of related order/RFQ/etc
	IsArchived  bool           `gorm:"default:false" json:"is_archived"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Participants []ThreadParticipant `gorm:"foreignKey:ThreadID" json:"participants,omitempty"`
	Messages     []ChatMessage       `gorm:"foreignKey:ThreadID" json:"messages,omitempty"`
}

func (ChatThread) TableName() string {
	return "collaboration.chat_threads"
}

type ThreadParticipant struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ThreadID  uuid.UUID `gorm:"type:uuid;not null;index" json:"thread_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Role      string    `gorm:"type:varchar(50)" json:"role"` // buyer, supplier, admin
	JoinedAt  time.Time `json:"joined_at"`

	// Relationships
	Thread ChatThread `gorm:"foreignKey:ThreadID" json:"thread,omitempty"`
}

func (ThreadParticipant) TableName() string {
	return "collaboration.thread_participants"
}

type ChatMessage struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ThreadID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"thread_id"`
	SenderID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"sender_id"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	MessageType string         `gorm:"type:varchar(50);default:'text'" json:"message_type"` // text, file, system
	IsRead      bool           `gorm:"default:false" json:"is_read"`
	IsEdited    bool           `gorm:"default:false" json:"is_edited"`
	IsDeleted   bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Thread ChatThread    `gorm:"foreignKey:ThreadID" json:"thread,omitempty"`
	Files  []MessageFile `gorm:"foreignKey:MessageID" json:"files,omitempty"`
}

func (ChatMessage) TableName() string {
	return "collaboration.chat_messages"
}

type MessageFile struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MessageID uuid.UUID `gorm:"type:uuid;not null;index" json:"message_id"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL   string    `gorm:"type:text;not null" json:"file_url"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `gorm:"type:varchar(100)" json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Message ChatMessage `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}

func (MessageFile) TableName() string {
	return "collaboration.message_files"
}

type Dispute struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	DisputeType string         `gorm:"type:varchar(50);not null" json:"dispute_type"` // quality, delivery, payment, other
	Status      string         `gorm:"type:varchar(50);default:'open';index" json:"status"` // open, in_review, resolved, closed
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	RaisedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"raised_by"`
	ResolvedBy  *uuid.UUID     `gorm:"type:uuid" json:"resolved_by,omitempty"`
	Resolution  string         `gorm:"type:text" json:"resolution,omitempty"`
	ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Messages []ChatMessage `gorm:"foreignKey:ThreadID" json:"messages,omitempty"`
}

func (Dispute) TableName() string {
	return "collaboration.disputes"
}

type Rating struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	RatedBy    uuid.UUID `gorm:"type:uuid;not null;index" json:"rated_by"`
	RatedEntityType string `gorm:"type:varchar(50);not null" json:"rated_entity_type"` // supplier, buyer, product
	RatedEntityID uuid.UUID `gorm:"type:uuid;not null;index" json:"rated_entity_id"`
	Rating      int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment     string    `gorm:"type:text" json:"comment"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	IsModerated bool      `gorm:"default:false" json:"is_moderated"`
	ModeratedBy *uuid.UUID `gorm:"type:uuid" json:"moderated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
}

func (Rating) TableName() string {
	return "collaboration.ratings"
}
