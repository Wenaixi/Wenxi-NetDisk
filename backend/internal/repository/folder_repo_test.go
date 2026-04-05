package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
)

// FolderRepository Tests

func TestFolderRepository_Create(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	folder := &model.Folder{UserID: 1, Name: "test_folder"}
	err := repo.Create(folder)
	assert.NoError(t, err)
	assert.NotZero(t, folder.ID)
}

func TestFolderRepository_FindByID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	folder := &model.Folder{UserID: 1, Name: "find_folder"}
	repo.Create(folder)

	found, err := repo.FindByID(folder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "find_folder", found.Name)
}

func TestFolderRepository_FindByID_NotFound(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	found, err := repo.FindByID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFolderRepository_FindByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	repo.Create(&model.Folder{UserID: 1, Name: "user1_folder1"})
	repo.Create(&model.Folder{UserID: 1, Name: "user1_folder2"})
	repo.Create(&model.Folder{UserID: 2, Name: "user2_folder"})

	folders, err := repo.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, folders, 2)
}

func TestFolderRepository_FindByParentID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	parentID := uint(5)
	repo.Create(&model.Folder{UserID: 1, Name: "child1", ParentID: &parentID})
	repo.Create(&model.Folder{UserID: 1, Name: "child2", ParentID: &parentID})
	repo.Create(&model.Folder{UserID: 1, Name: "root"})

	children, err := repo.FindByParentID(1, &parentID)
	assert.NoError(t, err)
	assert.Len(t, children, 2)

	roots, err := repo.FindByParentID(1, nil)
	assert.NoError(t, err)
	assert.Len(t, roots, 1)
}

func TestFolderRepository_Update(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	folder := &model.Folder{UserID: 1, Name: "old_name", Description: "old desc"}
	repo.Create(folder)

	folder.Name = "new_name"
	folder.Description = "new desc"
	err := repo.Update(folder)
	assert.NoError(t, err)

	found, _ := repo.FindByID(folder.ID)
	assert.Equal(t, "new_name", found.Name)
	assert.Equal(t, "new desc", found.Description)
}

func TestFolderRepository_Delete(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	folder := &model.Folder{UserID: 1, Name: "deleteme"}
	repo.Create(folder)

	err := repo.Delete(folder.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(folder.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFolderRepository_DeleteByUserID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	repo.Create(&model.Folder{UserID: 1, Name: "del1"})
	repo.Create(&model.Folder{UserID: 1, Name: "del2"})
	repo.Create(&model.Folder{UserID: 2, Name: "keep"})

	err := repo.DeleteByUserID(1)
	assert.NoError(t, err)

	user1, _ := repo.FindByUserID(1)
	user2, _ := repo.FindByUserID(2)
	assert.Len(t, user1, 0)
	assert.Len(t, user2, 1)
}

func TestFolderRepository_DeleteByParentID(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := NewFolderRepository(db)

	parentID := uint(5)
	repo.Create(&model.Folder{UserID: 1, Name: "child1", ParentID: &parentID})
	repo.Create(&model.Folder{UserID: 1, Name: "child2", ParentID: &parentID})

	err := repo.DeleteByParentID(parentID)
	assert.NoError(t, err)

	children, _ := repo.FindByParentID(1, &parentID)
	assert.Len(t, children, 0)
}
