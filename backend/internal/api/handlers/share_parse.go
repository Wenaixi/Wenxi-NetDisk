package handlers

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ShareParseHandler struct {
	parseSvc *service.ShareParseService
}

func NewShareParseHandler(svc *service.ShareParseService) *ShareParseHandler {
	return &ShareParseHandler{parseSvc: svc}
}

// ParseShare 解析分享链接
func (h *ShareParseHandler) ParseShare(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
		Pwd string `json:"pwd"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "分享链接不能为空")
		return
	}

	result, err := h.parseSvc.ParseShareLink(req.URL, req.Pwd)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetShareDownloadURL 获取分享文件的下载直链
func (h *ShareParseHandler) GetShareDownloadURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
		Pwd string `json:"pwd"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "分享链接不能为空")
		return
	}

	downloadURL, err := h.parseSvc.GetShareDownloadLink(req.URL, req.Pwd)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"download_url": downloadURL})
}

// ValidateShareURL 验证分享链接是否有效
func (h *ShareParseHandler) ValidateShareURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "分享链接不能为空")
		return
	}

	valid := h.parseSvc.ValidateShareLink(req.URL)
	if !valid {
		response.BadRequest(c, "无效的蓝奏云分享链接")
		return
	}

	response.Success(c, gin.H{"valid": true})
}
