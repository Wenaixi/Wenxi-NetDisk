package lanzou

import "time"

// FileType 文件类型
type FileType int

const (
	FileTypeFile   FileType = 1
	FileTypeFolder FileType = 2
)

// Token 蓝奏云认证Token
type Token struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Cookie    string    `json:"cookie"` // PC版Cookie
	TokenVal  string    `json:"token_value"` // Token值
	UID       string    `json:"uid"` // 用户ID
	ExpiresAt time.Time `json:"expires_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileInfo 文件信息
type FileInfo struct {
	ID         string `json:"id"`          // 文件ID
	Name       string `json:"name"`        // 文件名
	NameAll    string `json:"name_all"`    // 完整文件名
	Size       int64  `json:"size"`        // 文件大小
	Time       string `json:"time"`        // 上传时间
	Icon       string `json:"icon"`         // 图标
	Downs      int    `json:"downs"`       // 下载次数
	Type       FileType `json:"type"`       // 1文件 2文件夹
}

// FolderInfo 文件夹信息
type FolderInfo struct {
	ID        uint   `json:"id"`        // 文件夹ID
	UserID    uint   `json:"user_id"`    // 所属用户
	ParentID  uint   `json:"parent_id"`  // 父文件夹ID
	Name      string `json:"name"`       // 文件夹名
	Time      string `json:"time"`      // 创建时间
	FileCount int    `json:"file_count"` // 文件数量
}

// Task5Response task=5 列出文件响应
type Task5Response struct {
	Info interface{} `json:"info"`
	Text []FileInfo  `json:"text"`
	Zt   int         `json:"zt"` // 状态码
}

// Task47Response task=47 列出文件夹响应
type Task47Response struct {
	Info interface{} `json:"info"`
	Text []FolderInfo `json:"text"`
	Zt   int         `json:"zt"`
}

// Task2Response task=2 创建文件夹响应
type Task2Response struct {
	Info interface{} `json:"info"`
	Text struct {
		ID   uint   `json:"folder_id"`
		Name string `json:"folder_name"`
	} `json:"text"`
	Zt int `json:"zt"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	Info interface{} `json:"info"`
	Text struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Time     string `json:"time"`
		Size     int64  `json:"size"`
		Desc     string `json:"desc"`
		DownURL  string `json:"downurl"`
		Filename string `json:"filename"`
	} `json:"text"`
	Zt int `json:"zt"`
}

// DeleteResponse 删除响应
type DeleteResponse struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}

// RenameResponse 重命名响应
type RenameResponse struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}

// MoveResponse 移动响应
type MoveResponse struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}

// ShareResponse 分享响应
type ShareResponse struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
	Text struct {
		URL   string `json:"url"`   // 分享链接
		Pwd   string `json:"pwd"`   // 密码
		Descr string `json:"descr"` // 描述
	} `json:"text"`
}

// Task23Response task=23 设置文件访问密码响应
type Task23Response struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}

// Task16Response task=16 设置文件夹访问密码响应
type Task16Response struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}

// Task4Response task=4 重命名文件夹响应
type Task4Response struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
}
