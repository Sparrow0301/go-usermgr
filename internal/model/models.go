package model

import "time"

const (
	UserStatusEnabled  = "enabled"
	UserStatusDisabled = "disabled"
)

type User struct {
	ID           uint       `gorm:"primaryKey"`
	Username     string     `gorm:"size=50;uniqueIndex;not null"`
	Email        string     `gorm:"size=255;uniqueIndex;not null"`
	PasswordHash string     `gorm:"size=255;not null"`
	FullName     string     `gorm:"size=100"`
	Status       string     `gorm:"size=20;default:'enabled'"`
	LastLoginAt  *time.Time `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size=50;uniqueIndex;not null"`
	Description string `gorm:"size=255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"size=100;uniqueIndex;not null"`
	Description string `gorm:"size=255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRole struct {
	UserID    uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
	CreatedAt    time.Time
}
