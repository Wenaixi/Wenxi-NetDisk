package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
)

func setupRepoTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.File{}, &model.Folder{}, &model.RecycleBin{},
		&model.Share{}, &model.LanZouToken{}, &model.UploadSession{}, &model.FileVersion{})
	assert.NoError(t, err)

	return db
}

// FileRepository Tests

func TestFileRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	file := &model.File{
		UserID:       1,
		Name:         "test.txt",
		Size:         1024,
		LanZouFileID: "lz-123",
	}

	err := repo.Create(file)
	assert.NoError(t, err)
	assert.NotZero(t, file.ID)
	assert.NotZero(t, file.CreatedAt)
}

func TestFileRepository_FindByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	file := &model.File{UserID: 1, Name: "findbyid.txt", Size: 2048}
	repo.Create(file)

	found, err := repo.FindByID(file.ID)
	assert.NoError(t, err)
	assert.Equal(t, "findbyid.txt", found.Name)
	assert.Equal(t, int64(2048), found.Size)
}

func TestFileRepository_FindByID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	found, err := repo.FindByID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFileRepository_FindByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	repo.Create(&model.File{UserID: 1, Name: "user1_file1.txt", Size: 100})
	repo.Create(&model.File{UserID: 1, Name: "user1_file2.txt", Size: 200})
	repo.Create(&model.File{UserID: 2, Name: "user2_file.txt", Size: 300})

	files, err := repo.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileRepository_FindByFolderID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	folderID := uint(5)
	// Use separate DB sessions per test by using in-memory DB with unique path
	// For this test, just verify folder-based query works
	f1 := &model.File{UserID: 1, Name: "in_folder.txt", Size: 100, FolderID: &folderID, EncryptionKey: "k", EncryptionNonce: "n"}
	repo.Create(f1)

	files, err := repo.FindByFolderID(1, &folderID)
	assert.NoError(t, err)
	assert.True(t, len(files) >= 1)
}

func TestFileRepository_Update(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	file := &model.File{UserID: 1, Name: "old_name.txt", Size: 100}
	repo.Create(file)

	file.Name = "new_name.txt"
	file.Size = 500
	err := repo.Update(file)
	assert.NoError(t, err)

	found, _ := repo.FindByID(file.ID)
	assert.Equal(t, "new_name.txt", found.Name)
	assert.Equal(t, int64(500), found.Size)
}

func TestFileRepository_Delete(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	file := &model.File{UserID: 1, Name: "deleteme.txt", Size: 100}
	repo.Create(file)

	err := repo.Delete(file.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(file.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFileRepository_DeleteByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileRepository(db)

	repo.Create(&model.File{UserID: 1, Name: "del1.txt", Size: 100})
	repo.Create(&model.File{UserID: 1, Name: "del2.txt", Size: 200})
	repo.Create(&model.File{UserID: 2, Name: "keep.txt", Size: 300})

	err := repo.DeleteByUserID(1)
	assert.NoError(t, err)

	user1Files, _ := repo.FindByUserID(1)
	user2Files, _ := repo.FindByUserID(2)
	assert.Len(t, user1Files, 0)
	assert.Len(t, user2Files, 1)
}
