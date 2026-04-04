package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type FileVersionRepository struct {
	db *gorm.DB
}

func NewFileVersionRepository(db *gorm.DB) *FileVersionRepository {
	return &FileVersionRepository{db: db}
}

func (r *FileVersionRepository) Create(version *model.FileVersion) error {
	return r.db.Create(version).Error
}

func (r *FileVersionRepository) FindByFileID(fileID uint) ([]model.FileVersion, error) {
	var versions []model.FileVersion
	err := r.db.Where("file_id = ?", fileID).Order("created_at DESC").Find(&versions).Error
	return versions, err
}

func (r *FileVersionRepository) FindByID(id uint) (*model.FileVersion, error) {
	var version model.FileVersion
	err := r.db.First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *FileVersionRepository) Delete(id uint) error {
	return r.db.Delete(&model.FileVersion{}, id).Error
}

func (r *FileVersionRepository) DeleteByFileID(fileID uint) error {
	return r.db.Where("file_id = ?", fileID).Delete(&model.FileVersion{}).Error
}
