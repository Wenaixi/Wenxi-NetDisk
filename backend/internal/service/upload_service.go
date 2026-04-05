package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

const (
	// ChunkSize 每个分块的大小 (2MB)
	ChunkSize = 2 * 1024 * 1024

	// StatusPending 待上传
	StatusPending = "pending"
	// StatusUploading 上传中
	StatusUploading = "uploading"
	// StatusCompleted 已完成
	StatusCompleted = "completed"
	// StatusFailed 失败
	StatusFailed = "failed"
)

// UploadSessionRepository 上传会话仓库接口
type UploadSessionRepository interface {
	Create(session *model.UploadSession) error
	FindByID(id uint) (*model.UploadSession, error)
	FindByUserIDAndHash(userID uint, fileHash string) (*model.UploadSession, error)
	Update(session *model.UploadSession) error
	Delete(id uint) error
	DeleteByUserID(userID uint) error
}

// LanZouClientProvider 蓝奏云客户端提供者接口
type LanZouClientProvider interface {
	IsConnected(userID uint) bool
	GetClient(userID uint) (*lanzou.Client, error)
}

// FileMetadataCreator 文件元数据创建器接口
type FileMetadataCreator interface {
	CreateMetadata(userID uint, req *CreateFileRequest) (*model.File, error)
}

type UploadService struct {
	uploadRepo   UploadSessionRepository
	lanzouSvc    LanZouClientProvider
	fileSvc      FileMetadataCreator
}

func NewUploadService(
	uploadRepo UploadSessionRepository,
	lanzouSvc LanZouClientProvider,
	fileSvc FileMetadataCreator,
) *UploadService {
	return &UploadService{
		uploadRepo:   uploadRepo,
		lanzouSvc:    lanzouSvc,
		fileSvc:      fileSvc,
	}
}

// InitializeUploadRequest 初始化上传请求
type InitializeUploadRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
	MimeType string `json:"mime_type"`
	FolderID uint   `json:"folder_id"`
}

// InitializeUploadResponse 初始化上传响应
type InitializeUploadResponse struct {
	SessionID  uint   `json:"session_id"`
	UploadURL  string `json:"upload_url"`
	ChunkSize  int    `json:"chunk_size"`
	TotalChunks int   `json:"total_chunks"`
	FinalName   string `json:"final_name"` // 最终上传的文件名（可能被混淆）
}

// InitializeUpload 初始化上传会话
func (s *UploadService) InitializeUpload(userID uint, req *InitializeUploadRequest) (*InitializeUploadResponse, error) {
	// 检查用户是否已连接蓝奏云
	if !s.lanzouSvc.IsConnected(userID) {
		return nil, errors.New("lanzou cloud not connected")
	}

	// 计算文件哈希（使用文件名+大小+时间戳的临时哈希）
	hash := calculateFileHash(req.FileName, req.FileSize)

	// 检查是否存在未完成的上传会话
	existingSession, err := s.uploadRepo.FindByUserIDAndHash(userID, hash)
	if err == nil && existingSession != nil && existingSession.Status != StatusCompleted {
		// 返回已有的上传会话
		finalName := existingSession.FinalFileName
		if finalName == "" {
			finalName = existingSession.FileName
		}
		return &InitializeUploadResponse{
			SessionID:   existingSession.ID,
			ChunkSize:   ChunkSize,
			TotalChunks: existingSession.ChunksTotal,
			FinalName:   finalName,
		}, nil
	}

	// 计算分块数量
	totalChunks := int(math.Ceil(float64(req.FileSize) / float64(ChunkSize)))

	// 处理文件名：检查扩展名是否支持，不支持则添加混淆后缀
	finalFileName := req.FileName
	ext := lanzou.GetFileExtension(req.FileName)
	if ext != "" && !lanzou.IsSupportedExtension(ext) {
		finalFileName = lanzou.CreateSpecificName(req.FileName)
	}

	// 创建新的上传会话
	session := &model.UploadSession{
		UserID:         userID,
		FileName:       req.FileName,
		FinalFileName:  finalFileName,
		FileSize:       req.FileSize,
		FileHash:       hash,
		ChunksTotal:    totalChunks,
		ChunksUploaded: 0,
		Status:         StatusUploading,
	}

	if err := s.uploadRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create upload session: %w", err)
	}

	return &InitializeUploadResponse{
		SessionID:   session.ID,
		ChunkSize:   ChunkSize,
		TotalChunks: totalChunks,
		FinalName:   finalFileName,
	}, nil
}

// UploadChunkRequest 上传分块请求
type UploadChunkRequest struct {
	ChunkIndex int    `json:"chunk_index" binding:"required"`
	Data       []byte `json:"data" binding:"required"`
	FolderID   int    `json:"folder_id"` // 目标文件夹ID
}

// UploadChunkResponse 上传分块响应
type UploadChunkResponse struct {
	LanZouFileID  string `json:"lanzou_file_id"`
	DownloadURL   string `json:"download_url"`
	FileName      string `json:"file_name"`
}

// UploadChunk 处理分块上传
func (s *UploadService) UploadChunk(userID uint, sessionID uint, req *UploadChunkRequest) (*UploadChunkResponse, error) {
	session, err := s.uploadRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("upload session not found")
	}

	if session.UserID != userID {
		return nil, errors.New("access denied")
	}

	if session.Status == StatusCompleted {
		return nil, errors.New("upload already completed")
	}

	// 验证分块索引
	if req.ChunkIndex < 0 || req.ChunkIndex >= session.ChunksTotal {
		return nil, errors.New("invalid chunk index")
	}

	// 获取蓝奏云客户端
	client, err := s.lanzouSvc.GetClient(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lanzou client: %w", err)
	}

	// 使用最终文件名（可能已被混淆）
	fileName := session.FinalFileName
	if fileName == "" {
		fileName = session.FileName
	}

	// 目标文件夹ID
	folderID := req.FolderID
	if folderID == 0 {
		folderID = -1 // 默认根目录
	}

	// 调用Html5Upload上传
	reader := bytes.NewReader(req.Data)
	resp, err := client.Html5Upload(reader, fileName, int64(len(req.Data)), folderID)
	if err != nil {
		return nil, fmt.Errorf("upload to lanzou failed: %w", err)
	}

	if resp.Zt == 0 {
		return nil, fmt.Errorf("upload failed: %s", resp.Info)
	}

	// 提取上传结果
	if len(resp.Text) == 0 {
		return nil, errors.New("upload succeeded but no file info returned")
	}

	fileInfo := resp.Text[0]

	// 更新已上传分块数
	session.ChunksUploaded++
	if session.ChunksUploaded >= session.ChunksTotal {
		session.Status = StatusCompleted
		// 保存蓝奏云文件ID
		session.LanZouUploadURL = fileInfo.ID
	}

	if err := s.uploadRepo.Update(session); err != nil {
		return nil, err
	}

	return &UploadChunkResponse{
		LanZouFileID: fileInfo.ID,
		DownloadURL:  fileInfo.IsNew,
		FileName:     fileInfo.Name,
	}, nil
}

// CompleteUploadRequest 完成上传请求
type CompleteUploadRequest struct {
	EncryptionKey   string `json:"encryption_key" binding:"required"`
	EncryptionNonce string `json:"encryption_nonce" binding:"required"`
	LanZouFileID    string `json:"lanzou_file_id" binding:"required"`
}

// CompleteUpload 完成上传并保存文件元数据
func (s *UploadService) CompleteUpload(userID uint, sessionID uint, req *CompleteUploadRequest) (*model.File, error) {
	session, err := s.uploadRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("upload session not found")
	}

	if session.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 创建文件元数据
	folderID := ""
	fileReq := &CreateFileRequest{
		Name:            session.FileName,
		Size:            session.FileSize,
		LanZouFileID:    req.LanZouFileID,
		LanZouFolderID:  folderID,
		EncryptionKey:   req.EncryptionKey,
		EncryptionNonce: req.EncryptionNonce,
		MimeType:        "application/octet-stream",
	}

	file, err := s.fileSvc.CreateMetadata(userID, fileReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create file metadata: %w", err)
	}

	// 更新会话状态为已完成
	session.Status = StatusCompleted
	if err := s.uploadRepo.Update(session); err != nil {
		return nil, err
	}

	return file, nil
}

// GetUploadStatus 获取上传状态
func (s *UploadService) GetUploadStatus(userID uint, sessionID uint) (*model.UploadSession, error) {
	session, err := s.uploadRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("upload session not found")
	}

	if session.UserID != userID {
		return nil, errors.New("access denied")
	}

	return session, nil
}

// ResumeUpload 恢复上传
func (s *UploadService) ResumeUpload(userID uint, fileHash string) (*InitializeUploadResponse, error) {
	session, err := s.uploadRepo.FindByUserIDAndHash(userID, fileHash)
	if err != nil {
		return nil, errors.New("no upload session found for resume")
	}

	if session.Status == StatusCompleted {
		return nil, errors.New("upload already completed")
	}

	return &InitializeUploadResponse{
		SessionID:   session.ID,
		UploadURL:   session.LanZouUploadURL,
		ChunkSize:   ChunkSize,
		TotalChunks: session.ChunksTotal,
	}, nil
}

// calculateFileHash 计算文件哈希
func calculateFileHash(fileName string, fileSize int64) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%d", fileName, fileSize)))
	return hex.EncodeToString(h.Sum(nil))[:32]
}

// GetLanZouClient 获取蓝奏云客户端（用于直接上传）
func (s *UploadService) GetLanZouClient(userID uint) (*lanzou.Client, error) {
	return s.lanzouSvc.GetClient(userID)
}
