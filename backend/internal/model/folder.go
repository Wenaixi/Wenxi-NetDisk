package model

import (
	"time"
)

type Folder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"size:500" json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Parent    *Folder   `gorm:"foreignKey:ParentID" json:"-"`
}
