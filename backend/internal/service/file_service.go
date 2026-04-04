package service

import (
	"errors"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// FileRepository 文件仓库接口
type FileRepository interface {
	Create(file *model.File) error
	FindByUserID(userID uint) ([]model.File, error)
	FindByID(id uint) (*model.File, error)
	FindByFolderID(userID uint, folderID *uint) ([]model.File, error)
	FindByLanZouFileID(lanzouFileID string) (*model.File, error)
	Delete(id uint) error
	Update(file *model.File) error
}

type FileService struct {
	fileRepo FileRepository
}

func NewFileService(fileRepo FileRepository) *FileService {
	return &FileService{fileRepo: fileRepo}
}

// FindByID 根据ID查找文件（供其他服务使用）
func (s *FileService) FindByID(id uint) (*model.File, error) {
	return s.fileRepo.FindByID(id)
}

// FindByLanZouFileID 根据蓝奏云文件ID查找（供其他服务使用）
func (s *FileService) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return s.fileRepo.FindByLanZouFileID(lanzouFileID)
}

type CreateFileRequest struct {
	Name            string `json:"name" binding:"required"`
	Size            int64  `json:"size" binding:"required"`
	LanZouFileID    string `json:"lanzou_file_id" binding:"required"`
	LanZouFolderID  string `json:"lanzou_folder_id"`
	EncryptionKey   string `json:"encryption_key" binding:"required"`
	EncryptionNonce string `json:"encryption_nonce" binding:"required"`
	MimeType        string `json:"mime_type"`
}

func (s *FileService) CreateMetadata(userID uint, req *CreateFileRequest) (*model.File, error) {
	file := &model.File{
		UserID:           userID,
		Name:             req.Name,
		Size:             req.Size,
		LanZouFileID:     req.LanZouFileID,
		LanZouFolderID:   req.LanZouFolderID,
		EncryptionKey:    req.EncryptionKey,
		EncryptionNonce:  req.EncryptionNonce,
		MimeType:         req.MimeType,
	}

	if err := s.fileRepo.Create(file); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) ListFiles(userID uint) ([]model.File, error) {
	return s.fileRepo.FindByUserID(userID)
}

func (s *FileService) GetFile(userID, fileID uint) (*model.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}
	return file, nil
}

func (s *FileService) DeleteFile(userID, fileID uint) error {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return err
	}
	if file.UserID != userID {
		return errors.New("access denied")
	}
	return s.fileRepo.Delete(fileID)
}

// UpdateFileDescription 更新文件描述
func (s *FileService) UpdateFileDescription(userID, fileID uint, description string) (*model.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	file.Description = description
	if err := s.fileRepo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}
