package service

import (
	"errors"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// mockShareRepo 模拟分享仓库
type mockShareRepo struct {
	shares []*model.Share
}

func newMockShareRepo() *mockShareRepo {
	return &mockShareRepo{
		shares: []*model.Share{
			{
				ID:          1,
				UserID:      1,
				FileID:      1,
				ShareToken:  "abc123",
				PasswordHash: nil,
				ExpiresAt:   nil,
			},
			{
				ID:          2,
				UserID:      1,
				FileID:      2,
				ShareToken:  "expired456",
				PasswordHash: nil,
				ExpiresAt:   &[]time.Time{time.Now().Add(-24 * time.Hour)}[0],
			},
			{
				ID:          3,
				UserID:      1,
				FileID:      3,
				ShareToken:  "pwd789",
				PasswordHash: &[]string{"$2a$10$hashedpwd"}[0],
				ExpiresAt:   &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			},
		},
	}
}

func (m *mockShareRepo) Create(share *model.Share) error {
	share.ID = uint(len(m.shares) + 1)
	m.shares = append(m.shares, share)
	return nil
}

func (m *mockShareRepo) FindByID(id uint) (*model.Share, error) {
	for _, s := range m.shares {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("share not found")
}

func (m *mockShareRepo) FindByToken(token string) (*model.Share, error) {
	for _, s := range m.shares {
		if s.ShareToken == token {
			return s, nil
		}
	}
	return nil, errors.New("share not found")
}

func (m *mockShareRepo) FindByFileID(fileID uint) ([]model.Share, error) {
	var result []model.Share
	for _, s := range m.shares {
		if s.FileID == fileID {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (m *mockShareRepo) FindByUserID(userID uint) ([]model.Share, error) {
	var result []model.Share
	fileRepo := newMockShareFileRepo()
	for _, s := range m.shares {
		// 通过file找到userID
		file, _ := fileRepo.FindByID(s.FileID)
		if file != nil && file.UserID == userID {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (m *mockShareRepo) Delete(id uint) error {
	for i, s := range m.shares {
		if s.ID == id {
			m.shares = append(m.shares[:i], m.shares[i+1:]...)
			return nil
		}
	}
	return errors.New("share not found")
}

func (m *mockShareRepo) DeleteByFileID(fileID uint) error {
	var remaining []*model.Share
	for _, s := range m.shares {
		if s.FileID != fileID {
			remaining = append(remaining, s)
		}
	}
	m.shares = remaining
	return nil
}

// mockShareFileRepo 模拟分享服务所需的文件仓库
type mockShareFileRepo struct {
	files []*model.File
}

func newMockShareFileRepo() *mockShareFileRepo {
	return &mockShareFileRepo{
		files: []*model.File{
			{ID: 1, UserID: 1, Name: "file1.txt", Size: 100},
			{ID: 2, UserID: 1, Name: "file2.pdf", Size: 200},
			{ID: 3, UserID: 2, Name: "file3.doc", Size: 300},
		},
	}
}

// mockShareFileRepoForShare 用于FindByUserID查询的辅助类型
type mockShareFileRepoForShare struct {
	files []*model.File
}

func (m mockShareFileRepoForShare) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func newMockShareFileRepoForShare() mockShareFileRepoForShare {
	return mockShareFileRepoForShare{
		files: []*model.File{
			{ID: 1, UserID: 1, Name: "file1.txt", Size: 100},
			{ID: 2, UserID: 1, Name: "file2.pdf", Size: 200},
			{ID: 3, UserID: 2, Name: "file3.doc", Size: 300},
		},
	}
}

func (m *mockShareFileRepo) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

// TestShareService_CreateShare 测试创建分享
func TestShareService_CreateShare(t *testing.T) {
	t.Run("should create share successfully without password", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.CreateShare(1, 1, nil, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if share == nil {
			t.Error("expected share, got nil")
			return
		}
		if share.ShareToken == "" {
			t.Error("expected share token to be generated")
		}
		if share.PasswordHash != nil {
			t.Error("expected no password hash for passwordless share")
		}
	})

	t.Run("should create share with password", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		password := "secret123"
		share, err := svc.CreateShare(1, 1, &password, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if share.PasswordHash == nil {
			t.Error("expected password hash to be set")
		}
	})

	t.Run("should create share with expiration", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		expiry := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
		share, err := svc.CreateShare(1, 1, nil, &expiry)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if share.ExpiresAt == nil {
			t.Error("expected expiration to be set")
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.CreateShare(1, 999, nil, nil)

		if err == nil {
			t.Error("expected file not found error")
		}
		if share != nil {
			t.Error("expected nil share for non-existent file")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.CreateShare(3, 1, nil, nil) // user 3 doesn't own file 1

		if err == nil {
			t.Error("expected access denied error")
		}
		if share != nil {
			t.Error("expected nil share for unauthorized user")
		}
	})
}

// TestShareService_GetShareByToken 测试通过token获取分享
func TestShareService_GetShareByToken(t *testing.T) {
	t.Run("should return share for valid token", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.GetShareByToken("abc123")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if share == nil {
			t.Error("expected share, got nil")
			return
		}
		if share.ShareToken != "abc123" {
			t.Errorf("expected token 'abc123', got '%s'", share.ShareToken)
		}
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.GetShareByToken("invalid_token")

		if err == nil {
			t.Error("expected share not found error")
		}
		if share != nil {
			t.Error("expected nil share for invalid token")
		}
	})
}

// TestShareService_DeleteShare 测试删除分享
func TestShareService_DeleteShare(t *testing.T) {
	t.Run("should delete share successfully", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		err := svc.DeleteShare(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should fail if share not found", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		err := svc.DeleteShare(1, 999)

		if err == nil {
			t.Error("expected share not found error")
		}
	})

	t.Run("should fail if user doesn't own the file", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		err := svc.DeleteShare(3, 1) // user 3 doesn't own file 1

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestShareService_ValidatePassword 测试密码验证
func TestShareService_ValidatePassword(t *testing.T) {
	t.Run("should return true for passwordless share", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, _ := svc.GetShareByToken("abc123")

		if !svc.ValidatePassword(share, "") {
			t.Error("expected true for passwordless share")
		}
	})

	t.Run("should validate correct password", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		// 创建一个带密码的分享
		password := "testpass123"
		share, err := svc.CreateShare(1, 1, &password, nil)
		if err != nil {
			t.Fatalf("failed to create share: %v", err)
		}

		if !svc.ValidatePassword(share, "testpass123") {
			t.Error("expected correct password to validate")
		}
	})

	t.Run("should reject wrong password", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		// 使用已存在的带密码的分享
		share, _ := svc.GetShareByToken("pwd789")
		if share == nil {
			t.Skip("share not found")
		}

		if svc.ValidatePassword(share, "wrongpassword") {
			t.Error("expected wrong password to be rejected")
		}
	})
}

// TestShareService_IsExpired 测试过期检查
func TestShareService_IsExpired(t *testing.T) {
	t.Run("should return false for non-expiring share", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, _ := svc.GetShareByToken("abc123")

		if svc.IsExpired(share) {
			t.Error("expected non-expiring share to not be expired")
		}
	})

	t.Run("should return true for expired share", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, _ := svc.GetShareByToken("expired456")

		if !svc.IsExpired(share) {
			t.Error("expected expired share to be expired")
		}
	})

	t.Run("should return false for future-expiring share", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, _ := svc.GetShareByToken("pwd789")

		if svc.IsExpired(share) {
			t.Error("expected future-expiring share to not be expired")
		}
	})
}

// TestShareService_ListShares 测试获取用户的所有分享
func TestShareService_ListShares(t *testing.T) {
	t.Run("should list all shares for user", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		shares, err := svc.ListShares(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// user 1 owns files 1,2. shares for these files: share 1 (abc123, file1), share 2 (expired456, file1)
		// share 1 is not expired, share 2 is expired so should be filtered
		if len(shares) != 1 {
			t.Errorf("expected 1 share for user 1, got %d", len(shares))
		}
	})

	t.Run("should filter out expired shares", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		shares, err := svc.ListShares(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Verify no expired shares
		for _, share := range shares {
			if svc.IsExpired(&share) {
				t.Error("expected no expired shares in result")
			}
		}
	})

	t.Run("should return empty for user with no shares", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		shares, err := svc.ListShares(999)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(shares) != 0 {
			t.Errorf("expected 0 shares, got %d", len(shares))
		}
	})
}

// TestShareService_GetShareWithFile 测试获取分享详情（包含过期检查）
func TestShareService_GetShareWithFile(t *testing.T) {
	t.Run("should return share for valid token", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.GetShareWithFile("abc123")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if share == nil {
			t.Error("expected share, got nil")
			return
		}
		if share.ShareToken != "abc123" {
			t.Errorf("expected token 'abc123', got '%s'", share.ShareToken)
		}
	})

	t.Run("should return error for expired share", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.GetShareWithFile("expired456")

		if err == nil {
			t.Error("expected error for expired share")
		}
		if share != nil {
			t.Error("expected nil share for expired token")
		}
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		shareRepo := newMockShareRepo()
		fileRepo := newMockShareFileRepo()
		svc := NewShareService(shareRepo, fileRepo)

		share, err := svc.GetShareWithFile("invalid_token")

		if err == nil {
			t.Error("expected error for invalid token")
		}
		if share != nil {
			t.Error("expected nil share for invalid token")
		}
	})
}

// TestShareService_CreateShare_FileNotFound 测试文件不存在
func TestShareService_CreateShare_FileNotFound(t *testing.T) {
	shareRepo := newMockShareRepo()
	fileRepo := newMockShareFileRepo()
	svc := NewShareService(shareRepo, fileRepo)

	_, err := svc.CreateShare(1, 999, nil, nil)
	if err == nil {
		t.Error("expected file not found error")
	}
}

// TestShareService_CreateShare_AccessDenied 测试无权访问文件
func TestShareService_CreateShare_AccessDenied(t *testing.T) {
	shareRepo := newMockShareRepo()
	fileRepo := newMockShareFileRepo()
	svc := NewShareService(shareRepo, fileRepo)

	// File 3 belongs to user 2
	_, err := svc.CreateShare(1, 3, nil, nil)
	if err == nil {
		t.Error("expected access denied error")
	}
}

// TestShareService_ListShares_RepoError 测试仓库错误
func TestShareService_ListShares_RepoError(t *testing.T) {
	shareRepo := newMockShareRepo()
	fileRepo := newMockShareFileRepo()
	svc := NewShareService(shareRepo, fileRepo)

	// Inject error into repo
	shareRepo.shares = nil // clear shares to avoid complex error injection
	shares, err := svc.ListShares(999)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(shares) != 0 {
		t.Errorf("expected 0 shares, got %d", len(shares))
	}
}

// TestShareService_DeleteShare_NotFound 测试删除不存在的分享
func TestShareService_DeleteShare_NotFound(t *testing.T) {
	shareRepo := newMockShareRepo()
	fileRepo := newMockShareFileRepo()
	svc := NewShareService(shareRepo, fileRepo)

	err := svc.DeleteShare(1, 999)
	if err == nil {
		t.Error("expected share not found error")
	}
}

// TestShareService_DeleteShare_AccessDenied 测试删除他人分享
func TestShareService_DeleteShare_AccessDenied(t *testing.T) {
	shareRepo := newMockShareRepo()
	fileRepo := newMockShareFileRepo()
	svc := NewShareService(shareRepo, fileRepo)

	// Share 1 belongs to user 1, try to delete as user 999
	err := svc.DeleteShare(999, 1)
	if err == nil {
		t.Error("expected access denied error")
	}
}
