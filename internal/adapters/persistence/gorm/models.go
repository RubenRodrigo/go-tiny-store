package gorm

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for GORM models
type Base struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// UserModel represents the GORM model for users
type UserModel struct {
	Base
	Email               string                    `gorm:"unique;not null"`
	Username            string                    `gorm:"not null"`
	Password            string                    `gorm:"not null"`
	FirstName           string                    `gorm:""`
	LastName            string                    `gorm:""`
	RefreshTokens       []RefreshTokenModel       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PasswordResetTokens []PasswordResetTokenModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Roles               []RoleModel               `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName overrides the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// RoleModel represents the GORM model for roles
type RoleModel struct {
	Base
	Name        string `gorm:"unique;not null"`
	Description string `gorm:""`
}

// TableName overrides the table name for RoleModel
func (RoleModel) TableName() string {
	return "roles"
}

// RefreshTokenModel represents the GORM model for refresh tokens
type RefreshTokenModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:""`
	UpdatedAt time.Time `gorm:""`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:""`
	UserID    string    `gorm:"type:uuid;not null;index"`
}

// TableName overrides the table name for RefreshTokenModel
func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

// PasswordResetTokenModel represents the GORM model for password reset tokens
type PasswordResetTokenModel struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time  `gorm:""`
	UpdatedAt time.Time  `gorm:""`
	TokenHash string     `gorm:"unique;not null"`
	ExpiresAt time.Time  `gorm:""`
	UsedAt    *time.Time `gorm:""`
	UserID    string     `gorm:"type:uuid;not null;index"`
}

// TableName overrides the table name for PasswordResetTokenModel
func (PasswordResetTokenModel) TableName() string {
	return "password_reset_tokens"
}

// ProductModel represents the GORM model for products
type ProductModel struct {
	Base
	Name       string              `gorm:"not null;size:255"`
	Price      float64             `gorm:"not null"`
	Disabled   bool                `gorm:"default:false"`
	Stock      int                 `gorm:"default:0"`
	CategoryID string              `gorm:"type:uuid"`
	Images     []ProductImageModel `gorm:"foreignKey:ProductID"`
}

// TableName overrides the table name for ProductModel
func (ProductModel) TableName() string {
	return "products"
}

// ProductImageModel represents the GORM model for product images
type ProductImageModel struct {
	Base
	ProductID string `gorm:"type:uuid;not null;index"`
	URL       string `gorm:"not null;size:500"`
	AltText   string `gorm:"size:255"`
}

// TableName overrides the table name for ProductImageModel
func (ProductImageModel) TableName() string {
	return "product_images"
}

// CategoryModel represents the GORM model for categories
type CategoryModel struct {
	Base
	Name     string         `gorm:"not null;size:100"`
	Products []ProductModel `gorm:"foreignKey:CategoryID"`
}

// TableName overrides the table name for CategoryModel
func (CategoryModel) TableName() string {
	return "categories"
}

// AllModels returns all GORM models for schema migration tools (Atlas, etc.)
func AllModels() []interface{} {
	return []interface{}{
		&UserModel{},
		&RoleModel{},
		&RefreshTokenModel{},
		&PasswordResetTokenModel{},
		&ProductModel{},
		&ProductImageModel{},
		&CategoryModel{},
	}
}
