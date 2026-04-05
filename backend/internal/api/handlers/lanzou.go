package handlers

import (
	"strconv"

	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type LanZouHandler struct {
	lanzouSvc  *service.LanZouService
	uploadSvc  *service.UploadService
}

func NewLanZouHandler(lanzouSvc *service.LanZouService, uploadSvc *service.UploadService) *LanZouHandler {
	return &LanZouHandler{
		lanzouSvc: lanzouSvc,
		uploadSvc: uploadSvc,
	}
}

func (h *LanZouHandler) Connect(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req struct {
		Cookie     string `json:"cookie" binding:"required"`
		TokenValue string `json:"token_value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.lanzouSvc.SaveToken(uid, req.Cookie, req.TokenValue); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "connected"})
}

func (h *LanZouHandler) GetStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	token, err := h.lanzouSvc.GetToken(uid)
	if err != nil {
		response.Success(c, gin.H{"connected": false})
		return
	}

	response.Success(c, gin.H{
		"connected":  true,
		"expires_at": token.ExpiresAt,
	})
}

// Disconnect 断开蓝奏云连接
func (h *LanZouHandler) Disconnect(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	if err := h.lanzouSvc.DeleteToken(uid); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "disconnected"})
}

// ListFiles 获取蓝奏云文件列表和文件夹列表
func (h *LanZouHandler) ListFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, _ := strconv.Atoi(c.Query("folder_id"))
	if folderID == 0 {
		folderID = -1 // 根目录
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	client, err := h.lanzouSvc.GetClient(uid)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// 获取文件列表
	fileResp, err := client.Task5(folderID, page)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取文件夹列表
	folderResp, err := client.Task47(folderID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"files":   fileResp.Text,
		"folders": folderResp.Text,
		"zt":      fileResp.Zt,
		"info":    fileResp.Info,
	})
}

// ListFolders 获取蓝奏云文件夹列表
func (h *LanZouHandler) ListFolders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, _ := strconv.Atoi(c.Query("folder_id"))
	if folderID == 0 {
		folderID = -1
	}

	client, err := h.lanzouSvc.GetClient(uid)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	resp, err := client.Task47(folderID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// CreateFolder 创建文件夹
func (h *LanZouHandler) CreateFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req struct {
		ParentID int    `json:"parent_id"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	client, err := h.lanzouSvc.GetClient(uid)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	resp, err := client.Task2(req.ParentID, req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// InitializeUpload 初始化上传
func (h *LanZouHandler) InitializeUpload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req service.InitializeUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.uploadSvc.InitializeUpload(uid, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// CompleteUpload 完成上传
func (h *LanZouHandler) CompleteUpload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid session id")
		return
	}

	var req service.CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, err := h.uploadSvc.CompleteUpload(uid, uint(sessionID), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, file)
}

// UploadStatus 获取上传状态
func (h *LanZouHandler) UploadStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid session id")
		return
	}

	session, err := h.uploadSvc.GetUploadStatus(uid, uint(sessionID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, session)
}

// CreateShare 创建分享链接
func (h *LanZouHandler) CreateShare(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req struct {
		FileID  int `json:"file_id" binding:"required"`
		Minutes int `json:"minutes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.lanzouSvc.CreateShare(uid, req.FileID, req.Minutes)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetFileURL 获取文件下载直链
func (h *LanZouHandler) GetFileURL(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	_, url, err := h.lanzouSvc.GetFileURL(uid, fileID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}

// SetAccess 设置文件/文件夹访问密码
func (h *LanZouHandler) SetAccess(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req struct {
		ID        int    `json:"id" binding:"required"`
		Type      string `json:"type" binding:"required"` // "file" 或 "folder"
		Shows     int    `json:"shows"`                   // 1=公开, 2=密码
		Shownames string `json:"shownames"`               // 密码
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var err error
	if req.Type == "file" {
		_, err = h.lanzouSvc.SetFileAccess(uid, req.ID, req.Shows, req.Shownames)
	} else if req.Type == "folder" {
		_, err = h.lanzouSvc.SetFolderAccess(uid, req.ID, req.Shows, req.Shownames)
	} else {
		response.BadRequest(c, "invalid type, must be 'file' or 'folder'")
		return
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "access updated"})
}

// GetFileDescription 获取文件描述
func (h *LanZouHandler) GetFileDescription(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	resp, err := h.lanzouSvc.GetFileDescription(uid, fileID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// SetFileDescription 设置文件描述
func (h *LanZouHandler) SetFileDescription(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	var req struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.lanzouSvc.SetFileDescription(uid, fileID, req.Description)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Rename 重命名文件/文件夹
func (h *LanZouHandler) Rename(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Type string `json:"type" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	client, err := h.lanzouSvc.GetClient(uid)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	var resp interface{}
	if req.Type == "file" {
		resp, err = client.Task14(id, req.Name)
	} else if req.Type == "folder" {
		resp, err = client.Task48FolderRename(id, req.Name)
	} else {
		response.BadRequest(c, "invalid type, must be 'file' or 'folder'")
		return
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}
