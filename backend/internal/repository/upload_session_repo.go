package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type UploadSessionRepository struct {
	db *gorm.DB
}

func NewUploadSessionRepository(db *gorm.DB) *UploadSessionRepository {
	return &UploadSessionRepository{db: db}
}

func (r *UploadSessionRepository) Create(session *model.UploadSession) error {
	return r.db.Create(session).Error
}

func (r *UploadSessionRepository) FindByID(id uint) (*model.UploadSession, error) {
	var session model.UploadSession
	err := r.db.First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UploadSessionRepository) FindByUserIDAndHash(userID uint, fileHash string) (*model.UploadSession, error) {
	var session model.UploadSession
	err := r.db.Where("user_id = ? AND file_hash = ?", userID, fileHash).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UploadSessionRepository) Update(session *model.UploadSession) error {
	return r.db.Save(session).Error
}

func (r *UploadSessionRepository) Delete(id uint) error {
	return r.db.Delete(&model.UploadSession{}, id).Error
}

func (r *UploadSessionRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UploadSession{}).Error
}
