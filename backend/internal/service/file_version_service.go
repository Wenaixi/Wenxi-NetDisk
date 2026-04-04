package service

import (
	"errors"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

type FileVersionService struct {
	versionRepo FileVersionRepository
	fileRepo    VersionFileRepository
}

type FileVersionRepository interface {
	Create(version *model.FileVersion) error
	FindByFileID(fileID uint) ([]model.FileVersion, error)
	FindByID(id uint) (*model.FileVersion, error)
	Delete(id uint) error
	DeleteByFileID(fileID uint) error
}

type VersionFileRepository interface {
	FindByID(id uint) (*model.File, error)
	Update(file *model.File) error
}

func NewFileVersionService(versionRepo FileVersionRepository, fileRepo VersionFileRepository) *FileVersionService {
	return &FileVersionService{
		versionRepo: versionRepo,
		fileRepo:    fileRepo,
	}
}

// CreateFileVersionRequest 创建文件版本请求
type CreateFileVersionRequest struct {
	LanZouFileID    string `json:"lanzou_file_id" binding:"required"`
	Size            int64  `json:"size" binding:"required"`
	EncryptionKey   string `json:"encryption_key" binding:"required"`
	EncryptionNonce string `json:"encryption_nonce" binding:"required"`
	Description     string `json:"description"`
}

// CreateVersion 创建文件版本
func (s *FileVersionService) CreateVersion(userID uint, fileID uint, req *CreateFileVersionRequest) (*model.FileVersion, error) {
	// 验证文件所有权
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	version := &model.FileVersion{
		FileID:          fileID,
		UserID:          userID,
		LanZouFileID:    req.LanZouFileID,
		Size:            req.Size,
		EncryptionKey:   req.EncryptionKey,
		EncryptionNonce: req.EncryptionNonce,
		Description:     req.Description,
	}

	if err := s.versionRepo.Create(version); err != nil {
		return nil, err
	}

	return version, nil
}

// ListVersions 获取文件的所有版本
func (s *FileVersionService) ListVersions(userID uint, fileID uint) ([]model.FileVersion, error) {
	// 验证文件所有权
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	return s.versionRepo.FindByFileID(fileID)
}

// RestoreVersion 恢复指定版本
func (s *FileVersionService) RestoreVersion(userID uint, fileID uint, versionID uint) (*model.File, error) {
	// 验证文件所有权
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 获取版本
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, errors.New("version not found")
	}
	if version.FileID != fileID {
		return nil, errors.New("version does not belong to file")
	}

	// 恢复版本到当前文件
	file.LanZouFileID = version.LanZouFileID
	file.Size = version.Size
	file.EncryptionKey = version.EncryptionKey
	file.EncryptionNonce = version.EncryptionNonce

	if err := s.fileRepo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteVersion 删除指定版本
func (s *FileVersionService) DeleteVersion(userID uint, fileID uint, versionID uint) error {
	// 验证文件所有权
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}
	if file.UserID != userID {
		return errors.New("access denied")
	}

	// 验证版本存在
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return errors.New("version not found")
	}
	if version.FileID != fileID {
		return errors.New("version does not belong to file")
	}

	return s.versionRepo.Delete(versionID)
}
