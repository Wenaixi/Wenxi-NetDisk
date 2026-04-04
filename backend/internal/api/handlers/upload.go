package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	fileSvc     *service.FileService
	versionSvc  *service.FileVersionService
}

func NewUploadHandler(fileSvc *service.FileService, versionSvc *service.FileVersionService) *UploadHandler {
	return &UploadHandler{
		fileSvc:    fileSvc,
		versionSvc: versionSvc,
	}
}

// UploadFile 处理 multipart 文件上传（本地存储/简化上传）
func (h *UploadHandler) UploadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	// 解析表单
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	// 获取文件夹ID（可选）
	folderIDStr := c.PostForm("folder_id")
	var folderID *uint
	if folderIDStr != "" {
		id, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err == nil {
			pid := uint(id)
			folderID = &pid
		}
	}

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		response.InternalError(c, "failed to read file")
		return
	}

	// 计算文件哈希
	hash := sha256.Sum256(data)
	fileHash := fmt.Sprintf("%x", hash)[:32]

	// 文件大小
	fileSize := int64(len(data))

	// 生成lanzou_file_id（临时使用文件哈希）
	lanzouFileID := fileHash

	// MIME类型
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		ext := filepath.Ext(header.Filename)
		mimeType = guessMimeType(ext)
	}

	// 创建文件元数据
	fileReq := &service.CreateFileRequest{
		Name:            header.Filename,
		Size:            fileSize,
		LanZouFileID:    lanzouFileID,
		LanZouFolderID:  "",
		EncryptionKey:   fileHash, // 使用文件哈希作为临时密钥
		EncryptionNonce: fmt.Sprintf("%x", hash[:12]),
		MimeType:        mimeType,
	}

	createdFile, err := h.fileSvc.CreateMetadata(uid, fileReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 如果有folder_id，移动文件到指定文件夹
	if folderID != nil {
		_, _ = h.fileSvc.MoveFile(uid, createdFile.ID, folderID)
	}

	response.Success(c, gin.H{
		"id":          createdFile.ID,
		"name":        createdFile.Name,
		"size":        createdFile.Size,
		"mime_type":   createdFile.MimeType,
		"created_at":  createdFile.CreatedAt,
	})
}

// GetUploadURL 获取上传URL（用于直传蓝奏云）
func (h *UploadHandler) GetUploadURL(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req struct {
		FileName string `json:"file_name" binding:"required"`
		Size     int64  `json:"size" binding:"required"`
		MimeType string `json:"mime_type"`
		FolderID *uint  `json:"folder_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 生成上传URL（简化实现，实际应调用蓝奏云API获取直传URL）
	uploadURL := fmt.Sprintf("https://pc.woozooo.com/fileup.php?task=1&user=%d&file=%s&size=%d",
		uid, req.FileName, req.Size)

	// 创建文件元数据（标记为待上传状态）
	tempHash := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%d", req.FileName, req.Size, time.Now().UnixNano())))
	fileHash := fmt.Sprintf("%x", tempHash)[:32]

	fileReq := &service.CreateFileRequest{
		Name:            req.FileName,
		Size:            req.Size,
		LanZouFileID:    "pending:" + fileHash,
		LanZouFolderID:  "",
		EncryptionKey:   fileHash,
		EncryptionNonce: fmt.Sprintf("%x", tempHash[:12]),
		MimeType:        req.MimeType,
	}

	createdFile, err := h.fileSvc.CreateMetadata(uid, fileReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 如果有folder_id，移动文件
	if req.FolderID != nil {
		_, _ = h.fileSvc.MoveFile(uid, createdFile.ID, req.FolderID)
	}

	response.Success(c, gin.H{
		"upload_url": uploadURL,
		"file_id":    createdFile.ID,
	})
}

// ListVersions 获取文件版本列表
func (h *UploadHandler) ListVersions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	versions, err := h.versionSvc.ListVersions(uid, uint(fileID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, versions)
}

// RestoreVersion 恢复文件版本
func (h *UploadHandler) RestoreVersion(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	versionID, err := strconv.ParseUint(c.Param("version_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	file, err := h.versionSvc.RestoreVersion(uid, uint(fileID), uint(versionID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":         file.ID,
		"name":       file.Name,
		"size":       file.Size,
		"updated_at": file.UpdatedAt,
	})
}

// guessMimeType 根据文件扩展名猜测MIME类型
func guessMimeType(ext string) string {
	mimeTypes := map[string]string{
		".txt":  "text/plain",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
	}
	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}
