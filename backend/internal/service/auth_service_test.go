package service

import (
	"errors"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/crypto"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/jwt"
	"gorm.io/gorm"
)

// mockUserRepo 模拟用户仓库
type mockUserRepo struct {
	users []*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: []*model.User{
			{
				ID:           1,
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "$2a$10$hashedpassword",
			},
		},
	}
}

func (m *mockUserRepo) Create(user *model.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Update(user *model.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepo) Delete(id uint) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

// mockJWTManager 模拟JWT管理器 - unused, using real JWTManager instead
type mockJWTManager struct {
	secret string
}

// TestAuthService_Register tests user registration
func TestAuthService_Register(t *testing.T) {
	t.Run("should register new user successfully", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.Register("newuser", "new@example.com", "password123")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user, got nil")
			return
		}
		if user.Username != "newuser" {
			t.Errorf("expected username 'newuser', got '%s'", user.Username)
		}
		if user.Email != "new@example.com" {
			t.Errorf("expected email 'new@example.com', got '%s'", user.Email)
		}
		if user.PasswordHash == "" {
			t.Error("expected password hash to be set")
		}
	})

	t.Run("should fail if username exists", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.Register("testuser", "another@example.com", "password123")

		if err == nil {
			t.Error("expected duplicate username error")
		}
		if user != nil {
			t.Error("expected nil user on duplicate username")
		}
	})

	t.Run("should fail if email exists", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.Register("anotheruser", "test@example.com", "password123")

		if err == nil {
			t.Error("expected duplicate email error")
		}
		if user != nil {
			t.Error("expected nil user on duplicate email")
		}
	})
}

// TestAuthService_GetUserByID tests getting user by ID
func TestAuthService_GetUserByID(t *testing.T) {
	t.Run("should return existing user", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.GetUserByID(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user, got nil")
			return
		}
		if user.Username != "testuser" {
			t.Errorf("expected username 'testuser', got '%s'", user.Username)
		}
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.GetUserByID(999)

		if err == nil {
			t.Error("expected error for non-existent user")
		}
		if user != nil {
			t.Error("expected nil user for non-existent ID")
		}
	})
}

// TestAuthService_ValidatePassword tests password validation
func TestAuthService_ValidatePassword(t *testing.T) {
	t.Run("should hash and verify password correctly", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.Register("passuser", "pass@example.com", "strongpassword")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
			return
		}

		if !crypto.CheckPassword("strongpassword", user.PasswordHash) {
			t.Error("expected password to match")
		}
	})

	t.Run("should not match wrong password", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		user, err := svc.Register("matchuser", "match@example.com", "correctpassword")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
			return
		}

		if crypto.CheckPassword("wrongpassword", user.PasswordHash) {
			t.Error("expected wrong password to not match")
		}
	})
}

// TestAuthService_Login tests user login
func TestAuthService_Login(t *testing.T) {
	t.Run("should login with email successfully", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		hashedPassword, err := crypto.HashPassword("Test123456")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		repo.users[0].PasswordHash = hashedPassword

		token, err := svc.Login("test@example.com", "Test123456")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token, got empty string")
		}
	})

	t.Run("should login with username successfully", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		hashedPassword, err := crypto.HashPassword("Test123456")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		repo.users[0].PasswordHash = hashedPassword

		token, err := svc.Login("testuser", "Test123456")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token, got empty string")
		}
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		hashedPassword, err := crypto.HashPassword("CorrectPassword")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		repo.users[0].PasswordHash = hashedPassword

		token, err := svc.Login("test@example.com", "WrongPassword")

		if err == nil {
			t.Error("expected invalid credentials error")
		}
		if token != "" {
			t.Error("expected empty token on failed login")
		}
	})

	t.Run("should fail with non-existent email", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		token, err := svc.Login("nonexistent@example.com", "Test123456")

		if err == nil {
			t.Error("expected invalid credentials error")
		}
		if token != "" {
			t.Error("expected empty token on failed login")
		}
	})

	t.Run("should fail with non-existent username", func(t *testing.T) {
		repo := newMockUserRepo()
		jwtManager := jwt.NewJWTManager("testsecret", 1440)
		svc := NewAuthService(repo, jwtManager)

		token, err := svc.Login("nonexistentuser", "Test123456")

		if err == nil {
			t.Error("expected invalid credentials error")
		}
		if token != "" {
			t.Error("expected empty token on failed login")
		}
	})
}

// mockFileRepo 模拟文件仓库 (for complete file service test)
type mockFileRepoComplete struct {
	files []*model.File
}

func newMockFileRepoComplete() *mockFileRepoComplete {
	return &mockFileRepoComplete{
		files: []*model.File{
			{ID: 1, UserID: 1, Name: "file1.txt", Size: 100},
			{ID: 2, UserID: 1, Name: "file2.pdf", Size: 200},
			{ID: 3, UserID: 2, Name: "file3.doc", Size: 300},
		},
	}
}

func (m *mockFileRepoComplete) Create(file *model.File) error {
	file.ID = uint(len(m.files) + 1)
	m.files = append(m.files, file)
	return nil
}

func (m *mockFileRepoComplete) FindByUserID(userID uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}

func (m *mockFileRepoComplete) FindByID(id uint) (*model.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockFileRepoComplete) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return m.FindByUserID(userID)
}

func (m *mockFileRepoComplete) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func (m *mockFileRepoComplete) Delete(id uint) error {
	for i, f := range m.files {
		if f.ID == id {
			m.files = append(m.files[:i], m.files[i+1:]...)
			return nil
		}
	}
	return errors.New("file not found")
}

func (m *mockFileRepoComplete) Update(file *model.File) error {
	for i, f := range m.files {
		if f.ID == file.ID {
			m.files[i] = file
			return nil
		}
	}
	return errors.New("file not found")
}

func (m *mockFileRepoComplete) UpdateFileName(id uint, name string) error {
	for i, f := range m.files {
		if f.ID == id {
			m.files[i].Name = name
			return nil
		}
	}
	return errors.New("file not found")
}

func (m *mockFileRepoComplete) MoveFile(id uint, folderID *uint) error {
	for i, f := range m.files {
		if f.ID == id {
			m.files[i].FolderID = folderID
			return nil
		}
	}
	return errors.New("file not found")
}
