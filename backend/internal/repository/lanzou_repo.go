package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type LanZouTokenRepository struct {
	db *gorm.DB
}

func NewLanZouTokenRepository(db *gorm.DB) *LanZouTokenRepository {
	return &LanZouTokenRepository{db: db}
}

func (r *LanZouTokenRepository) Upsert(token *model.LanZouToken) error {
	return r.db.Save(token).Error
}

func (r *LanZouTokenRepository) FindByUserID(userID uint) (*model.LanZouToken, error) {
	var token model.LanZouToken
	err := r.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *LanZouTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.LanZouToken{}).Error
}
