package model

import "time"

// FileVersion 文件版本表
type FileVersion struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	FileID          uint      `gorm:"not null;index" json:"file_id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	LanZouFileID    string    `gorm:"size:100" json:"lanzou_file_id"`
	Size            int64     `gorm:"not null" json:"size"`
	EncryptionKey   string    `gorm:"size:64;not null" json:"-"`
	EncryptionNonce string    `gorm:"size:24;not null" json:"-"`
	Description     string    `gorm:"size:500" json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}
