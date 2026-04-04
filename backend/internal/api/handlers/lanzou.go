package handlers

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type LanZouHandler struct {
	lanzouSvc *service.LanZouService
}

func NewLanZouHandler(lanzouSvc *service.LanZouService) *LanZouHandler {
	return &LanZouHandler{lanzouSvc: lanzouSvc}
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
