package service

import (
	"errors"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"
)

type FileService struct {
	fileRepo *repository.FileRepository
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	return &FileService{fileRepo: fileRepo}
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
