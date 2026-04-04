package service

import (
	"errors"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// FolderRepository 文件夹仓库接口
type FolderRepository interface {
	Create(folder *model.Folder) error
	FindByID(id uint) (*model.Folder, error)
	FindByUserID(userID uint) ([]model.Folder, error)
	FindByParentID(userID uint, parentID *uint) ([]model.Folder, error)
	Update(folder *model.Folder) error
	Delete(id uint) error
	DeleteByUserID(userID uint) error
	DeleteByParentID(parentID uint) error
}

type FolderService struct {
	folderRepo FolderRepository
}

func NewFolderService(folderRepo FolderRepository) *FolderService {
	return &FolderService{folderRepo: folderRepo}
}

// CreateFolderRequest 创建文件夹请求
type CreateFolderRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

func (s *FolderService) CreateFolder(userID uint, req *CreateFolderRequest) (*model.Folder, error) {
	// 如果有父文件夹，验证其存在且属于该用户
	if req.ParentID != nil {
		parent, err := s.folderRepo.FindByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("parent folder not found")
		}
		if parent.UserID != userID {
			return nil, errors.New("access denied")
		}
	}

	folder := &model.Folder{
		UserID:   userID,
		ParentID: req.ParentID,
		Name:     req.Name,
	}

	if err := s.folderRepo.Create(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *FolderService) ListFolders(userID uint) ([]model.Folder, error) {
	return s.folderRepo.FindByUserID(userID)
}

func (s *FolderService) ListFoldersByParent(userID uint, parentID *uint) ([]model.Folder, error) {
	return s.folderRepo.FindByParentID(userID, parentID)
}

func (s *FolderService) GetFolder(userID, folderID uint) (*model.Folder, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}
	return folder, nil
}

func (s *FolderService) UpdateFolder(userID, folderID uint, name string) (*model.Folder, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}

	folder.Name = name
	if err := s.folderRepo.Update(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *FolderService) DeleteFolder(userID, folderID uint) error {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return err
	}
	if folder.UserID != userID {
		return errors.New("access denied")
	}

	// 删除该文件夹下的所有子文件夹
	if err := s.folderRepo.DeleteByParentID(folderID); err != nil {
		return err
	}

	return s.folderRepo.Delete(folderID)
}

func (s *FolderService) MoveFolder(userID, folderID uint, newParentID *uint) (*model.Folder, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 如果有新的父文件夹，验证其存在且属于该用户
	if newParentID != nil {
		parent, err := s.folderRepo.FindByID(*newParentID)
		if err != nil {
			return nil, errors.New("parent folder not found")
		}
		if parent.UserID != userID {
			return nil, errors.New("access denied")
		}

		// 防止将文件夹移动到自身或其子文件夹下
		if *newParentID == folderID {
			return nil, errors.New("cannot move folder to itself")
		}
		if s.isDescendant(folderID, *newParentID) {
			return nil, errors.New("cannot move folder to its descendant")
		}
	}

	folder.ParentID = newParentID
	if err := s.folderRepo.Update(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// isDescendant 检查 targetID 是否是 ancestorID 的后代
func (s *FolderService) isDescendant(ancestorID, targetID uint) bool {
	folder, err := s.folderRepo.FindByID(targetID)
	if err != nil {
		return false
	}
	if folder.ParentID == nil {
		return false
	}
	if *folder.ParentID == ancestorID {
		return true
	}
	return s.isDescendant(ancestorID, *folder.ParentID)
}
