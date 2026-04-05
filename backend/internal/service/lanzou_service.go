package service

import (
	"errors"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

// LanZouTokenRepository 蓝奏云token仓库接口
type LanZouTokenRepository interface {
	Upsert(token *model.LanZouToken) error
	FindByUserID(userID uint) (*model.LanZouToken, error)
	DeleteByUserID(userID uint) error
}

type LanZouService struct {
	tokenRepo LanZouTokenRepository
	baseURL   string // 用于测试覆盖默认baseURL
}

func NewLanZouService(tokenRepo LanZouTokenRepository) *LanZouService {
	return &LanZouService{tokenRepo: tokenRepo}
}

// SetBaseURL 设置蓝奏云API基础URL（用于测试）
func (s *LanZouService) SetBaseURL(url string) {
	s.baseURL = url
}

func (s *LanZouService) SaveToken(userID uint, cookie, tokenValue string) error {
	token := &model.LanZouToken{
		UserID:     userID,
		Cookie:     cookie,
		TokenValue: tokenValue,
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	}
	return s.tokenRepo.Upsert(token)
}

func (s *LanZouService) GetToken(userID uint) (*model.LanZouToken, error) {
	return s.tokenRepo.FindByUserID(userID)
}

func (s *LanZouService) DeleteToken(userID uint) error {
	return s.tokenRepo.DeleteByUserID(userID)
}

// IsConnected 检查用户是否已连接蓝奏云
func (s *LanZouService) IsConnected(userID uint) bool {
	token, err := s.tokenRepo.FindByUserID(userID)
	if err != nil {
		return false
	}
	// 检查token是否过期
	if time.Now().After(token.ExpiresAt) {
		return false
	}
	return true
}

// GetClient 获取蓝奏云客户端
func (s *LanZouService) GetClient(userID uint) (*lanzou.Client, error) {
	token, err := s.tokenRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("lanzou not connected")
	}
	if token == nil {
		return nil, errors.New("lanzou not connected")
	}

	// 检查token是否过期
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("lanzou token expired")
	}

	client := lanzou.NewClient(token.Cookie)
	if s.baseURL != "" {
		client.SetBaseURL(s.baseURL)
	}
	return client, nil
}

// CreateShare 创建分享链接
func (s *LanZouService) CreateShare(userID uint, fileId int, minutes int) (*lanzou.ShareResponse, error) {
	client, err := s.GetClient(userID)
	if err != nil {
		return nil, err
	}
	return client.Task39(fileId, minutes)
}

// GetFileURL 获取文件下载直链
func (s *LanZouService) GetFileURL(userID uint, fileId int) (string, string, error) {
	client, err := s.GetClient(userID)
	if err != nil {
		return "", "", err
	}
	// 蓝奏云下载需要先获取文件详情
	resp, err := client.Task22(fileId)
	if err != nil {
		return "", "", err
	}

	// 从响应中提取直链
	// 蓝奏云API返回的down_url字段
	if url, ok := resp["down_url"].(string); ok {
		return "", url, nil
	}
	return "", "", errors.New("failed to get download URL")
}

// SetFileAccess 设置文件访问密码
func (s *LanZouService) SetFileAccess(userID uint, fileId int, shows int, shownames string) (*lanzou.Task23Response, error) {
	client, err := s.GetClient(userID)
	if err != nil {
		return nil, err
	}
	return client.Task23(fileId, shows, shownames)
}

// SetFolderAccess 设置文件夹访问密码
func (s *LanZouService) SetFolderAccess(userID uint, folderId int, shows int, shownames string) (*lanzou.Task16Response, error) {
	client, err := s.GetClient(userID)
	if err != nil {
		return nil, err
	}
	return client.Task16(folderId, shows, shownames)
}
