package model

import (
	"time"
)

type UploadSession struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	FileName        string    `gorm:"size:255;not null" json:"file_name"`
	FinalFileName   string    `gorm:"size:255" json:"final_file_name"`
	ObfuscatedName  string    `gorm:"size:255" json:"obfuscated_name"`
	FileSize        int64     `gorm:"not null" json:"file_size"`
	FileHash        string    `gorm:"size:64;not null" json:"file_hash"`
	ChunksTotal     int       `gorm:"not null" json:"chunks_total"`
	ChunksUploaded  int       `gorm:"default:0" json:"chunks_uploaded"`
	LanZouUploadURL string    `gorm:"size:500" json:"lanzou_upload_url"`
	LanZouFileID    string    `gorm:"size:50" json:"lanzou_file_id"`
	Status          string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
