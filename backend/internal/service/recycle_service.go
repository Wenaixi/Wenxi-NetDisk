package service

import (
	"errors"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

type RecycleBinService struct {
	recycleRepo RecycleBinRepository
	fileRepo    RecycleFileRepo
	folderRepo  RecycleFolderRepo
}

type RecycleBinRepository interface {
	Create(item *model.RecycleBin) error
	List(userID uint) ([]model.RecycleBin, error)
	GetByID(id uint, userID uint) (*model.RecycleBin, error)
	Restore(id uint, userID uint) error
	DeletePermanently(id uint, userID uint) error
	ClearAll(userID uint) error
}

type RecycleFileRepo interface {
	FindByID(id uint) (*model.File, error)
	Create(file *model.File) error
	Delete(id uint) error
}

type RecycleFolderRepo interface {
	FindByID(id uint) (*model.Folder, error)
	Create(folder *model.Folder) error
	Delete(id uint) error
}

func NewRecycleBinService(
	recycleRepo RecycleBinRepository,
	fileRepo RecycleFileRepo,
	folderRepo RecycleFolderRepo,
) *RecycleBinService {
	return &RecycleBinService{
		recycleRepo: recycleRepo,
		fileRepo:    fileRepo,
		folderRepo:  folderRepo,
	}
}

// MoveToRecycleBin 将文件移入回收站
func (s *RecycleBinService) MoveToRecycleBin(userID uint, fileID uint) (*model.RecycleBin, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 删除文件
	if err := s.fileRepo.Delete(fileID); err != nil {
		return nil, err
	}

	// 创建回收站记录
	item := &model.RecycleBin{
		UserID:          userID,
		OriginalName:    file.Name,
		ItemType:        "file",
		ItemID:          fileID,
		LanZouFileID:    file.LanZouFileID,
		Size:            file.Size,
		EncryptionKey:   file.EncryptionKey,
		EncryptionNonce: file.EncryptionNonce,
		DeletedAt:       time.Now(),
	}

	if err := s.recycleRepo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// MoveFolderToRecycleBin 将文件夹移入回收站
func (s *RecycleBinService) MoveFolderToRecycleBin(userID uint, folderID uint) (*model.RecycleBin, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}

	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 删除文件夹
	if err := s.folderRepo.Delete(folderID); err != nil {
		return nil, err
	}

	// 创建回收站记录
	item := &model.RecycleBin{
		UserID:         userID,
		OriginalName:   folder.Name,
		ItemType:       "folder",
		ItemID:         folderID,
		DeletedAt:      time.Now(),
	}

	if folder.ParentID != nil {
		item.ParentID = *folder.ParentID
	}

	if err := s.recycleRepo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// List 获取回收站列表
func (s *RecycleBinService) List(userID uint) ([]model.RecycleBin, error) {
	return s.recycleRepo.List(userID)
}

// Restore 恢复回收站项目
func (s *RecycleBinService) Restore(userID uint, recycleID uint) error {
	item, err := s.recycleRepo.GetByID(recycleID, userID)
	if err != nil {
		return errors.New("recycle item not found")
	}

	if item.UserID != userID {
		return errors.New("access denied")
	}

	// 根据类型恢复
	if item.ItemType == "file" {
		file := &model.File{
			ID:              item.ItemID,
			UserID:          userID,
			Name:            item.OriginalName,
			Size:            item.Size,
			LanZouFileID:    item.LanZouFileID,
			EncryptionKey:   item.EncryptionKey,
			EncryptionNonce: item.EncryptionNonce,
		}
		if err := s.fileRepo.Create(file); err != nil {
			return err
		}
	} else if item.ItemType == "folder" {
		folder := &model.Folder{
			ID:       item.ItemID,
			UserID:   userID,
			Name:     item.OriginalName,
		}
		if item.ParentID > 0 {
			pid := item.ParentID
			folder.ParentID = &pid
		}
		if err := s.folderRepo.Create(folder); err != nil {
			return err
		}
	}

	// 删除回收站记录
	return s.recycleRepo.Restore(recycleID, userID)
}

// PermanentDelete 永久删除
func (s *RecycleBinService) PermanentDelete(userID uint, recycleID uint) error {
	_, err := s.recycleRepo.GetByID(recycleID, userID)
	if err != nil {
		return errors.New("recycle item not found")
	}

	return s.recycleRepo.DeletePermanently(recycleID, userID)
}

// ClearAll 清空回收站
func (s *RecycleBinService) ClearAll(userID uint) error {
	return s.recycleRepo.ClearAll(userID)
}
