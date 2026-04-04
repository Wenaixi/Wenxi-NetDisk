package service

import (
	"errors"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// mockFolderRepo 模拟文件夹仓库
type mockFolderRepo struct {
	folders []*model.Folder
}

func newMockFolderRepo() *mockFolderRepo {
	rootID := uint(1)
	return &mockFolderRepo{
		folders: []*model.Folder{
			{ID: 1, UserID: 1, ParentID: nil, Name: "Root Folder"},
			{ID: 2, UserID: 1, ParentID: &rootID, Name: "Sub Folder 1"},
			{ID: 3, UserID: 1, ParentID: &rootID, Name: "Sub Folder 2"},
			{ID: 4, UserID: 2, ParentID: nil, Name: "Other User Folder"},
		},
	}
}

func (m *mockFolderRepo) Create(folder *model.Folder) error {
	folder.ID = uint(len(m.folders) + 1)
	m.folders = append(m.folders, folder)
	return nil
}

func (m *mockFolderRepo) FindByID(id uint) (*model.Folder, error) {
	for _, f := range m.folders {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("folder not found")
}

func (m *mockFolderRepo) FindByUserID(userID uint) ([]model.Folder, error) {
	var result []model.Folder
	for _, f := range m.folders {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}

func (m *mockFolderRepo) FindByParentID(userID uint, parentID *uint) ([]model.Folder, error) {
	var result []model.Folder
	for _, f := range m.folders {
		if f.UserID != userID {
			continue
		}
		if parentID == nil {
			if f.ParentID == nil {
				result = append(result, *f)
			}
		} else {
			if f.ParentID != nil && *f.ParentID == *parentID {
				result = append(result, *f)
			}
		}
	}
	return result, nil
}

func (m *mockFolderRepo) Update(folder *model.Folder) error {
	for i, f := range m.folders {
		if f.ID == folder.ID {
			m.folders[i] = folder
			return nil
		}
	}
	return errors.New("folder not found")
}

func (m *mockFolderRepo) Delete(id uint) error {
	for i, f := range m.folders {
		if f.ID == id {
			m.folders = append(m.folders[:i], m.folders[i+1:]...)
			return nil
		}
	}
	return errors.New("folder not found")
}

func (m *mockFolderRepo) DeleteByUserID(userID uint) error {
	var remaining []*model.Folder
	for _, f := range m.folders {
		if f.UserID != userID {
			remaining = append(remaining, f)
		}
	}
	m.folders = remaining
	return nil
}

func (m *mockFolderRepo) DeleteByParentID(parentID uint) error {
	var remaining []*model.Folder
	for _, f := range m.folders {
		if f.ParentID == nil || *f.ParentID != parentID {
			remaining = append(remaining, f)
		}
	}
	m.folders = remaining
	return nil
}

// TestFolderService_CreateFolder 测试创建文件夹
func TestFolderService_CreateFolder(t *testing.T) {
	t.Run("should create folder at root", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		req := &CreateFolderRequest{Name: "New Folder"}
		folder, err := svc.CreateFolder(1, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder == nil {
			t.Error("expected folder, got nil")
			return
		}
		if folder.Name != "New Folder" {
			t.Errorf("expected name 'New Folder', got '%s'", folder.Name)
		}
		if folder.ParentID != nil {
			t.Error("expected nil parent for root folder")
		}
	})

	t.Run("should create folder with parent", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		parentID := uint(1)
		req := &CreateFolderRequest{Name: "Child Folder", ParentID: &parentID}
		folder, err := svc.CreateFolder(1, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder == nil {
			t.Error("expected folder, got nil")
			return
		}
		if folder.ParentID == nil || *folder.ParentID != 1 {
			t.Error("expected parent_id 1")
		}
	})

	t.Run("should fail if parent not found", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		nonExistentID := uint(999)
		req := &CreateFolderRequest{Name: "Test", ParentID: &nonExistentID}
		_, err := svc.CreateFolder(1, req)

		if err == nil {
			t.Error("expected error for non-existent parent")
		}
	})

	t.Run("should fail if parent belongs to another user", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		// Folder 4 belongs to user 2
		parentID := uint(4)
		req := &CreateFolderRequest{Name: "Test", ParentID: &parentID}
		_, err := svc.CreateFolder(1, req)

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFolderService_ListFolders 测试列出文件夹
func TestFolderService_ListFolders(t *testing.T) {
	t.Run("should list all folders for user", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folders, err := svc.ListFolders(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(folders) != 3 {
			t.Errorf("expected 3 folders, got %d", len(folders))
		}
	})

	t.Run("should list folders by parent", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		parentID := uint(1)
		folders, err := svc.ListFoldersByParent(1, &parentID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(folders) != 2 {
			t.Errorf("expected 2 sub folders, got %d", len(folders))
		}
	})

	t.Run("should list root folders", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folders, err := svc.ListFoldersByParent(1, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(folders) != 1 {
			t.Errorf("expected 1 root folder, got %d", len(folders))
		}
	})
}

// TestFolderService_GetFolder 测试获取文件夹
func TestFolderService_GetFolder(t *testing.T) {
	t.Run("should get existing folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.GetFolder(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder == nil {
			t.Error("expected folder, got nil")
			return
		}
		if folder.Name != "Root Folder" {
			t.Errorf("expected 'Root Folder', got '%s'", folder.Name)
		}
	})

	t.Run("should return error for non-existent folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		_, err := svc.GetFolder(1, 999)

		if err == nil {
			t.Error("expected error for non-existent folder")
		}
	})

	t.Run("should return error for folder of another user", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		_, err := svc.GetFolder(1, 4)

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFolderService_UpdateFolder 测试更新文件夹
func TestFolderService_UpdateFolder(t *testing.T) {
	t.Run("should update folder name", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.UpdateFolder(1, 1, "Renamed Folder")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder.Name != "Renamed Folder" {
			t.Errorf("expected 'Renamed Folder', got '%s'", folder.Name)
		}
	})

	t.Run("should fail for non-owned folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		_, err := svc.UpdateFolder(1, 4, "Hack")

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFolderService_DeleteFolder 测试删除文件夹
func TestFolderService_DeleteFolder(t *testing.T) {
	t.Run("should delete folder and its children", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		err := svc.DeleteFolder(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify children are also deleted
		children, _ := repo.FindByParentID(1, nil)
		for _, child := range children {
			if child.ParentID != nil && *child.ParentID == 1 {
				t.Error("child folder should have been deleted")
			}
		}
	})

	t.Run("should fail for non-owned folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		err := svc.DeleteFolder(1, 4)

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFolderService_MoveFolder 测试移动文件夹
func TestFolderService_MoveFolder(t *testing.T) {
	t.Run("should move folder to new parent", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		newParentID := uint(3)
		folder, err := svc.MoveFolder(1, 2, &newParentID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder.ParentID == nil || *folder.ParentID != 3 {
			t.Error("expected parent_id 3")
		}
	})

	t.Run("should move folder to root", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.MoveFolder(1, 2, nil)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder.ParentID != nil {
			t.Error("expected nil parent for root")
		}
	})

	t.Run("should fail when moving to itself", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		parentID := uint(1)
		_, err := svc.MoveFolder(1, 1, &parentID)

		if err == nil {
			t.Error("expected error when moving folder to itself")
		}
	})

	t.Run("should fail when moving to descendant", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		// Folder 2 is child of 1, so moving 1 to 2 should fail
		parentID := uint(2)
		_, err := svc.MoveFolder(1, 1, &parentID)

		if err == nil {
			t.Error("expected error when moving folder to its descendant")
		}
	})

	t.Run("should fail for non-owned folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		newParentID := uint(1)
		_, err := svc.MoveFolder(1, 4, &newParentID)

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestFolderService_UpdateFolderDescription 测试更新文件夹描述
func TestFolderService_UpdateFolderDescription(t *testing.T) {
	t.Run("should update description for owner", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.UpdateFolderDescription(1, 1, "test folder description")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder == nil {
			t.Error("expected folder, got nil")
			return
		}
		if folder.Description != "test folder description" {
			t.Errorf("expected description 'test folder description', got '%s'", folder.Description)
		}
	})

	t.Run("should deny update for non-owner", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.UpdateFolderDescription(999, 1, "new desc")

		if err == nil {
			t.Error("expected access denied error")
		}
		if folder != nil {
			t.Error("expected nil folder for non-owner")
		}
	})

	t.Run("should return error for non-existent folder", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		folder, err := svc.UpdateFolderDescription(1, 999, "new desc")

		if err == nil {
			t.Error("expected error for non-existent folder")
		}
		if folder != nil {
			t.Error("expected nil folder for non-existent folder")
		}
	})

	t.Run("should clear description with empty string", func(t *testing.T) {
		repo := newMockFolderRepo()
		svc := NewFolderService(repo)

		// First set a description
		svc.UpdateFolderDescription(1, 1, "initial description")

		// Then clear it
		folder, err := svc.UpdateFolderDescription(1, 1, "")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if folder.Description != "" {
			t.Errorf("expected empty description, got '%s'", folder.Description)
		}
	})
}
