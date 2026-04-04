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
func (s *DownloadService) getDownloadLink(client *lanzou.Client, file *model.File) (string, error) {
	// 蓝奏云下载流程：
	// 1. 如果有分享链接，直接使用分享链接获取下载直链
	// 2. 如果只有文件ID，需要先获取文件信息

	// 对于已上传到蓝奏云的文件，可以通过文件ID获取下载直链
	// 但蓝奏云API并不直接支持通过文件ID获取下载，需要通过分享页面

	// 简化实现：返回空字符串，客户端需要先创建分享再下载
	// 实际生产环境可能需要不同的处理方式
	if file.LanZouFileID != "" {
		// 尝试通过文件ID获取下载信息
		// 注意：蓝奏云实际是通过分享页面下载的
		return "", nil
	}

	return "", errors.New("download link not available")
}

// GetShareDownloadURL 通过分享链接获取下载直链
func (s *DownloadService) GetShareDownloadURL(shareToken, password string) (string, error) {
	// 通过分享链接获取下载直链
	// 这个方法用于公开分享下载（不验证用户身份）
	// 实际蓝奏云API：https://pc.woozooo.com/d.php?...

	// 简化实现
	return "", nil
}
