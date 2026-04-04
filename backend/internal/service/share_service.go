package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type ShareService struct {
	shareRepo *repository.ShareRepository
	fileRepo  *repository.FileRepository
}

func NewShareService(shareRepo *repository.ShareRepository, fileRepo *repository.FileRepository) *ShareService {
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
		hash, err := hashPassword(*password)
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
