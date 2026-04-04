package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type ShareRepository struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) *ShareRepository {
	return &ShareRepository{db: db}
}

func (r *ShareRepository) Create(share *model.Share) error {
	return r.db.Create(share).Error
}

func (r *ShareRepository) FindByID(id uint) (*model.Share, error) {
	var share model.Share
	err := r.db.First(&share, id).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *ShareRepository) FindByToken(token string) (*model.Share, error) {
	var share model.Share
	err := r.db.Preload("File").Where("share_token = ?", token).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *ShareRepository) FindByFileID(fileID uint) ([]model.Share, error) {
	var shares []model.Share
	err := r.db.Where("file_id = ?", fileID).Find(&shares).Error
	return shares, err
}

func (r *ShareRepository) FindByUserID(userID uint) ([]model.Share, error) {
	var shares []model.Share
	err := r.db.Preload("File").Where("user_id = ?", userID).Find(&shares).Error
	return shares, err
}

func (r *ShareRepository) Delete(id uint) error {
	return r.db.Delete(&model.Share{}, id).Error
}

func (r *ShareRepository) DeleteByFileID(fileID uint) error {
	return r.db.Where("file_id = ?", fileID).Delete(&model.Share{}).Error
}
