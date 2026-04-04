package service

import (
	"errors"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// mockFileRepo 模拟文件仓库
type mockFileRepo struct {
	files []*model.File
}

func newMockFileRepo() *mockFileRepo {
	return &mockFileRepo{
		files: []*model.File{
			{
				ID:     1,
				UserID: 1,
				Name:   "document.pdf",
				Size:   1024,
			},
			{
				ID:     2,
				UserID: 1,
				Name:   "image.png",
				Size:   2048,
			},
			{
				ID:     3,
				UserID: 2,
				Name:   "other.txt",
				Size:   512,
			},
		},
	}
}

func (m *mockFileRepo) Create(file *model.File) error {
	file.ID = uint(len(m.files) + 1)
	m.files = append(m.files, file)
	return nil
}

func (m *mockFileRepo) FindByUserID(userID uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}

func (m *mockFileRepo) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockFileRepo) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}

func (m *mockFileRepo) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockFileRepo) Delete(id uint) error {
	for i, f := range m.files {
		if f.ID == id {
			m.files = append(m.files[:i], m.files[i+1:]...)
			return nil
		}
	}
	return errors.New("file not found")
}

func (m *mockFileRepo) Update(file *model.File) error {
	for i, f := range m.files {
		if f.ID == file.ID {
			m.files[i] = file
			return nil
		}
	}
	return errors.New("file not found")
}

// TestFileService_ListFiles tests listing files
func TestFileService_ListFiles(t *testing.T) {
	t.Run("should return files for user", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		files, err := svc.ListFiles(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(files) != 2 {
			t.Errorf("expected 2 files for user 1, got %d", len(files))
		}
	})

	t.Run("should return empty list for user with no files", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		files, err := svc.ListFiles(999)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(files) != 0 {
			t.Errorf("expected 0 files for user 999, got %d", len(files))
		}
	})
}

// TestFileService_GetFile tests getting a single file
func TestFileService_GetFile(t *testing.T) {
	t.Run("should return file for owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.GetFile(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.Name != "document.pdf" {
			t.Errorf("expected name 'document.pdf', got '%s'", file.Name)
		}
	})

	t.Run("should deny access for non-owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.GetFile(999, 1)

		if err == nil {
			t.Error("expected access denied error")
		}
		if file != nil {
			t.Error("expected nil file for non-owner")
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.GetFile(1, 999)

		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if file != nil {
			t.Error("expected nil file for non-existent file")
		}
	})
}

// TestFileService_CreateMetadata tests file metadata creation
func TestFileService_CreateMetadata(t *testing.T) {
	t.Run("should create file metadata", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		req := &CreateFileRequest{
			Name:            "new-file.txt",
			Size:            4096,
			LanZouFileID:    "lanzou-123",
			EncryptionKey:   "encrypted-key",
			EncryptionNonce: "nonce-123",
			MimeType:        "text/plain",
		}

		file, err := svc.CreateMetadata(1, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.Name != "new-file.txt" {
			t.Errorf("expected name 'new-file.txt', got '%s'", file.Name)
		}
		if file.UserID != 1 {
			t.Errorf("expected user ID 1, got %d", file.UserID)
		}
		if file.EncryptionKey != "encrypted-key" {
			t.Errorf("expected encryption key, got '%s'", file.EncryptionKey)
		}
	})
}

// TestFileService_DeleteFile tests file deletion
func TestFileService_DeleteFile(t *testing.T) {
	t.Run("should delete file for owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		err := svc.DeleteFile(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify file is deleted
		_, err = svc.GetFile(1, 1)
		if err == nil {
			t.Error("expected error after deletion")
		}
	})

	t.Run("should deny delete for non-owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		err := svc.DeleteFile(999, 1)

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFileService_UpdateFileDescription tests updating file description
func TestFileService_UpdateFileDescription(t *testing.T) {
	t.Run("should update description for owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileDescription(1, 1, "test description")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.Description != "test description" {
			t.Errorf("expected description 'test description', got '%s'", file.Description)
		}
	})

	t.Run("should deny update for non-owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileDescription(999, 1, "new desc")

		if err == nil {
			t.Error("expected access denied error")
		}
		if file != nil {
			t.Error("expected nil file for non-owner")
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileDescription(1, 999, "new desc")

		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if file != nil {
			t.Error("expected nil file for non-existent file")
		}
	})

	t.Run("should clear description with empty string", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		// First set a description
		svc.UpdateFileDescription(1, 1, "initial description")

		// Then clear it
		file, err := svc.UpdateFileDescription(1, 1, "")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file.Description != "" {
			t.Errorf("expected empty description, got '%s'", file.Description)
		}
	})
}

// TestFileService_UpdateFileName tests updating file name
func TestFileService_UpdateFileName(t *testing.T) {
	t.Run("should update name for owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileName(1, 1, "renamed.pdf")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.Name != "renamed.pdf" {
			t.Errorf("expected name 'renamed.pdf', got '%s'", file.Name)
		}
	})

	t.Run("should deny update for non-owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileName(999, 1, "renamed.pdf")

		if err == nil {
			t.Error("expected access denied error")
		}
		if file != nil {
			t.Error("expected nil file for non-owner")
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.UpdateFileName(1, 999, "renamed.pdf")

		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if file != nil {
			t.Error("expected nil file for non-existent file")
		}
	})
}

// TestFileService_MoveFile tests moving files to folders
func TestFileService_MoveFile(t *testing.T) {
	t.Run("should move file to new folder for owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		folderID := uint(5)
		file, err := svc.MoveFile(1, 1, &folderID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.FolderID == nil || *file.FolderID != 5 {
			t.Errorf("expected folder ID 5, got %v", file.FolderID)
		}
	})

	t.Run("should move file to root (nil folder)", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		file, err := svc.MoveFile(1, 1, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
		}
	})

	t.Run("should deny move for non-owner", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		folderID := uint(5)
		file, err := svc.MoveFile(999, 1, &folderID)

		if err == nil {
			t.Error("expected access denied error")
		}
		if file != nil {
			t.Error("expected nil file for non-owner")
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		repo := newMockFileRepo()
		svc := NewFileService(repo)

		folderID := uint(5)
		file, err := svc.MoveFile(1, 999, &folderID)

		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if file != nil {
			t.Error("expected nil file for non-existent file")
		}
	})
}
