package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// FileVersionRepository Tests

func TestFileVersionRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	version := &model.FileVersion{
		FileID:   1,
		UserID:   1,
		Size:     1024,
		EncryptionKey:   "key1",
		EncryptionNonce: "nonce1",
	}

	err := repo.Create(version)
	assert.NoError(t, err)
	assert.NotZero(t, version.ID)
}

func TestFileVersionRepository_FindByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	version := &model.FileVersion{
		FileID: 1, UserID: 1, Size: 2048,
		EncryptionKey: "k", EncryptionNonce: "n",
	}
	repo.Create(version)

	found, err := repo.FindByID(version.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2048), found.Size)
}

func TestFileVersionRepository_FindByID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	found, err := repo.FindByID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFileVersionRepository_FindByFileID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	repo.Create(&model.FileVersion{FileID: 1, UserID: 1, Size: 100, EncryptionKey: "k", EncryptionNonce: "n"})
	repo.Create(&model.FileVersion{FileID: 1, UserID: 1, Size: 200, EncryptionKey: "k", EncryptionNonce: "n"})
	repo.Create(&model.FileVersion{FileID: 2, UserID: 1, Size: 300, EncryptionKey: "k", EncryptionNonce: "n"})

	versions, err := repo.FindByFileID(1)
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
}

func TestFileVersionRepository_Delete(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	version := &model.FileVersion{
		FileID: 1, UserID: 1, Size: 100,
		EncryptionKey: "k", EncryptionNonce: "n",
	}
	repo.Create(version)

	err := repo.Delete(version.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(version.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFileVersionRepository_DeleteByFileID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFileVersionRepository(db)

	repo.Create(&model.FileVersion{FileID: 1, UserID: 1, Size: 100, EncryptionKey: "k", EncryptionNonce: "n"})
	repo.Create(&model.FileVersion{FileID: 1, UserID: 1, Size: 200, EncryptionKey: "k", EncryptionNonce: "n"})

	err := repo.DeleteByFileID(1)
	assert.NoError(t, err)

	versions, _ := repo.FindByFileID(1)
	assert.Len(t, versions, 0)
}
