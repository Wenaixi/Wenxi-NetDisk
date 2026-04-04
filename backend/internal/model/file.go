package model

import (
	"time"
)

type File struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Size            int64     `gorm:"not null" json:"size"`
	LanZouFileID    string    `gorm:"size:100" json:"lanzou_file_id"`
	LanZouFolderID  string    `gorm:"size:100" json:"lanzou_folder_id"`
	EncryptionKey   string    `gorm:"size:64;not null" json:"-"`
	EncryptionNonce string    `gorm:"size:24;not null" json:"-"`
	MimeType        string    `gorm:"size:100" json:"mime_type"`
	Description     string    `gorm:"size:500" json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	User            User      `gorm:"foreignKey:UserID" json:"-"`
}
