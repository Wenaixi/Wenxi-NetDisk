package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupTestRouterWithUser() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	return r
}

// TestAuthHandler_Register_Validation 测试注册handler的输入验证
func TestAuthHandler_Register_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{"invalid email format", `{"email":"invalid","password":"Test123456"}`, http.StatusBadRequest},
		{"short password", `{"email":"test@example.com","password":"123"}`, http.StatusBadRequest},
		{"missing email", `{"password":"Test123456"}`, http.StatusBadRequest},
		{"missing password", `{"email":"test@example.com"}`, http.StatusBadRequest},
		{"short username", `{"email":"test@example.com","username":"ab","password":"Test123456"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			handler := &AuthHandler{}
			r.POST("/register", handler.Register)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

// TestAuthHandler_Login_Validation 测试登录handler的输入验证
func TestAuthHandler_Login_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{"missing password", `{"email":"test@example.com"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
		{"neither email nor username", `{"password":"Test123456"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			handler := &AuthHandler{}
			r.POST("/login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectCode, w.Code, w.Body.String())
			}
		})
	}

	t.Run("email+password passes validation", func(t *testing.T) {
		r := setupTestRouter()
		handler := &AuthHandler{}
		r.POST("/login", handler.Login)

		w := httptest.NewRecorder()
		body := `{"email":"test@example.com","password":"Test123456"}`
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		defer func() {
			if r := recover(); r != nil {
				// Expected panic because authSvc is nil, but validation passed
			}
		}()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusBadRequest {
			t.Error("validation should have passed")
		}
	})

	t.Run("username+password passes validation", func(t *testing.T) {
		r := setupTestRouter()
		handler := &AuthHandler{}
		r.POST("/login", handler.Login)

		w := httptest.NewRecorder()
		body := `{"username":"testuser","password":"Test123456"}`
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		defer func() {
			if r := recover(); r != nil {
				// Expected panic because authSvc is nil, but validation passed
			}
		}()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusBadRequest {
			t.Error("validation should have passed")
		}
	})
}

// TestShareHandler_CreateShareViaBody_Validation 测试CreateShareViaBody的验证
func TestShareHandler_CreateShareViaBody_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{"missing file_id", `{}`, http.StatusBadRequest},
		{"invalid json", `not json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouterWithUser()
			handler := &ShareHandler{}
			r.POST("/shares", handler.CreateShareViaBody)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/shares", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

// TestShareHandler_CreateShare_Validation 测试CreateShare的验证
func TestShareHandler_CreateShare_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &ShareHandler{}
		r.POST("/files/:id/share", handler.CreateShare)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/files/abc/share", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFolderHandler_UpdateFolder_Validation 测试文件夹更新handler的验证
func TestFolderHandler_UpdateFolder_Validation(t *testing.T) {
	t.Run("should fail with invalid folder id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FolderHandler{}
		r.PUT("/folders/abc", handler.UpdateFolder)

		w := httptest.NewRecorder()
		body := `{"name":"test"}`
		req := httptest.NewRequest("PUT", "/folders/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should return 400 with empty body", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FolderHandler{}
		r.PUT("/folders/1", handler.UpdateFolder)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("PUT", "/folders/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Empty body returns 400 because neither name nor description is provided
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFolderHandler_UpdateFolderDescription_Validation 测试文件夹描述更新验证
func TestFolderHandler_UpdateFolderDescription_Validation(t *testing.T) {
	t.Run("should fail with invalid folder id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FolderHandler{}
		r.PUT("/folders/abc/description", handler.UpdateFolderDescription)

		w := httptest.NewRecorder()
		body := `{"description":"test"}`
		req := httptest.NewRequest("PUT", "/folders/abc/description", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestUploadHandler_UploadFile_Validation 测试文件上传handler验证
func TestUploadHandler_UploadFile_Validation(t *testing.T) {
	t.Run("should fail with no file", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.POST("/files/upload", handler.UploadFile)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/files/upload", strings.NewReader(""))
		req.Header.Set("Content-Type", "multipart/form-data")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestUploadHandler_GetUploadURL_Validation 测试获取上传URL验证
func TestUploadHandler_GetUploadURL_Validation(t *testing.T) {
	t.Run("should fail with missing file_name", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.POST("/files/upload-url", handler.GetUploadURL)

		w := httptest.NewRecorder()
		body := `{"size":1000}`
		req := httptest.NewRequest("POST", "/files/upload-url", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should fail with missing size", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.POST("/files/upload-url", handler.GetUploadURL)

		w := httptest.NewRecorder()
		body := `{"file_name":"test.txt"}`
		req := httptest.NewRequest("POST", "/files/upload-url", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestUploadHandler_ListVersions_Validation 测试版本列表验证
func TestUploadHandler_ListVersions_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.GET("/files/abc/versions", handler.ListVersions)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/files/abc/versions", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestUploadHandler_RestoreVersion_Validation 测试版本恢复验证
func TestUploadHandler_RestoreVersion_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.POST("/files/abc/versions/1/restore", handler.RestoreVersion)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/files/abc/versions/1/restore", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should fail with invalid version id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &UploadHandler{}
		r.POST("/files/1/versions/abc/restore", handler.RestoreVersion)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/files/1/versions/abc/restore", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestDownloadHandler_GetDownloadURL_Validation 测试下载URL验证
func TestDownloadHandler_GetDownloadURL_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &DownloadHandler{}
		r.GET("/files/abc/download", handler.GetDownloadURL)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/files/abc/download", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestHandler_ResponseFormat 验证handler响应格式正确
func TestHandler_ResponseFormat(t *testing.T) {
	t.Run("should return proper error format for bad request", func(t *testing.T) {
		r := setupTestRouter()
		handler := &AuthHandler{}
		r.POST("/register", handler.Register)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("response should be valid JSON: %v", err)
		}

		if _, ok := resp["code"]; !ok {
			t.Error("response should contain 'code' field")
		}
		if _, ok := resp["msg"]; !ok {
			t.Error("response should contain 'msg' field")
		}
	})
}

// mockUserRepoCreateError 模拟Create用户返回错误的仓库
type mockUserRepoCreateError struct{}

func (m *mockUserRepoCreateError) Create(user *model.User) error {
	return errors.New("db write failed")
}
func (m *mockUserRepoCreateError) FindByID(id uint) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoCreateError) FindByUsername(username string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoCreateError) FindByEmail(email string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoCreateError) Update(user *model.User) error { return nil }
func (m *mockUserRepoCreateError) Delete(id uint) error { return nil }

// TestAuthHandler_Register_DBError 测试注册时数据库错误
func TestAuthHandler_Register_DBError(t *testing.T) {
	userRepo := &mockUserRepoCreateError{}
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.Register)

	w := httptest.NewRecorder()
	body := `{"email":"new@test.com","password":"SecurePass123"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestAuthHandler_Register_WithEmail 测试带邮箱前缀的用户名自动派生
func TestAuthHandler_Register_WithEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.Register)

	w := httptest.NewRecorder()
	body := `{"email":"myuser@test.com","password":"SecurePass123"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	data := resp["data"].(map[string]interface{})
	if data["username"] != "myuser" {
		t.Errorf("expected username 'myuser', got '%v'", data["username"])
	}
}

// TestAuthHandler_Login_WrongCredentials 测试登录密码错误
func TestAuthHandler_Login_WrongCredentials(t *testing.T) {
	userRepo := NewMockUserRepository()
	userRepo.users[1] = &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
	}
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", handler.Login)

	w := httptest.NewRecorder()
	body := `{"email":"test@example.com","password":"WrongPassword"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestAuthHandler_Login_UserNotFound 测试登录用户不存在
func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	userRepo := NewMockUserRepository()
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", handler.Login)

	w := httptest.NewRecorder()
	body := `{"email":"nouser@example.com","password":"SomePass"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestAuthHandler_GetCurrentUser_Success_Path 测试获取当前用户成功
func TestAuthHandler_GetCurrentUser_Success_Path(t *testing.T) {
	userRepo := NewMockUserRepository()
	userRepo.users[1] = &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/auth/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/me", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	data := resp["data"].(map[string]interface{})
	if data["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", data["username"])
	}
	if data["email"] != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%v'", data["email"])
	}
}

// mockUserRepoNotFoundError 模拟FindByID返回错误的仓库
type mockUserRepoNotFoundError struct{}

func (m *mockUserRepoNotFoundError) Create(user *model.User) error { return nil }
func (m *mockUserRepoNotFoundError) FindByID(id uint) (*model.User, error) {
	return nil, errors.New("user not found")
}
func (m *mockUserRepoNotFoundError) FindByUsername(username string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoNotFoundError) FindByEmail(email string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoNotFoundError) Update(user *model.User) error { return nil }
func (m *mockUserRepoNotFoundError) Delete(id uint) error { return nil }

// TestAuthHandler_GetCurrentUser_NotFound_Path 测试获取不存在的用户
func TestAuthHandler_GetCurrentUser_NotFound_Path(t *testing.T) {
	userRepo := &mockUserRepoNotFoundError{}
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})
	r.GET("/auth/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/me", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// mockUserRepoDuplicate 模拟邮箱/用户名已存在的仓库
type mockUserRepoDuplicate struct{}

func (m *mockUserRepoDuplicate) Create(user *model.User) error { return nil }
func (m *mockUserRepoDuplicate) FindByID(id uint) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoDuplicate) FindByUsername(username string) (*model.User, error) {
	return &model.User{ID: 1, Username: username, Email: "existing@test.com"}, nil
}
func (m *mockUserRepoDuplicate) FindByEmail(email string) (*model.User, error) {
	return &model.User{ID: 1, Username: "existing", Email: email}, nil
}
func (m *mockUserRepoDuplicate) Update(user *model.User) error { return nil }
func (m *mockUserRepoDuplicate) Delete(id uint) error { return nil }

// TestAuthHandler_Register_DuplicateEmail 测试重复邮箱注册
func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	userRepo := &mockUserRepoDuplicate{}
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.Register)

	w := httptest.NewRecorder()
	body := `{"email":"existing@test.com","password":"SecurePass123"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
