package handlers

import (
	"strconv"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileSvc *service.FileService
}

func NewFileHandler(fileSvc *service.FileService) *FileHandler {
	return &FileHandler{fileSvc: fileSvc}
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uint)

	files, err := h.fileSvc.ListFiles(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, files)
}

func (h *FileHandler) CreateFileMetadata(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uint)

	var req service.CreateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, err := h.fileSvc.CreateMetadata(id, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, file)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	file, err := h.fileSvc.GetFile(uid, uint(fileID))
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}

	response.Success(c, file)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	if err := h.fileSvc.DeleteFile(uid, uint(fileID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "file deleted"})
}
