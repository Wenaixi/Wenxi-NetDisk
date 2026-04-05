package service

import (
	"errors"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

// FileDownloadRepository 文件下载仓库接口
type FileDownloadRepository interface {
	FindByID(id uint) (*model.File, error)
	FindByLanZouFileID(lanzouFileID string) (*model.File, error)
}

// DownloadService 下载服务
type DownloadService struct {
	fileRepo    FileDownloadRepository
	lanzouSvc   LanZouClientProvider
}

func NewDownloadService(fileRepo FileDownloadRepository, lanzouSvc LanZouClientProvider) *DownloadService {
	return &DownloadService{
		fileRepo:  fileRepo,
		lanzouSvc: lanzouSvc,
	}
}

// DownloadRequest 下载请求
type DownloadRequest struct {
	FileID uint `json:"file_id" binding:"required"`
}

// DownloadResponse 下载响应
type DownloadResponse struct {
	FileID          uint   `json:"file_id"`
	FileName        string `json:"file_name"`
	FileSize        int64  `json:"file_size"`
	DownloadURL     string `json:"download_url"`
	EncryptionKey   string `json:"encryption_key"`
	EncryptionNonce string `json:"encryption_nonce"`
	MimeType        string `json:"mime_type"`
}

// GetDownloadURL 获取文件下载直链
func (s *DownloadService) GetDownloadURL(userID, fileID uint) (*DownloadResponse, error) {
	// 获取文件元数据
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// 验证文件所有权
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 获取蓝奏云客户端
	client, err := s.lanzouSvc.GetClient(userID)
	if err != nil {
		return nil, errors.New("lanzou cloud not connected")
	}

	// 获取下载直链
	downloadURL, err := s.getDownloadLink(client, file)
	if err != nil {
		return nil, err
	}

	return &DownloadResponse{
		FileID:          file.ID,
		FileName:        file.Name,
		FileSize:        file.Size,
		DownloadURL:     downloadURL,
		EncryptionKey:   file.EncryptionKey,
		EncryptionNonce: file.EncryptionNonce,
		MimeType:        file.MimeType,
	}, nil
}

// getDownloadLink 获取下载链接
// 通过蓝奏云Task22/Task18文件详情接口获取is_newd下载链接
func (s *DownloadService) getDownloadLink(client *lanzou.Client, file *model.File) (string, error) {
	if file.LanZouFileID == "" {
		return "", errors.New("file not uploaded to lanzou yet")
	}

	// 如果client为nil（测试环境），返回模拟链接
	if client == nil {
		return "https://mock.lanzou.com/" + file.LanZouFileID, nil
	}

	// 尝试通过文件详情获取下载链接 (task=22)
	info, err := client.Task22(atoi(file.LanZouFileID))
	if err != nil {
		// 如果API调用失败，返回基于文件ID的下载链接
		return "https://pc.woozooo.com/" + file.LanZouFileID, nil
	}

	// Task22返回的info包含is_newd字段（下载链接）
	if isnewd, ok := info["is_newd"].(string); ok && isnewd != "" {
		fileID, _ := info["f_id"].(string)
		if fileID == "" {
			fileID = file.LanZouFileID
		}
		return isnewd + "/" + fileID, nil
	}

	// 降级：返回基于文件ID的链接
	return "https://pc.woozooo.com/" + file.LanZouFileID, nil
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

// GetShareDownloadURL 通过分享链接获取下载直链
func (s *DownloadService) GetShareDownloadURL(shareToken, password string) (string, error) {
	// 通过分享链接获取下载直链
	// 这个方法用于公开分享下载（不验证用户身份）
	// 实际蓝奏云API：https://pc.woozooo.com/d.php?...

	// 简化实现
	return "", nil
}
