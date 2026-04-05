package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// LanZouTokenRepository Tests

func TestLanZouTokenRepository_Upsert(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewLanZouTokenRepository(db)

	token := &model.LanZouToken{
		UserID:    1,
		Cookie:    "test-cookie",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := repo.Upsert(token)
	assert.NoError(t, err)
	assert.NotZero(t, token.ID)
}

func TestLanZouTokenRepository_FindByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewLanZouTokenRepository(db)

	token := &model.LanZouToken{
		UserID:    1,
		Cookie:    "find-me-cookie",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	repo.Upsert(token)

	found, err := repo.FindByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, "find-me-cookie", found.Cookie)
	assert.Equal(t, uint(1), found.UserID)
}

func TestLanZouTokenRepository_FindByUserID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewLanZouTokenRepository(db)

	found, err := repo.FindByUserID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestLanZouTokenRepository_DeleteByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewLanZouTokenRepository(db)

	token := &model.LanZouToken{
		UserID:    1,
		Cookie:    "delete-me-cookie",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	repo.Upsert(token)

	err := repo.DeleteByUserID(1)
	assert.NoError(t, err)

	found, err := repo.FindByUserID(1)
	assert.Error(t, err)
	assert.Nil(t, found)
}
