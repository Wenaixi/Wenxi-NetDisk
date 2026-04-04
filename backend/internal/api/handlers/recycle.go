package handlers

import (
	"strconv"

	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type RecycleHandler struct {
	recycleSvc *service.RecycleBinService
}

func NewRecycleHandler(recycleSvc *service.RecycleBinService) *RecycleHandler {
	return &RecycleHandler{recycleSvc: recycleSvc}
}

// List 获取回收站列表
func (h *RecycleHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	items, err := h.recycleSvc.List(uid)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, items)
}

// Restore 恢复回收站项目
func (h *RecycleHandler) Restore(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.recycleSvc.Restore(uid, uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "restored"})
}

// Delete 永久删除回收站项目
func (h *RecycleHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.recycleSvc.PermanentDelete(uid, uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "deleted"})
}

// Clear 清空回收站
func (h *RecycleHandler) Clear(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	if err := h.recycleSvc.ClearAll(uid); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "cleared"})
}
