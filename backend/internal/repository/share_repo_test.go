package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// ShareRepository Tests

func TestShareRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	share := &model.Share{UserID: 1, FileID: 1, ShareToken: "token123"}
	err := repo.Create(share)
	assert.NoError(t, err)
	assert.NotZero(t, share.ID)
}

func TestShareRepository_FindByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	share := &model.Share{UserID: 1, FileID: 1, ShareToken: "findme"}
	repo.Create(share)

	found, err := repo.FindByID(share.ID)
	assert.NoError(t, err)
	assert.Equal(t, "findme", found.ShareToken)
}

func TestShareRepository_FindByToken(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	share := &model.Share{UserID: 1, FileID: 1, ShareToken: "unique-token-abc"}
	repo.Create(share)

	found, err := repo.FindByToken("unique-token-abc")
	assert.NoError(t, err)
	assert.Equal(t, "unique-token-abc", found.ShareToken)
}

func TestShareRepository_FindByToken_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	found, err := repo.FindByToken("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestShareRepository_FindByFileID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	repo.Create(&model.Share{UserID: 1, FileID: 10, ShareToken: "tok1"})
	repo.Create(&model.Share{UserID: 1, FileID: 10, ShareToken: "tok2"})
	repo.Create(&model.Share{UserID: 1, FileID: 20, ShareToken: "tok3"})

	shares, err := repo.FindByFileID(10)
	assert.NoError(t, err)
	assert.Len(t, shares, 2)
}

func TestShareRepository_FindByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	repo.Create(&model.Share{UserID: 1, FileID: 1, ShareToken: "t1"})
	repo.Create(&model.Share{UserID: 1, FileID: 2, ShareToken: "t2"})
	repo.Create(&model.Share{UserID: 2, FileID: 3, ShareToken: "t3"})

	shares, err := repo.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, shares, 2)
}

func TestShareRepository_Delete(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	share := &model.Share{UserID: 1, FileID: 1, ShareToken: "deleteme"}
	repo.Create(share)

	err := repo.Delete(share.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(share.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestShareRepository_DeleteByFileID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewShareRepository(db)

	repo.Create(&model.Share{UserID: 1, FileID: 100, ShareToken: "t1"})
	repo.Create(&model.Share{UserID: 1, FileID: 100, ShareToken: "t2"})
	repo.Create(&model.Share{UserID: 1, FileID: 200, ShareToken: "t3"})

	err := repo.DeleteByFileID(100)
	assert.NoError(t, err)

	file100, _ := repo.FindByFileID(100)
	file200, _ := repo.FindByFileID(200)
	assert.Len(t, file100, 0)
	assert.Len(t, file200, 1)
}
