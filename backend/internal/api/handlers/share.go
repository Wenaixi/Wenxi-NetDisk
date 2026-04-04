package handlers

import (
	"net/http"
	"strconv"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	shareSvc *service.ShareService
}

func NewShareHandler(shareSvc *service.ShareService) *ShareHandler {
	return &ShareHandler{shareSvc: shareSvc}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	var req struct {
		Password  *string `json:"password"`
		ExpiresAt *string `json:"expires_at"`
	}
	c.ShouldBindJSON(&req)

	share, err := h.shareSvc.CreateShare(uid, uint(fileID), req.Password, req.ExpiresAt)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"share_token": share.ShareToken,
		"share_url":   "/api/shares/" + share.ShareToken,
	})
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	token := c.Param("token")

	share, err := h.shareSvc.GetShareWithFile(token)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":                share.ID,
		"file_id":          share.FileID,
		"file_name":        share.File.Name,
		"file_size":        share.File.Size,
		"requires_password": share.PasswordHash != nil,
		"expires_at":       share.ExpiresAt,
		"created_at":       share.CreatedAt,
	})
}

func (h *ShareHandler) ListShares(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	shares, err := h.shareSvc.ListShares(uid)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 转换为响应格式
	var result []gin.H
	for _, share := range shares {
		result = append(result, gin.H{
			"id":                share.ID,
			"file_id":          share.FileID,
			"file_name":        share.File.Name,
			"file_size":        share.File.Size,
			"requires_password": share.PasswordHash != nil,
			"expires_at":       share.ExpiresAt,
			"created_at":       share.CreatedAt,
		})
	}

	response.Success(c, result)
}

func (h *ShareHandler) DeleteShare(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	shareID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid share id")
		return
	}

	if err := h.shareSvc.DeleteShare(uid, uint(shareID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "share deleted"})
}
