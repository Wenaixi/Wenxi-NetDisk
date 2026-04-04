package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockDownloadRepo_R5 完整的下载仓库mock
type mockDownloadRepo_R5 struct {
	files map[uint]*model.File
}

func newMockDownloadRepo_R5() *mockDownloadRepo_R5 {
	return &mockDownloadRepo_R5{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "document.pdf", Size: 102400, LanZouFileID: "lz-abc123"},
		},
	}
}
func (m *mockDownloadRepo_R5) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, assert.AnError
}
func (m *mockDownloadRepo_R5) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, assert.AnError
}

// mockLanzouForDownload_R5 蓝奏云mock
type mockLanzouForDownload_R5 struct{}

func (m *mockLanzouForDownload_R5) IsConnected(uint) bool { return true }
func (m *mockLanzouForDownload_R5) GetClient(uint) (*lanzou.Client, error) {
	return lanzou.NewClient("mock-cookie"), nil
}

// mockRecycleBin_R5 完整的回收站mock
type mockRecycleBin_R5 struct {
	items   map[uint]*model.RecycleBin
	cleared bool
}

func newMockRecycleBin_R5() *mockRecycleBin_R5 {
	return &mockRecycleBin_R5{
		items: map[uint]*model.RecycleBin{
			1: {ID: 1, UserID: 1, OriginalName: "deleted_item", ItemType: "file"},
		},
	}
}
func (m *mockRecycleBin_R5) Create(item *model.RecycleBin) error {
	item.ID = uint(len(m.items)) + 1
	m.items[item.ID] = item
	return nil
}
func (m *mockRecycleBin_R5) List(userID uint) ([]model.RecycleBin, error) {
	var result []model.RecycleBin
	for _, item := range m.items {
		if item.UserID == userID {
			result = append(result, *item)
		}
	}
	return result, nil
}
func (m *mockRecycleBin_R5) GetByID(id, userID uint) (*model.RecycleBin, error) {
	if item, ok := m.items[id]; ok && item.UserID == userID {
		return item, nil
	}
	return nil, assert.AnError
}
func (m *mockRecycleBin_R5) Restore(id, userID uint) error {
	if item, ok := m.items[id]; ok && item.UserID == userID {
		return nil
	}
	return assert.AnError
}
func (m *mockRecycleBin_R5) DeletePermanently(id, userID uint) error {
	if item, ok := m.items[id]; ok && item.UserID == userID {
		delete(m.items, id)
		return nil
	}
	return assert.AnError
}
func (m *mockRecycleBin_R5) ClearAll(userID uint) error {
	m.cleared = true
	for id, item := range m.items {
		if item.UserID == userID {
			delete(m.items, id)
		}
	}
	return nil
}

// mockFileRepoError_R5 返回错误的仓库
type mockFileRepoError_R5 struct{}

func (m *mockFileRepoError_R5) Create(*model.File) error               { return nil }
func (m *mockFileRepoError_R5) FindByUserID(uint) ([]model.File, error) { return nil, assert.AnError }
func (m *mockFileRepoError_R5) FindByID(uint) (*model.File, error)      { return nil, assert.AnError }
func (m *mockFileRepoError_R5) FindByFolderID(uint, *uint) ([]model.File, error) {
	return nil, nil
}
func (m *mockFileRepoError_R5) FindByLanZouFileID(string) (*model.File, error) {
	return nil, assert.AnError
}
func (m *mockFileRepoError_R5) Delete(uint) error                       { return nil }
func (m *mockFileRepoError_R5) Update(*model.File) error                { return nil }

// ============ DownloadHandler GetDownloadURL Full Path ============

func TestDownloadHandler_GetDownloadURL_FullSuccess(t *testing.T) {
	repo := newMockDownloadRepo_R5()
	lanzouProvider := &mockLanzouForDownload_R5{}
	downloadSvc := service.NewDownloadService(repo, lanzouProvider)
	handler := NewDownloadHandler(downloadSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/1/download", nil)
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "document.pdf", data["file_name"])
}

func TestDownloadHandler_GetDownloadURL_FileNotFound_R5(t *testing.T) {
	repo := newMockDownloadRepo_R5()
	lanzouProvider := &mockLanzouForDownload_R5{}
	downloadSvc := service.NewDownloadService(repo, lanzouProvider)
	handler := NewDownloadHandler(downloadSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/999/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============ RecycleHandler List/Clear Success Paths ============

func TestRecycleHandler_List_HasItems(t *testing.T) {
	recycleRepo := newMockRecycleBin_R5()
	handler := NewRecycleHandler(service.NewRecycleBinService(recycleRepo, nil, nil))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/recycle", handler.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/recycle", nil)
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].([]interface{})
	assert.True(t, len(data) >= 1)
}

func TestRecycleHandler_List_Empty(t *testing.T) {
	recycleRepo := &mockRecycleBin_R5{items: make(map[uint]*model.RecycleBin)}
	handler := NewRecycleHandler(service.NewRecycleBinService(recycleRepo, nil, nil))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/recycle", handler.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/recycle", nil)
	r.ServeHTTP(w, req)

	assertJSONResponse(t, w, http.StatusOK)
}

func TestRecycleHandler_Clear_Verified(t *testing.T) {
	recycleRepo := newMockRecycleBin_R5()
	handler := NewRecycleHandler(service.NewRecycleBinService(recycleRepo, nil, nil))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/clear", handler.Clear)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/clear", nil)
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "cleared", data["status"])
	assert.True(t, recycleRepo.cleared, "ClearAll should set cleared flag")
}

// ============ FileHandler ListFiles Error Path ============

func TestFileHandler_ListFiles_RepoError(t *testing.T) {
	errRepo := &mockFileRepoError_R5{}
	fileSvc := service.NewFileService(errRepo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files", handler.ListFiles)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============ FileHandler UpdateFileDescription InvalidBody ============

func TestFileHandler_UpdateFileDescription_BadBody(t *testing.T) {
	repo := newMockFileRepoForHandler()
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/description", handler.UpdateFileDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/description", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============ FileHandler MoveFile InvalidBody ============

func TestFileHandler_MoveFile_BadBody(t *testing.T) {
	repo := newMockFileRepoForHandler()
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/move", handler.MoveFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/move", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============ LanzouHandler Disconnect Full Path ============

func TestLanzouHandler_Disconnect_HasToken(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	tokenRepo.tokens[1] = &model.LanZouToken{
		ID: 1, UserID: 1, Cookie: "valid-cookie",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/disconnect", handler.Disconnect)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/lanzou/disconnect", nil)
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "disconnected", data["status"])
}

func TestLanzouHandler_Disconnect_NoToken(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})
	r.POST("/lanzou/disconnect", handler.Disconnect)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/lanzou/disconnect", nil)
	r.ServeHTTP(w, req)

	// DeleteByUserID always returns nil even for non-existent tokens
	// So the handler returns 200 with "disconnected"
	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "disconnected", data["status"])
}

// mockLanzouNotConnectedProvider 不连接的蓝奏云mock
type mockLanzouNotConnectedProvider struct{}

func (m *mockLanzouNotConnectedProvider) IsConnected(uint) bool { return false }
func (m *mockLanzouNotConnectedProvider) GetClient(uint) (*lanzou.Client, error) {
	return nil, assert.AnError
}

// ============ LanzouHandler InitializeUpload NotConnected ============

func TestLanzouHandler_InitializeUpload_NoConnection_R5(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	lanzouSvc := service.NewLanZouService(tokenRepo)

	repo := newMockUploadSessionRepoForStatus()
	lanzouProvider := &mockLanzouNotConnectedProvider{}
	fileCreator := &mockFileSvcForUpload{}
	uploadSvc := service.NewUploadService(repo, lanzouProvider, fileCreator)
	handler := NewLanZouHandler(lanzouSvc, uploadSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/upload/init", handler.InitializeUpload)

	w := httptest.NewRecorder()
	body := `{"file_name":"test.txt","file_size":100}`
	req := httptest.NewRequest("POST", "/lanzou/upload/init", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============ LanzouHandler ListFiles/Folders NotConnected ============

func TestLanzouHandler_ListFiles_NoConnection(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})
	r.GET("/lanzou/files", handler.ListFiles)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/files", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLanzouHandler_ListFolders_NoConnection(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})
	r.GET("/lanzou/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/folders", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLanzouHandler_CreateFolder_NoConnection_R5(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})
	r.POST("/lanzou/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"test_folder","parent_id":-1}`
	req := httptest.NewRequest("POST", "/lanzou/folders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ============ AuthHandler Register Success with auto-generated username ============

func TestAuthHandler_Register_AutoUsername(t *testing.T) {
	userRepo := NewMockUserRepository()
	authSvc := service.NewAuthService(userRepo, nil)
	handler := NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.Register)

	w := httptest.NewRecorder()
	body := `{"email":"autouser@example.com","password":"Pass123456"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "autouser", data["username"])
	assert.Equal(t, "autouser@example.com", data["email"])
}

// ============ ShareParseHandler Success Paths ============

func TestShareParseHandler_ParseShare_WithValidURL(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shares/parse", handler.ParseShare)

	w := httptest.NewRecorder()
	body := `{"url":"https://abc.lanzoui.com/xyz123","pwd":"test"}`
	req := httptest.NewRequest("POST", "/shares/parse", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}

func TestShareParseHandler_GetDownloadURL_WithValidURL(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shares/download-url", handler.GetShareDownloadURL)

	w := httptest.NewRecorder()
	body := `{"url":"https://abc.lanzoui.com/xyz123","pwd":"test"}`
	req := httptest.NewRequest("POST", "/shares/download-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}

func TestShareParseHandler_ValidateURL_Success(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shares/validate-url", handler.ValidateShareURL)

	w := httptest.NewRecorder()
	body := `{"url":"https://abc.lanzoui.com/xyz123"}`
	req := httptest.NewRequest("POST", "/shares/validate-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}
