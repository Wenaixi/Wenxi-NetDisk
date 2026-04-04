package model

import (
	"time"
)

type LanZouToken struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Cookie     string    `gorm:"type:text;not null" json:"-"`
	TokenValue string    `gorm:"type:text" json:"-"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}
