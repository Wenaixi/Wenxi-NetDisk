package model

import (
	"time"
)

type Share struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	FileID       uint      `gorm:"not null;index" json:"file_id"`
	ShareToken   string    `gorm:"uniqueIndex;size:32;not null" json:"share_token"`
	PasswordHash *string   `gorm:"size:100" json:"-"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	File         File      `gorm:"foreignKey:FileID" json:"-"`
}
