package repository

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(folder *model.Folder) error {
	return r.db.Create(folder).Error
}

func (r *FolderRepository) FindByID(id uint) (*model.Folder, error) {
	var folder model.Folder
	err := r.db.First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepository) FindByUserID(userID uint) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&folders).Error
	return folders, err
}

func (r *FolderRepository) FindByParentID(userID uint, parentID *uint) ([]model.Folder, error) {
	var folders []model.Folder
	q := r.db.Where("user_id = ?", userID)
	if parentID != nil {
		q = q.Where("parent_id = ?", *parentID)
	} else {
		q = q.Where("parent_id IS NULL")
	}
	err := q.Order("created_at DESC").Find(&folders).Error
	return folders, err
}

func (r *FolderRepository) Update(folder *model.Folder) error {
	return r.db.Save(folder).Error
}

func (r *FolderRepository) Delete(id uint) error {
	return r.db.Delete(&model.Folder{}, id).Error
}

func (r *FolderRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.Folder{}).Error
}

func (r *FolderRepository) DeleteByParentID(parentID uint) error {
	return r.db.Where("parent_id = ?", parentID).Delete(&model.Folder{}).Error
}
