package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// RecycleBinRepository Tests

func TestRecycleBinRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	item := &model.RecycleBin{
		UserID:       1,
		OriginalName: "deleted_file.txt",
		ItemType:     "file",
		ItemID:       1,
		DeletedAt:    time.Now(),
		ExpiresAt:    time.Now().AddDate(0, 0, 30),
	}

	err := repo.Create(item)
	assert.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestRecycleBinRepository_List(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "file1", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "file2", ItemType: "file", ItemID: 2,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	repo.Create(&model.RecycleBin{
		UserID: 2, OriginalName: "other", ItemType: "file", ItemID: 3,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})

	items, err := repo.List(1)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestRecycleBinRepository_GetByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	item := &model.RecycleBin{
		UserID: 1, OriginalName: "getme", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	}
	repo.Create(item)

	found, err := repo.GetByID(item.ID, 1)
	assert.NoError(t, err)
	assert.Equal(t, "getme", found.OriginalName)
}

func TestRecycleBinRepository_GetByID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	found, err := repo.GetByID(99999, 1)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestRecycleBinRepository_Restore(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	item := &model.RecycleBin{
		UserID: 1, OriginalName: "restore_me", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	}
	repo.Create(item)

	err := repo.Restore(item.ID, 1)
	assert.NoError(t, err)

	found, err := repo.GetByID(item.ID, 1)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestRecycleBinRepository_DeletePermanently(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	item := &model.RecycleBin{
		UserID: 1, OriginalName: "perm_delete", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	}
	repo.Create(item)

	err := repo.DeletePermanently(item.ID, 1)
	assert.NoError(t, err)
}

func TestRecycleBinRepository_ClearAll(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "clear1", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "clear2", ItemType: "file", ItemID: 2,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})

	err := repo.ClearAll(1)
	assert.NoError(t, err)

	items, _ := repo.List(1)
	assert.Len(t, items, 0)
}

func TestRecycleBinRepository_CleanExpired(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	// Expired item
	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "expired", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now().AddDate(0, 0, -31),
		ExpiresAt: time.Now().AddDate(0, 0, -1),
	})
	// Valid item
	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "valid", ItemType: "file", ItemID: 2,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})

	err := repo.CleanExpired()
	assert.NoError(t, err)

	items, _ := repo.List(1)
	assert.Len(t, items, 1)
	assert.Equal(t, "valid", items[0].OriginalName)
}

func TestRecycleBinRepository_GetByItemType(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewRecycleBinRepository(db)

	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "file_item", ItemType: "file", ItemID: 1,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	repo.Create(&model.RecycleBin{
		UserID: 1, OriginalName: "folder_item", ItemType: "folder", ItemID: 2,
		DeletedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
	})

	files, err := repo.GetByItemType(1, "file")
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	folders, err := repo.GetByItemType(1, "folder")
	assert.NoError(t, err)
	assert.Len(t, folders, 1)
}
