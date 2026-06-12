package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement;type:int unsigned" json:"id"`
	Name         string         `gorm:"type:varchar(150);not null;index:idx_users_name" json:"name"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex:uq_users_email" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"type:enum('client','admin');not null;default:'client'" json:"role"`
	CreatedAt    time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime;index:idx_users_deleted_at" json:"deleted_at,omitempty"`
}
