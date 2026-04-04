package service

import (
	"errors"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"
)

type LanZouService struct {
	tokenRepo *repository.LanZouTokenRepository
}

func NewLanZouService(tokenRepo *repository.LanZouTokenRepository) *LanZouService {
	return &LanZouService{tokenRepo: tokenRepo}
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

	// 检查token是否过期
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("lanzou token expired")
	}

	client := lanzou.NewClient(token.Cookie)
	return client, nil
}
