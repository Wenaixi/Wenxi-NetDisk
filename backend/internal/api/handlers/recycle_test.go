package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecycleHandler_Restore_InvalidID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/abc/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Restore_NonNumericID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/notanumber/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Delete_InvalidID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Delete_NonNumericID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Restore_FloatingPointID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/12.5/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Delete_FloatingPointID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Restore_NegativeID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/-1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Delete_NegativeID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Restore_EmptyID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle//restore", nil)
	r.ServeHTTP(w, req)

	// Empty ID results in strconv.ParseUint error on empty string -> 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecycleHandler_Delete_EmptyID(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRecycleHandler_List(t *testing.T) {
	svc := service.NewRecycleBinService(nil, nil, nil)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/recycle", handler.List)

	// List requires no path params, just user_id in context
	// But service will panic with nil repo, so we just verify handler setup
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.recycleSvc)
}

// mockRecycleSvcForHandler 模拟回收站服务用于handler测试
type mockRecycleSvcForHandler struct{}

func (m *mockRecycleSvcForHandler) List(userID uint) ([]interface{}, error) {
	return []interface{}{}, nil
}
func (m *mockRecycleSvcForHandler) Restore(userID uint, recycleID uint) error {
	return nil
}
func (m *mockRecycleSvcForHandler) PermanentDelete(userID uint, recycleID uint) error {
	return nil
}
func (m *mockRecycleSvcForHandler) ClearAll(userID uint) error {
	return nil
}

func TestRecycleHandler_List_WithMockRepos(t *testing.T) {
	// Use real service with mock-like approach
	// RecycleBinService uses concrete type, so test with real repo
	fileRepo := &mockRecycleFileRepoForHandler{}
	folderRepo := &mockRecycleFolderRepoForHandler{}
	recycleRepo := &mockRecycleBinRepoForHandler{}
	svc := service.NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	handler := NewRecycleHandler(svc)

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

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecycleHandler_Restore_WithMockRepos(t *testing.T) {
	fileRepo := &mockRecycleFileRepoForHandler{}
	folderRepo := &mockRecycleFolderRepoForHandler{}
	recycleRepo := &mockRecycleBinRepoForHandler{}
	svc := service.NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/recycle/:id/restore", handler.Restore)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/recycle/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecycleHandler_Delete_WithMockRepos(t *testing.T) {
	fileRepo := &mockRecycleFileRepoForHandler{}
	folderRepo := &mockRecycleFolderRepoForHandler{}
	recycleRepo := &mockRecycleBinRepoForHandler{}
	svc := service.NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	handler := NewRecycleHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecycleHandler_Clear_WithMockRepos(t *testing.T) {
	fileRepo := &mockRecycleFileRepoForHandler{}
	folderRepo := &mockRecycleFolderRepoForHandler{}
	recycleRepo := &mockRecycleBinRepoForHandler{}
	svc := service.NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	handler := NewRecycleHandler(svc)

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

	assert.Equal(t, http.StatusOK, w.Code)
}

// mock repositories for handler tests
type mockRecycleBinRepoForHandler struct{}

func (m *mockRecycleBinRepoForHandler) Create(item *model.RecycleBin) error {
	return nil
}
func (m *mockRecycleBinRepoForHandler) List(userID uint) ([]model.RecycleBin, error) {
	return []model.RecycleBin{}, nil
}
func (m *mockRecycleBinRepoForHandler) GetByID(id uint, userID uint) (*model.RecycleBin, error) {
	return &model.RecycleBin{ID: id, UserID: userID, OriginalName: "test.txt", ItemType: "file"}, nil
}
func (m *mockRecycleBinRepoForHandler) Restore(id uint, userID uint) error {
	return nil
}
func (m *mockRecycleBinRepoForHandler) DeletePermanently(id uint, userID uint) error {
	return nil
}
func (m *mockRecycleBinRepoForHandler) ClearAll(userID uint) error {
	return nil
}

type mockRecycleFileRepoForHandler struct{}

func (m *mockRecycleFileRepoForHandler) FindByID(id uint) (*model.File, error) {
	return &model.File{ID: id, UserID: 1, Name: "test.txt"}, nil
}
func (m *mockRecycleFileRepoForHandler) Create(file *model.File) error {
	return nil
}
func (m *mockRecycleFileRepoForHandler) Delete(id uint) error {
	return nil
}

type mockRecycleFolderRepoForHandler struct{}

func (m *mockRecycleFolderRepoForHandler) FindByID(id uint) (*model.Folder, error) {
	return &model.Folder{ID: id, UserID: 1, Name: "test-folder"}, nil
}
func (m *mockRecycleFolderRepoForHandler) Create(folder *model.Folder) error {
	return nil
}
func (m *mockRecycleFolderRepoForHandler) Delete(id uint) error {
	return nil
}
