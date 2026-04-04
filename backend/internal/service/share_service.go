package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/crypto"
)

// ShareRepository 分享仓库接口
type ShareRepository interface {
	Create(share *model.Share) error
	FindByID(id uint) (*model.Share, error)
	FindByToken(token string) (*model.Share, error)
	FindByFileID(fileID uint) ([]model.Share, error)
	Delete(id uint) error
	DeleteByFileID(fileID uint) error
}

// ShareFileRepository 分享服务所需的文件仓库接口
type ShareFileRepository interface {
	FindByID(id uint) (*model.File, error)
}

type ShareService struct {
	shareRepo ShareRepository
	fileRepo  ShareFileRepository
}

func NewShareService(shareRepo ShareRepository, fileRepo ShareFileRepository) *ShareService {
	return &ShareService{
		shareRepo: shareRepo,
		fileRepo:  fileRepo,
	}
}

func (s *ShareService) CreateShare(userID, fileID uint, password *string, expiresAt *string) (*model.Share, error) {
	// Verify file ownership
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}
	if file.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Generate share token
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	shareToken := hex.EncodeToString(tokenBytes)

	var passwordHash *string
	if password != nil && *password != "" {
		hash, err := crypto.HashPassword(*password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hash
	}

	var expires *time.Time
	if expiresAt != nil && *expiresAt != "" {
		t, err := time.Parse(time.RFC3339, *expiresAt)
		if err == nil {
			expires = &t
		}
	}

	share := &model.Share{
		FileID:        fileID,
		ShareToken:    shareToken,
		PasswordHash:  passwordHash,
		ExpiresAt:     expires,
	}

	if err := s.shareRepo.Create(share); err != nil {
		return nil, err
	}

	return share, nil
}

func (s *ShareService) GetShareByToken(token string) (*model.Share, error) {
	return s.shareRepo.FindByToken(token)
}

func (s *ShareService) DeleteShare(userID, shareID uint) error {
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		return err
	}

	// Verify ownership through file
	file, err := s.fileRepo.FindByID(share.FileID)
	if err != nil {
		return err
	}
	if file.UserID != userID {
		return errors.New("access denied")
	}

	return s.shareRepo.Delete(shareID)
}

// ValidatePassword 验证分享密码
func (s *ShareService) ValidatePassword(share *model.Share, password string) bool {
	if share.PasswordHash == nil {
		return true // 无密码分享
	}
	return crypto.CheckPassword(password, *share.PasswordHash)
}

// IsExpired 检查分享是否过期
func (s *ShareService) IsExpired(share *model.Share) bool {
	if share.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*share.ExpiresAt)
}
