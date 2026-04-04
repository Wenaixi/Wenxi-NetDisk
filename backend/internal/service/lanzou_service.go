package service

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
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
	}
	return s.tokenRepo.Upsert(token)
}

func (s *LanZouService) GetToken(userID uint) (*model.LanZouToken, error) {
	return s.tokenRepo.FindByUserID(userID)
}

func (s *LanZouService) DeleteToken(userID uint) error {
	return s.tokenRepo.DeleteByUserID(userID)
}
