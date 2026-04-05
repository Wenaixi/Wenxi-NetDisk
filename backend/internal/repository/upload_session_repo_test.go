package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// UploadSessionRepository Tests

func TestUploadSessionRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	session := &model.UploadSession{
		UserID:      1,
		FileName:    "test.txt",
		FileSize:    1024,
		ChunksTotal: 5,
		Status:      "pending",
	}

	err := repo.Create(session)
	assert.NoError(t, err)
	assert.NotZero(t, session.ID)
}

func TestUploadSessionRepository_FindByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	session := &model.UploadSession{
		UserID: 1, FileName: "findme.txt", FileSize: 2048,
		ChunksTotal: 3, Status: "pending",
	}
	repo.Create(session)

	found, err := repo.FindByID(session.ID)
	assert.NoError(t, err)
	assert.Equal(t, "findme.txt", found.FileName)
}

func TestUploadSessionRepository_FindByID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	found, err := repo.FindByID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUploadSessionRepository_FindByUserIDAndHash(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	s1 := &model.UploadSession{
		UserID: 1, FileName: "file.txt", FileSize: 100, FileHash: "abc123",
		ChunksTotal: 1, Status: "pending",
	}
	repo.Create(s1)

	// Different user, same hash
	repo.Create(&model.UploadSession{
		UserID: 2, FileName: "file.txt", FileSize: 100, FileHash: "abc123",
		ChunksTotal: 1, Status: "pending",
	})

	found, err := repo.FindByUserIDAndHash(1, "abc123")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), found.UserID)
}

func TestUploadSessionRepository_Update(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	session := &model.UploadSession{
		UserID: 1, FileName: "update.txt", FileSize: 100, FileHash: "h1",
		ChunksTotal: 5, ChunksUploaded: 0, Status: "pending",
	}
	repo.Create(session)

	session.ChunksUploaded = 3
	session.Status = "uploading"
	err := repo.Update(session)
	assert.NoError(t, err)

	found, _ := repo.FindByID(session.ID)
	assert.Equal(t, 3, found.ChunksUploaded)
	assert.Equal(t, "uploading", found.Status)
}

func TestUploadSessionRepository_Delete(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	session := &model.UploadSession{
		UserID: 1, FileName: "delete.txt", FileSize: 100, FileHash: "h2",
		ChunksTotal: 1, Status: "pending",
	}
	repo.Create(session)

	err := repo.Delete(session.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(session.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUploadSessionRepository_DeleteByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewUploadSessionRepository(db)

	repo.Create(&model.UploadSession{
		UserID: 1, FileName: "del1.txt", FileSize: 100, FileHash: "h3",
		ChunksTotal: 1, Status: "pending",
	})
	repo.Create(&model.UploadSession{
		UserID: 1, FileName: "del2.txt", FileSize: 200, FileHash: "h4",
		ChunksTotal: 1, Status: "pending",
	})

	err := repo.DeleteByUserID(1)
	assert.NoError(t, err)

	// All sessions for user 1 should be gone
	_, err = repo.FindByUserIDAndHash(1, "h3")
	assert.Error(t, err)
}
