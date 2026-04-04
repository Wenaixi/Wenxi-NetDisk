package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) FindByID(id uint) (*model.File, error) {
	var file model.File
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) FindByUserID(userID uint) ([]model.File, error) {
	var files []model.File
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (r *FileRepository) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	var file model.File
	err := r.db.Where("lanzou_file_id = ?", lanzouFileID).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) Update(file *model.File) error {
	return r.db.Save(file).Error
}

func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&model.File{}, id).Error
}

func (r *FileRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.File{}).Error
}
