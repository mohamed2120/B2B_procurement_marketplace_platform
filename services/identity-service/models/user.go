package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_email" json:"email"`
	PasswordHash string      `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	IsVerified bool          `gorm:"default:false" json:"is_verified"`
	LastLoginAt *time.Time   `json:"last_login_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	UserRoles []UserRole `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

func (User) TableName() string {
	return "identity.users"
}

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsSystem    bool      `gorm:"default:false" json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
}

func (Role) TableName() string {
	return "identity.roles"
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Resource    string    `gorm:"type:varchar(100);not null;index" json:"resource"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
}

func (Permission) TableName() string {
	return "identity.permissions"
}

type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Relationships
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "identity.role_permissions"
}

type UserRole struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "identity.user_roles"
}

type UserInvitation struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Email     string     `gorm:"type:varchar(255);not null" json:"email"`
	InvitedBy uuid.UUID  `gorm:"type:uuid;not null" json:"invited_by"`
	Token     string     `gorm:"type:varchar(255);not null;unique" json:"-"`
	RoleIDs   []uuid.UUID `gorm:"type:uuid[]" json:"role_ids"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	IsUsed    bool       `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (UserInvitation) TableName() string {
	return "identity.user_invitations"
}
