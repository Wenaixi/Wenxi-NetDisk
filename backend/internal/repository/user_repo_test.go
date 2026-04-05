package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("skipping test: cannot connect to sqlite (CGO required): %v", err)
		return nil
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "findbyid",
		Email:        "findbyid@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	found, err := repo.FindByID(99999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "findbyemail",
		Email:        "unique@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByEmail("unique@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Username, found.Username)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	found, err := repo.FindByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "updateme",
		Email:        "update@example.com",
		PasswordHash: "oldhash",
	}
	repo.Create(user)

	user.Username = "updateduser"
	user.PasswordHash = "newhash"
	err := repo.Update(user)
	assert.NoError(t, err)

	found, _ := repo.FindByID(user.ID)
	assert.Equal(t, "updateduser", found.Username)
	assert.Equal(t, "newhash", found.PasswordHash)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "deleteme",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	}
	repo.Create(user)

	err := repo.Delete(user.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_EmailUnique(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user1 := &model.User{
		Username:     "user1",
		Email:        "same@example.com",
		PasswordHash: "hash1",
	}
	err := repo.Create(user1)
	assert.NoError(t, err)

	user2 := &model.User{
		Username:     "user2",
		Email:        "same@example.com",
		PasswordHash: "hash2",
	}
	err = repo.Create(user2)
	assert.Error(t, err)
}

func TestUserRepository_UsernameUnique(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user1 := &model.User{
		Username:     "uniqueuser",
		Email:        "user1@example.com",
		PasswordHash: "hash1",
	}
	err := repo.Create(user1)
	assert.NoError(t, err)

	user2 := &model.User{
		Username:     "uniqueuser",
		Email:        "user2@example.com",
		PasswordHash: "hash2",
	}
	err = repo.Create(user2)
	assert.Error(t, err)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "findbyusername",
		Email:        "finduser@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByUsername("findbyusername")
	assert.NoError(t, err)
	assert.Equal(t, "findbyusername", found.Username)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	repo := NewUserRepository(db)

	found, err := repo.FindByUsername("nonexistent_user")
	assert.Error(t, err)
	assert.Nil(t, found)
}
