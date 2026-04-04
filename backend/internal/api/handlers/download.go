package handlers

import (
	"strconv"

	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	downloadSvc *service.DownloadService
}

func NewDownloadHandler(downloadSvc *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{downloadSvc: downloadSvc}
}

// GetDownloadURL 获取文件下载直链
func (h *DownloadHandler) GetDownloadURL(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	resp, err := h.downloadSvc.GetDownloadURL(uid, uint(fileID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}
