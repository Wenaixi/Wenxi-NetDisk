package handlers

import (
	"net/http"
	"strings"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"omitempty,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 如果未提供用户名，使用邮箱前缀作为用户名
	if req.Username == "" {
		req.Username = req.Email[:strings.Index(req.Email, "@")]
	}

	user, err := h.authSvc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 优先使用email，如果没有email则使用username
	identifier := req.Email
	if identifier == "" {
		identifier = req.Username
	}
	if identifier == "" {
		response.BadRequest(c, "email or username is required")
		return
	}

	token, err := h.authSvc.Login(identifier, req.Password)
	if err != nil {
		response.Unauthorized(c, "invalid credentials")
		return
	}

	response.Success(c, gin.H{
		"access_token": token,
		"token_type":   "bearer",
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uint)

	user, err := h.authSvc.GetUserByID(id)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
