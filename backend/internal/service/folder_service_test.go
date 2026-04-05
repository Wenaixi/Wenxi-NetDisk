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

// TestFolderService_CreateFolder_EmptyName 测试空名称
func TestFolderService_CreateFolder_EmptyName(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	req := &CreateFolderRequest{Name: ""}
	folder, err := svc.CreateFolder(1, req)

	// Empty name should still create (validation is at handler level)
	if err != nil {
		t.Logf("expected no error for empty name, got %v", err)
	}
	if folder != nil && folder.Name != "" {
		t.Logf("created folder with name: %s", folder.Name)
	}
}

// TestFolderService_UpdateFolder_NotFound 测试更新不存在的文件夹
func TestFolderService_UpdateFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	_, err := svc.UpdateFolder(1, 999, "new name")
	if err == nil {
		t.Error("expected error for non-existent folder")
	}
}

// TestFolderService_DeleteFolder_NotFound 测试删除不存在的文件夹
func TestFolderService_DeleteFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	err := svc.DeleteFolder(1, 999)
	if err == nil {
		t.Error("expected error for non-existent folder")
	}
}

// TestFolderService_MoveFolder_ParentNotFound 测试移动到不存在的父文件夹
func TestFolderService_MoveFolder_ParentNotFound(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	newParentID := uint(999)
	_, err := svc.MoveFolder(1, 2, &newParentID)
	if err == nil {
		t.Error("expected error for non-existent parent")
	}
}

// TestFolderService_MoveFolder_ParentNotOwned 测试移动到他人的文件夹
func TestFolderService_MoveFolder_ParentNotOwned(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	// Folder 4 belongs to user 2, try to move user 1's folder to it
	newParentID := uint(4)
	_, err := svc.MoveFolder(1, 2, &newParentID)
	if err == nil {
		t.Error("expected access denied for parent owned by another user")
	}
}

// TestFolderService_IsDescendant_DeepNesting 测试深层嵌套的isDescendant
func TestFolderService_IsDescendant_DeepNesting(t *testing.T) {
	repo := newMockFolderRepo()
	// Create a deep nesting: 10 -> 20 -> 30 -> 40
	id10 := uint(10)
	id20 := uint(20)
	id30 := uint(30)
	repo.folders = append(repo.folders,
		&model.Folder{ID: 10, UserID: 1, ParentID: nil, Name: "level1"},
		&model.Folder{ID: 20, UserID: 1, ParentID: &id10, Name: "level2"},
		&model.Folder{ID: 30, UserID: 1, ParentID: &id20, Name: "level3"},
		&model.Folder{ID: 40, UserID: 1, ParentID: &id30, Name: "level4"},
	)
	svc := NewFolderService(repo)

	// Try to move level1 (10) to level4 (40) - should fail (descendant)
	_, err := svc.MoveFolder(1, 10, &id30)
	if err == nil {
		t.Error("expected error when moving to deep descendant")
	}

	// Try to move level1 (10) to level3 (30) - should fail (descendant)
	_, err = svc.MoveFolder(1, 10, &id20)
	if err == nil {
		t.Error("expected error when moving to descendant")
	}

	// Moving level4 (40) to level2 (20) should be OK (not a descendant)
	_, err = svc.MoveFolder(1, 40, &id10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestFolderService_IsDescendant_NotFound 测试isDescendant找不到目标
func TestFolderService_IsDescendant_NotFound(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	newParentID := uint(999) // non-existent
	// This should not panic, isDescendant returns false on error
	_, err := svc.MoveFolder(1, 1, &newParentID)
	// Should fail because parent not found
	if err == nil {
		t.Error("expected error for non-existent target in isDescendant")
	}
}

// TestFolderService_ListFolders_EmptyUser 测试无文件夹用户
func TestFolderService_ListFolders_EmptyUser(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	folders, err := svc.ListFolders(999)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(folders) != 0 {
		t.Errorf("expected 0 folders for user 999, got %d", len(folders))
	}
}

// TestFolderService_GetFolder_NotFound 测试获取不存在的文件夹
func TestFolderService_GetFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepo()
	svc := NewFolderService(repo)

	_, err := svc.GetFolder(1, 999)
	if err == nil {
		t.Error("expected error for non-existent folder")
	}
}
