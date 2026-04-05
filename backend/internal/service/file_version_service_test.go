package service

import (
	"errors"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// mockFileVersionRepo 模拟文件版本仓库
type mockFileVersionRepo struct {
	versions []*model.FileVersion
}

func newMockFileVersionRepo() *mockFileVersionRepo {
	return &mockFileVersionRepo{
		versions: []*model.FileVersion{
			{
				ID:              1,
				FileID:          1,
				UserID:          1,
				LanZouFileID:    "lanzou_v1",
				Size:            100,
				EncryptionKey:   "key_v1",
				EncryptionNonce: "nonce_v1",
				Description:     "version 1",
			},
			{
				ID:              2,
				FileID:          1,
				UserID:          1,
				LanZouFileID:    "lanzou_v2",
				Size:            200,
				EncryptionKey:   "key_v2",
				EncryptionNonce: "nonce_v2",
				Description:     "version 2",
			},
			{
				ID:              3,
				FileID:          2,
				UserID:          1,
				LanZouFileID:    "lanzou_v3",
				Size:            300,
				EncryptionKey:   "key_v3",
				EncryptionNonce: "nonce_v3",
				Description:     "version 3",
			},
			{
				ID:              4,
				FileID:          3,
				UserID:          2,
				LanZouFileID:    "lanzou_v4",
				Size:            400,
				EncryptionKey:   "key_v4",
				EncryptionNonce: "nonce_v4",
				Description:     "version 4",
			},
		},
	}
}

func (m *mockFileVersionRepo) Create(version *model.FileVersion) error {
	version.ID = uint(len(m.versions) + 1)
	m.versions = append(m.versions, version)
	return nil
}

func (m *mockFileVersionRepo) FindByFileID(fileID uint) ([]model.FileVersion, error) {
	var result []model.FileVersion
	for _, v := range m.versions {
		if v.FileID == fileID {
			result = append(result, *v)
		}
	}
	return result, nil
}

func (m *mockFileVersionRepo) FindByID(id uint) (*model.FileVersion, error) {
	for _, v := range m.versions {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("version not found")
}

func (m *mockFileVersionRepo) Delete(id uint) error {
	for i, v := range m.versions {
		if v.ID == id {
			m.versions = append(m.versions[:i], m.versions[i+1:]...)
			return nil
		}
	}
	return errors.New("version not found")
}

func (m *mockFileVersionRepo) DeleteByFileID(fileID uint) error {
	var remaining []*model.FileVersion
	for _, v := range m.versions {
		if v.FileID != fileID {
			remaining = append(remaining, v)
		}
	}
	m.versions = remaining
	return nil
}

// mockVersionFileRepo 模拟版本服务所需的文件仓库
type mockVersionFileRepo struct {
	files []*model.File
}

func newMockVersionFileRepo() *mockVersionFileRepo {
	return &mockVersionFileRepo{
		files: []*model.File{
			{ID: 1, UserID: 1, Name: "file1.txt", Size: 100, LanZouFileID: "lanzou_1", EncryptionKey: "key_1", EncryptionNonce: "nonce_1"},
			{ID: 2, UserID: 1, Name: "file2.pdf", Size: 200, LanZouFileID: "lanzou_2", EncryptionKey: "key_2", EncryptionNonce: "nonce_2"},
			{ID: 3, UserID: 2, Name: "file3.doc", Size: 300, LanZouFileID: "lanzou_3", EncryptionKey: "key_3", EncryptionNonce: "nonce_3"},
		},
	}
}

func (m *mockVersionFileRepo) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockVersionFileRepo) Update(file *model.File) error {
	for i, f := range m.files {
		if f.ID == file.ID {
			m.files[i] = file
			return nil
		}
	}
	return errors.New("file not found")
}

// TestFileVersionService_CreateVersion 测试创建版本
func TestFileVersionService_CreateVersion(t *testing.T) {
	t.Run("should create version successfully", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		req := &CreateFileVersionRequest{
			LanZouFileID:    "new_lanzou_id",
			Size:            500,
			EncryptionKey:   "new_key",
			EncryptionNonce: "new_nonce",
			Description:     "new version",
		}

		version, err := svc.CreateVersion(1, 1, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if version == nil {
			t.Error("expected version, got nil")
			return
		}
		if version.FileID != 1 {
			t.Errorf("expected file_id 1, got %d", version.FileID)
		}
		if version.UserID != 1 {
			t.Errorf("expected user_id 1, got %d", version.UserID)
		}
		if version.LanZouFileID != "new_lanzou_id" {
			t.Errorf("expected lanzou_file_id 'new_lanzou_id', got '%s'", version.LanZouFileID)
		}
		if version.Size != 500 {
			t.Errorf("expected size 500, got %d", version.Size)
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		req := &CreateFileVersionRequest{
			LanZouFileID:    "new_lanzou_id",
			Size:            500,
			EncryptionKey:   "new_key",
			EncryptionNonce: "new_nonce",
		}

		version, err := svc.CreateVersion(1, 999, req)

		if err == nil {
			t.Error("expected file not found error")
		}
		if version != nil {
			t.Error("expected nil version")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		req := &CreateFileVersionRequest{
			LanZouFileID:    "new_lanzou_id",
			Size:            500,
			EncryptionKey:   "new_key",
			EncryptionNonce: "new_nonce",
		}

		version, err := svc.CreateVersion(3, 1, req) // user 3 doesn't own file 1

		if err == nil {
			t.Error("expected access denied error")
		}
		if version != nil {
			t.Error("expected nil version")
		}
	})
}

// TestFileVersionService_ListVersions 测试获取版本列表
func TestFileVersionService_ListVersions(t *testing.T) {
	t.Run("should list all versions for file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		versions, err := svc.ListVersions(1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(versions) != 2 {
			t.Errorf("expected 2 versions for file 1, got %d", len(versions))
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		versions, err := svc.ListVersions(1, 999)

		if err == nil {
			t.Error("expected file not found error")
		}
		if versions != nil {
			t.Error("expected nil versions")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		versions, err := svc.ListVersions(3, 1)

		if err == nil {
			t.Error("expected access denied error")
		}
		if versions != nil {
			t.Error("expected nil versions")
		}
	})

	t.Run("should return empty for file with no versions", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)
		// Add a new file to the repo that has no versions
		fileRepo.files = append(fileRepo.files, &model.File{ID: 99, UserID: 1, Name: "noversion.txt"})

		versions, err := svc.ListVersions(1, 99)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(versions) != 0 {
			t.Errorf("expected 0 versions, got %d", len(versions))
		}
	})
}

// TestFileVersionService_RestoreVersion 测试恢复版本
func TestFileVersionService_RestoreVersion(t *testing.T) {
	t.Run("should restore version successfully", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		file, err := svc.RestoreVersion(1, 1, 2) // restore version 2 of file 1

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		// Verify file was updated with version 2 data
		if file.LanZouFileID != "lanzou_v2" {
			t.Errorf("expected lanzou_file_id 'lanzou_v2', got '%s'", file.LanZouFileID)
		}
		if file.Size != 200 {
			t.Errorf("expected size 200, got %d", file.Size)
		}
		if file.EncryptionKey != "key_v2" {
			t.Errorf("expected encryption_key 'key_v2', got '%s'", file.EncryptionKey)
		}
		if file.EncryptionNonce != "nonce_v2" {
			t.Errorf("expected encryption_nonce 'nonce_v2', got '%s'", file.EncryptionNonce)
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		file, err := svc.RestoreVersion(1, 999, 1)

		if err == nil {
			t.Error("expected file not found error")
		}
		if file != nil {
			t.Error("expected nil file")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		file, err := svc.RestoreVersion(3, 1, 1)

		if err == nil {
			t.Error("expected access denied error")
		}
		if file != nil {
			t.Error("expected nil file")
		}
	})

	t.Run("should fail if version not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		file, err := svc.RestoreVersion(1, 1, 999)

		if err == nil {
			t.Error("expected version not found error")
		}
		if file != nil {
			t.Error("expected nil file")
		}
	})

	t.Run("should fail if version doesn't belong to file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		file, err := svc.RestoreVersion(1, 1, 3) // version 3 belongs to file 2

		if err == nil {
			t.Error("expected version mismatch error")
		}
		if file != nil {
			t.Error("expected nil file")
		}
	})
}

// TestFileVersionService_DeleteVersion 测试删除版本
func TestFileVersionService_DeleteVersion(t *testing.T) {
	t.Run("should delete version successfully", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		err := svc.DeleteVersion(1, 1, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify version was deleted
		versions, _ := svc.ListVersions(1, 1)
		if len(versions) != 1 {
			t.Errorf("expected 1 version after delete, got %d", len(versions))
		}
	})

	t.Run("should fail if file not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		err := svc.DeleteVersion(1, 999, 1)

		if err == nil {
			t.Error("expected file not found error")
		}
	})

	t.Run("should fail if user doesn't own file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		err := svc.DeleteVersion(3, 1, 1)

		if err == nil {
			t.Error("expected access denied error")
		}
	})

	t.Run("should fail if version not found", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		err := svc.DeleteVersion(1, 1, 999)

		if err == nil {
			t.Error("expected version not found error")
		}
	})

	t.Run("should fail if version doesn't belong to file", func(t *testing.T) {
		versionRepo := newMockFileVersionRepo()
		fileRepo := newMockVersionFileRepo()
		svc := NewFileVersionService(versionRepo, fileRepo)

		err := svc.DeleteVersion(1, 1, 3) // version 3 belongs to file 2

		if err == nil {
			t.Error("expected version mismatch error")
		}
	})
}

// TestFileVersionService_CreateVersion_RepoError 测试创建版本仓库错误
func TestFileVersionService_CreateVersion_RepoError(t *testing.T) {
	repo := &mockFileVersionRepoWithError{createErr: errors.New("db write failed")}
	fileRepo := newMockVersionFileRepo()
	svc := NewFileVersionService(repo, fileRepo)

	req := &CreateFileVersionRequest{
		LanZouFileID: "id", Size: 100, EncryptionKey: "k", EncryptionNonce: "n",
	}
	_, err := svc.CreateVersion(1, 1, req)
	if err == nil {
		t.Error("expected db write error")
	}
}

// TestFileVersionService_RestoreVersion_UpdateError 测试恢复版本时更新文件失败
func TestFileVersionService_RestoreVersion_UpdateError(t *testing.T) {
	versionRepo := newMockFileVersionRepo()
	fileRepo := &mockVersionFileRepoWithError{updateErr: errors.New("db update failed")}
	svc := NewFileVersionService(versionRepo, fileRepo)

	_, err := svc.RestoreVersion(1, 1, 1)
	if err == nil {
		t.Error("expected db update error")
	}
}

// mockFileVersionRepoWithError 模拟Create返回错误的版本仓库
type mockFileVersionRepoWithError struct {
	createErr error
}

func (m *mockFileVersionRepoWithError) Create(version *model.FileVersion) error {
	return m.createErr
}
func (m *mockFileVersionRepoWithError) FindByFileID(fileID uint) ([]model.FileVersion, error) {
	return nil, nil
}
func (m *mockFileVersionRepoWithError) FindByID(id uint) (*model.FileVersion, error) {
	return nil, errors.New("not found")
}
func (m *mockFileVersionRepoWithError) Delete(id uint) error {
	return errors.New("not found")
}
func (m *mockFileVersionRepoWithError) DeleteByFileID(fileID uint) error {
	return nil
}

// mockVersionFileRepoWithError 模拟Update返回错误的文件仓库
type mockVersionFileRepoWithError struct {
	updateErr error
}

func (m *mockVersionFileRepoWithError) FindByID(id uint) (*model.File, error) {
	return &model.File{ID: 1, UserID: 1, Name: "file.txt", Size: 100}, nil
}
func (m *mockVersionFileRepoWithError) Update(file *model.File) error {
	return m.updateErr
}
