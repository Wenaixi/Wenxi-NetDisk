package model

import "time"

// RecycleBin 回收站模型
// 删除的文件/文件夹先进入回收站，支持恢复
type RecycleBin struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	OriginalName  string    `gorm:"size:255;not null" json:"original_name"`
	OriginalPath  string    `gorm:"size:512" json:"original_path"` // 原始路径
	ItemType      string    `gorm:"size:20;not null" json:"item_type"` // 'file' or 'folder'
	ItemID        uint      `gorm:"not null;index" json:"item_id"` // 原始文件/文件夹ID
	LanZouFileID  string    `gorm:"size:100" json:"lanzou_file_id"`
	LanZouFolderID string   `gorm:"size:100" json:"lanzou_folder_id"`
	ParentID      uint      `gorm:"default:0" json:"parent_id"` // 原始父文件夹ID
	Size          int64     `gorm:"default:0" json:"size"`
	EncryptionKey string    `gorm:"size:64" json:"-"`
	EncryptionNonce string  `gorm:"size:24" json:"-"`
	DeletedAt     time.Time `json:"deleted_at"`
	ExpiresAt     time.Time `json:"expires_at"` // 回收站过期时间(30天)
	CreatedAt     time.Time `json:"created_at"`
}
