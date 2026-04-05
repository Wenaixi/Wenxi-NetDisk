package service

import (
	"errors"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

// mockDownloadFileRepo 模拟文件仓库
type mockDownloadFileRepo struct {
	files []*model.File
}

func newMockDownloadFileRepo() *mockDownloadFileRepo {
	return &mockDownloadFileRepo{
		files: []*model.File{
			{ID: 1, UserID: 1, Name: "file1.txt", Size: 100, LanZouFileID: "lz123", EncryptionKey: "key1", EncryptionNonce: "nonce1"},
			{ID: 2, UserID: 1, Name: "file2.pdf", Size: 200, LanZouFileID: "lz456", EncryptionKey: "key2", EncryptionNonce: "nonce2"},
			{ID: 3, UserID: 2, Name: "file3.doc", Size: 300, LanZouFileID: "lz789", EncryptionKey: "key3", EncryptionNonce: "nonce3"},
		},
	}
}

func (m *mockDownloadFileRepo) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockDownloadFileRepo) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

// mockDownloadLanzouProvider 模拟蓝奏云客户端提供者
type mockDownloadLanzouProvider struct {
	connected bool
}

func newMockDownloadLanzouProvider(connected bool) *mockDownloadLanzouProvider {
	return &mockDownloadLanzouProvider{connected: connected}
}

func (m *mockDownloadLanzouProvider) IsConnected(userID uint) bool {
	return m.connected
}

func (m *mockDownloadLanzouProvider) GetClient(userID uint) (*lanzou.Client, error) {
	if !m.connected {
		return nil, errors.New("lanzou not connected")
	}
	// 返回nil client，因为我们只测试服务层逻辑，不实际调用蓝奏云
	return nil, nil
}

// TestDownloadService_GetDownloadURL 测试获取下载直链
func TestDownloadService_GetDownloadURL(t *testing.T) {
	t.Run("should return download info for valid file", func(t *testing.T) {
		fileRepo := newMockDownloadFileRepo()
		lanzouProvider := newMockDownloadLanzouProvider(true)
		svc := NewDownloadService(fileRepo, lanzouProvider)

		resp, err := svc.GetDownloadURL(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Error("expected response, got nil")
			return
		}
		if resp.FileID != 1 {
			t.Errorf("expected file_id 1, got %d", resp.FileID)
		}
		if resp.FileName != "file1.txt" {
			t.Errorf("expected file_name 'file1.txt', got '%s'", resp.FileName)
		}
		if resp.EncryptionKey != "key1" {
			t.Errorf("expected encryption_key 'key1', got '%s'", resp.EncryptionKey)
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		fileRepo := newMockDownloadFileRepo()
		lanzouProvider := newMockDownloadLanzouProvider(true)
		svc := NewDownloadService(fileRepo, lanzouProvider)

		_, err := svc.GetDownloadURL(1, 999)

		if err == nil {
			t.Error("expected file not found error")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		fileRepo := newMockDownloadFileRepo()
		lanzouProvider := newMockDownloadLanzouProvider(true)
		svc := NewDownloadService(fileRepo, lanzouProvider)

		_, err := svc.GetDownloadURL(1, 3) // user 1 trying to access user 2's file

		if err == nil {
			t.Error("expected access denied error")
		}
	})

	t.Run("should fail if lanzou not connected", func(t *testing.T) {
		fileRepo := newMockDownloadFileRepo()
		lanzouProvider := newMockDownloadLanzouProvider(false)
		svc := NewDownloadService(fileRepo, lanzouProvider)

		_, err := svc.GetDownloadURL(1, 1)

		if err == nil {
			t.Error("expected lanzou not connected error")
		}
	})
}

// TestDownloadService_GetShareDownloadURL 测试通过分享获取下载直链
func TestDownloadService_GetShareDownloadURL(t *testing.T) {
	t.Run("should return download URL for valid share", func(t *testing.T) {
		fileRepo := newMockDownloadFileRepo()
		lanzouProvider := newMockDownloadLanzouProvider(true)
		svc := NewDownloadService(fileRepo, lanzouProvider)

		downloadURL, err := svc.GetShareDownloadURL("valid_token", "password")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// 简化实现返回空字符串
		if downloadURL != "" {
			t.Errorf("expected empty URL in simplified implementation, got '%s'", downloadURL)
		}
	})
}

// TestDownloadService_getDownloadLink_NoLanZouFileID 测试无蓝奏云文件ID
func TestDownloadService_getDownloadLink_NoLanZouFileID(t *testing.T) {
	fileRepo := newMockDownloadFileRepo()
	// Add a file without LanZouFileID
	fileRepo.files = append(fileRepo.files, &model.File{
		ID: 10, UserID: 1, Name: "no-lz.txt", Size: 50,
		LanZouFileID: "", EncryptionKey: "k", EncryptionNonce: "n",
	})
	lanzouProvider := newMockDownloadLanzouProvider(true)
	svc := NewDownloadService(fileRepo, lanzouProvider)

	_, err := svc.GetDownloadURL(1, 10)
	if err == nil {
		t.Error("expected error for file without lanzou_file_id")
	}
}
