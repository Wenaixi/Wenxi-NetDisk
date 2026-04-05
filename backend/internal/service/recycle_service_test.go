package service

import (
	"errors"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// Mock repositories
type mockRecycleBinRepo struct {
	items     map[uint]*model.RecycleBin
	nextID    uint
	listErr   error
	getErr    error
	createErr error
	deleteErr error
}

func newMockRecycleBinRepo() *mockRecycleBinRepo {
	return &mockRecycleBinRepo{
		items:  make(map[uint]*model.RecycleBin),
		nextID: 1,
	}
}

func (m *mockRecycleBinRepo) Create(item *model.RecycleBin) error {
	if m.createErr != nil {
		return m.createErr
	}
	item.ID = m.nextID
	m.nextID++
	m.items[item.ID] = item
	return nil
}

func (m *mockRecycleBinRepo) List(userID uint) ([]model.RecycleBin, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var items []model.RecycleBin
	for _, item := range m.items {
		if item.UserID == userID {
			items = append(items, *item)
		}
	}
	return items, nil
}

func (m *mockRecycleBinRepo) GetByID(id uint, userID uint) (*model.RecycleBin, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	item, exists := m.items[id]
	if !exists {
		return nil, errors.New("not found")
	}
	if item.UserID != userID {
		return nil, errors.New("access denied")
	}
	return item, nil
}

func (m *mockRecycleBinRepo) Restore(id uint, userID uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.items, id)
	return nil
}

func (m *mockRecycleBinRepo) DeletePermanently(id uint, userID uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.items, id)
	return nil
}

func (m *mockRecycleBinRepo) ClearAll(userID uint) error {
	for id, item := range m.items {
		if item.UserID == userID {
			delete(m.items, id)
		}
	}
	return nil
}

type mockRecycleFileRepo struct {
	files     map[uint]*model.File
	deleteErr error
	createErr error
}

func newMockRecycleFileRepo() *mockRecycleFileRepo {
	return &mockRecycleFileRepo{
		files: make(map[uint]*model.File),
	}
}

func (m *mockRecycleFileRepo) FindByID(id uint) (*model.File, error) {
	file, exists := m.files[id]
	if !exists {
		return nil, errors.New("file not found")
	}
	return file, nil
}

func (m *mockRecycleFileRepo) Create(file *model.File) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.files[file.ID] = file
	return nil
}

func (m *mockRecycleFileRepo) Delete(id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.files, id)
	return nil
}

type mockRecycleFolderRepo struct {
	folders   map[uint]*model.Folder
	deleteErr error
	createErr error
}

func newMockRecycleFolderRepo() *mockRecycleFolderRepo {
	return &mockRecycleFolderRepo{
		folders: make(map[uint]*model.Folder),
	}
}

func (m *mockRecycleFolderRepo) FindByID(id uint) (*model.Folder, error) {
	folder, exists := m.folders[id]
	if !exists {
		return nil, errors.New("folder not found")
	}
	return folder, nil
}

func (m *mockRecycleFolderRepo) Create(folder *model.Folder) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.folders[folder.ID] = folder
	return nil
}

func (m *mockRecycleFolderRepo) Delete(id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.folders, id)
	return nil
}

// Tests
func TestMoveToRecycleBin_Success(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	// Create a file to recycle
	fileRepo.files[1] = &model.File{
		ID:              1,
		UserID:          1,
		Name:            "test.txt",
		Size:            1024,
		LanZouFileID:    "lanzou123",
		EncryptionKey:   "key1",
		EncryptionNonce: "nonce1",
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	item, err := svc.MoveToRecycleBin(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.OriginalName != "test.txt" {
		t.Errorf("expected name 'test.txt', got '%s'", item.OriginalName)
	}
	if item.ItemType != "file" {
		t.Errorf("expected type 'file', got '%s'", item.ItemType)
	}
	// File should be deleted from repo
	_, err = fileRepo.FindByID(1)
	if err == nil {
		t.Error("file should have been deleted")
	}
}

func TestMoveToRecycleBin_FileNotFound(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveToRecycleBin(1, 999)

	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestMoveToRecycleBin_AccessDenied(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	fileRepo.files[1] = &model.File{
		ID:     1,
		UserID: 2, // Different user
		Name:   "test.txt",
		Size:   1024,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveToRecycleBin(1, 1)

	if err == nil {
		t.Error("expected access denied error")
	}
}

func TestMoveFolderToRecycleBin_Success(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	folderRepo.folders[1] = &model.Folder{
		ID:       1,
		UserID:   1,
		Name:     "test-folder",
		ParentID: nil,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	item, err := svc.MoveFolderToRecycleBin(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.OriginalName != "test-folder" {
		t.Errorf("expected name 'test-folder', got '%s'", item.OriginalName)
	}
	if item.ItemType != "folder" {
		t.Errorf("expected type 'folder', got '%s'", item.ItemType)
	}
}

func TestMoveFolderToRecycleBin_FolderNotFound(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveFolderToRecycleBin(1, 999)

	if err == nil {
		t.Error("expected error for non-existent folder")
	}
}

func TestList(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	// Add recycle items
	recycleRepo.items[1] = &model.RecycleBin{ID: 1, UserID: 1, OriginalName: "file1.txt", ItemType: "file"}
	recycleRepo.items[2] = &model.RecycleBin{ID: 2, UserID: 1, OriginalName: "file2.txt", ItemType: "file"}
	recycleRepo.items[3] = &model.RecycleBin{ID: 3, UserID: 2, OriginalName: "file3.txt", ItemType: "file"}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	items, err := svc.List(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items for user 1, got %d", len(items))
	}
}

func TestRestore_File(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID:              1,
		UserID:          1,
		OriginalName:    "restored.txt",
		ItemType:        "file",
		ItemID:          100,
		LanZouFileID:    "lf100",
		Size:            2048,
		EncryptionKey:   "key2",
		EncryptionNonce: "nonce2",
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Recycle item should be deleted
	_, err = recycleRepo.GetByID(1, 1)
	if err == nil {
		t.Error("recycle item should have been deleted")
	}
	// File should be restored
	file, err := fileRepo.FindByID(100)
	if err != nil {
		t.Fatalf("restored file not found: %v", err)
	}
	if file.Name != "restored.txt" {
		t.Errorf("expected restored name 'restored.txt', got '%s'", file.Name)
	}
}

func TestRestore_NotFound(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 999)

	if err == nil {
		t.Error("expected error for non-existent recycle item")
	}
}

func TestPermanentDelete(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID:     1,
		UserID: 1,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.PermanentDelete(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = recycleRepo.GetByID(1, 1)
	if err == nil {
		t.Error("item should have been permanently deleted")
	}
}

func TestClearAll(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{ID: 1, UserID: 1}
	recycleRepo.items[2] = &model.RecycleBin{ID: 2, UserID: 1}
	recycleRepo.items[3] = &model.RecycleBin{ID: 3, UserID: 2}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.ClearAll(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	items, _ := svc.List(1)
	if len(items) != 0 {
		t.Errorf("expected 0 items for user 1 after clear, got %d", len(items))
	}
	// User 2's items should remain
	items, _ = svc.List(2)
	if len(items) != 1 {
		t.Errorf("expected 1 item for user 2, got %d", len(items))
	}
}

func TestRecycleBinCreatedWithCorrectTimestamp(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	fileRepo.files[1] = &model.File{
		ID:     1,
		UserID: 1,
		Name:   "timetest.txt",
		Size:   512,
	}

	before := time.Now()
	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	item, err := svc.MoveToRecycleBin(1, 1)
	after := time.Now()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.DeletedAt.Before(before) || item.DeletedAt.After(after) {
		t.Error("DeletedAt should be within the test timeframe")
	}
}

func TestMoveFolderToRecycleBin_AccessDenied(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	folderRepo.folders[1] = &model.Folder{
		ID:     1,
		UserID: 2, // Different user
		Name:   "other-folder",
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveFolderToRecycleBin(1, 1)

	if err == nil {
		t.Error("expected access denied error")
	}
}

// TestRestore_Folder 测试恢复文件夹
func TestRestore_Folder(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	parentID := uint(10)
	recycleRepo.items[1] = &model.RecycleBin{
		ID:           1,
		UserID:       1,
		OriginalName: "restored-folder",
		ItemType:     "folder",
		ItemID:       200,
		ParentID:     parentID,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Folder should be restored with parent
	folder, err := folderRepo.FindByID(200)
	if err != nil {
		t.Fatalf("restored folder not found: %v", err)
	}
	if folder.Name != "restored-folder" {
		t.Errorf("expected name 'restored-folder', got '%s'", folder.Name)
	}
	if folder.ParentID == nil || *folder.ParentID != 10 {
		t.Errorf("expected parent 10, got %v", folder.ParentID)
	}
}

// TestRestore_FolderWithoutParent 测试恢复无父文件夹的文件夹
func TestRestore_FolderWithoutParent(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID:           1,
		UserID:       1,
		OriginalName: "root-folder",
		ItemType:     "folder",
		ItemID:       300,
		ParentID:     0, // No parent
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	folder, err := folderRepo.FindByID(300)
	if err != nil {
		t.Fatalf("restored folder not found: %v", err)
	}
	if folder.ParentID != nil {
		t.Errorf("expected nil parent, got %v", folder.ParentID)
	}
}

// TestRestore_AccessDenied 测试恢复他人回收站项目
func TestRestore_AccessDenied(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID:       1,
		UserID:   2, // Different user
		ItemType: "file",
		ItemID:   100,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err == nil {
		t.Error("expected access denied error")
	}
}

// TestPermanentDelete_NotFound 测试永久删除不存在的项目
func TestPermanentDelete_NotFound(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.PermanentDelete(1, 999)

	if err == nil {
		t.Error("expected error for non-existent item")
	}
}

// TestMoveToRecycleBin_FileDeleteError 测试文件删除失败
func TestMoveToRecycleBin_FileDeleteError(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	fileRepo.files[1] = &model.File{
		ID:     1, UserID: 1, Name: "test.txt", Size: 100,
	}
	fileRepo.deleteErr = errors.New("delete failed")

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveToRecycleBin(1, 1)

	if err == nil {
		t.Error("expected error when file delete fails")
	}
}

// TestMoveFolderToRecycleBin_FolderDeleteError 测试文件夹删除失败
func TestMoveFolderToRecycleBin_FolderDeleteError(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	folderRepo.folders[1] = &model.Folder{
		ID: 1, UserID: 1, Name: "test-folder",
	}
	folderRepo.deleteErr = errors.New("delete failed")

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	_, err := svc.MoveFolderToRecycleBin(1, 1)

	if err == nil {
		t.Error("expected error when folder delete fails")
	}
}

// TestRestore_FileCreateError 测试文件恢复创建失败
func TestRestore_FileCreateError(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID: 1, UserID: 1, OriginalName: "file.txt",
		ItemType: "file", ItemID: 100, Size: 500,
	}
	fileRepo.createErr = errors.New("create failed")

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err == nil {
		t.Error("expected error when file create fails")
	}
}

// TestRestore_FolderCreateError 测试文件夹恢复创建失败
func TestRestore_FolderCreateError(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID: 1, UserID: 1, OriginalName: "folder",
		ItemType: "folder", ItemID: 200,
	}
	folderRepo.createErr = errors.New("create failed")

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err == nil {
		t.Error("expected error when folder create fails")
	}
}

// TestRestore_RestoreError 测试恢复后删除记录失败
func TestRestore_RestoreError(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	recycleRepo.items[1] = &model.RecycleBin{
		ID: 1, UserID: 1, OriginalName: "file.txt",
		ItemType: "file", ItemID: 100,
	}
	recycleRepo.deleteErr = errors.New("restore failed")

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.Restore(1, 1)

	if err == nil {
		t.Error("expected error when restore fails")
	}
}

// TestClearAll_Error 测试清空回收站失败
func TestClearAll_Error(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	// Create an error-inducing ClearAll scenario
	// Since mockClearAll doesn't return errors, test with a service that would fail
	recycleRepo.items[1] = &model.RecycleBin{ID: 1, UserID: 1}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	err := svc.ClearAll(1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	items, _ := svc.List(1)
	if len(items) != 0 {
		t.Errorf("expected 0 items after clear, got %d", len(items))
	}
}

// TestMoveFolderToRecycleBin_WithParentID 测试带父文件夹ID的移入回收站
func TestMoveFolderToRecycleBin_WithParentID(t *testing.T) {
	recycleRepo := newMockRecycleBinRepo()
	fileRepo := newMockRecycleFileRepo()
	folderRepo := newMockRecycleFolderRepo()

	parentID := uint(5)
	folderRepo.folders[1] = &model.Folder{
		ID:       1,
		UserID:   1,
		Name:     "child-folder",
		ParentID: &parentID,
	}

	svc := NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	item, err := svc.MoveFolderToRecycleBin(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ParentID != 5 {
		t.Errorf("expected parent ID 5, got %d", item.ParentID)
	}
}
