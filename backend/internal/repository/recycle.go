package repository

import (
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
)

type RecycleBinRepository struct {
	db *gorm.DB
}

func NewRecycleBinRepository(db *gorm.DB) *RecycleBinRepository {
	return &RecycleBinRepository{db: db}
}

// Create 创建回收站记录
func (r *RecycleBinRepository) Create(item *model.RecycleBin) error {
	return r.db.Create(item).Error
}

// List 获取用户回收站列表
func (r *RecycleBinRepository) List(userID uint) ([]model.RecycleBin, error) {
	var items []model.RecycleBin
	err := r.db.Where("user_id = ? AND deleted_at >= ?", userID, getExpiryDate()).
		Order("deleted_at DESC").
		Find(&items).Error
	return items, err
}

// ListPaginated 分页获取回收站列表
func (r *RecycleBinRepository) ListPaginated(userID uint, page, pageSize int) ([]model.RecycleBin, int64, error) {
	var items []model.RecycleBin
	var total int64

	query := r.db.Model(&model.RecycleBin{}).Where("user_id = ? AND deleted_at >= ?", userID, getExpiryDate())

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("deleted_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetByID 根据ID获取回收站记录
func (r *RecycleBinRepository) GetByID(id uint, userID uint) (*model.RecycleBin, error) {
	var item model.RecycleBin
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Restore 恢复回收站记录（删除回收站条目）
func (r *RecycleBinRepository) Restore(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecycleBin{}).Error
}

// DeletePermanently 永久删除回收站记录
func (r *RecycleBinRepository) DeletePermanently(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Unscoped().Delete(&model.RecycleBin{}).Error
}

// CleanExpired 清理过期的回收站记录
func (r *RecycleBinRepository) CleanExpired() error {
	return r.db.Where("deleted_at < ?", getExpiryDate()).Unscoped().Delete(&model.RecycleBin{}).Error
}

// ClearAll 清空用户回收站
func (r *RecycleBinRepository) ClearAll(userID uint) error {
	return r.db.Where("user_id = ?", userID).Unscoped().Delete(&model.RecycleBin{}).Error
}

// GetByItemType 根据类型获取回收站记录
func (r *RecycleBinRepository) GetByItemType(userID uint, itemType string) ([]model.RecycleBin, error) {
	var items []model.RecycleBin
	err := r.db.Where("user_id = ? AND item_type = ? AND deleted_at >= ?", userID, itemType, getExpiryDate()).
		Order("deleted_at DESC").
		Find(&items).Error
	return items, err
}

// getExpiryDate 获取回收站过期时间(30天前)
func getExpiryDate() interface{} {
	// 返回30天前的时间，软删除的记录如果早于此时间将被清理
	return time.Now().AddDate(0, 0, -30)
}
