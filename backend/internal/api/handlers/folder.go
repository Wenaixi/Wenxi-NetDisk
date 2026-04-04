package handlers

import (
	"strconv"

	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	folderSvc *service.FolderService
}

func NewFolderHandler(folderSvc *service.FolderService) *FolderHandler {
	return &FolderHandler{folderSvc: folderSvc}
}

// ListFolders 获取文件夹列表
func (h *FolderHandler) ListFolders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	parentIDStr := c.Query("parent_id")
	var parentID *uint
	if parentIDStr != "" {
		id, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid parent_id")
			return
		}
		parentID = new(uint)
		*parentID = uint(id)
	}

	folders, err := h.folderSvc.ListFoldersByParent(uid, parentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, folders)
}

// CreateFolder 创建文件夹
func (h *FolderHandler) CreateFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req service.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folder, err := h.folderSvc.CreateFolder(uid, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, folder)
}

// GetFolder 获取文件夹详情
func (h *FolderHandler) GetFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid folder id")
		return
	}

	folder, err := h.folderSvc.GetFolder(uid, uint(folderID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, folder)
}

// UpdateFolder 更新文件夹
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid folder id")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folder, err := h.folderSvc.UpdateFolder(uid, uint(folderID), req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, folder)
}

// DeleteFolder 删除文件夹
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid folder id")
		return
	}

	if err := h.folderSvc.DeleteFolder(uid, uint(folderID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "folder deleted"})
}

// MoveFolder 移动文件夹
func (h *FolderHandler) MoveFolder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	folderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid folder id")
		return
	}

	var req struct {
		ParentID *uint `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folder, err := h.folderSvc.MoveFolder(uid, uint(folderID), req.ParentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, folder)
}
