package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShareHandler_CreateShare_InvalidID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShare_InvalidJSON(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/123", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShare_NonNumericID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/notanumber", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShare_FloatingPointID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/12.5", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShare_NegativeID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/-1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_ValidateShare_InvalidJSON(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:token/validate", handler.ValidateShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares/testtoken/validate", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_DeleteShare_InvalidID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/shares/:id", handler.DeleteShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/shares/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_DeleteShare_NonNumericID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/shares/:id", handler.DeleteShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/shares/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_DeleteShare_FloatingPointID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/shares/:id", handler.DeleteShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/shares/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_DeleteShare_NegativeID(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/shares/:id", handler.DeleteShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/shares/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShareViaBody_InvalidJSON(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares", handler.CreateShareViaBody)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareHandler_CreateShareViaBody_EmptyBody(t *testing.T) {
	svc := service.NewShareService(nil, nil)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares", handler.CreateShareViaBody)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/shares", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// mockShareRepoForHandler implements ShareRepository
type mockShareRepoForHandler struct {
	shares map[uint]*model.Share
}

func newMockShareRepoForHandler() *mockShareRepoForHandler {
	return &mockShareRepoForHandler{
		shares: map[uint]*model.Share{
			1: {ID: 1, UserID: 1, FileID: 1, ShareToken: "abc123", PasswordHash: nil, ExpiresAt: nil, File: model.File{ID: 1, UserID: 1, Name: "test.txt", Size: 1024}},
		},
	}
}

func (m *mockShareRepoForHandler) Create(share *model.Share) error {
	share.ID = uint(len(m.shares)) + 1
	m.shares[share.ID] = share
	return nil
}
func (m *mockShareRepoForHandler) FindByID(id uint) (*model.Share, error) {
	if s, ok := m.shares[id]; ok {
		return s, nil
	}
	return nil, errors.New("share not found")
}
func (m *mockShareRepoForHandler) FindByToken(token string) (*model.Share, error) {
	for _, s := range m.shares {
		if s.ShareToken == token {
			return s, nil
		}
	}
	return nil, errors.New("share not found")
}
func (m *mockShareRepoForHandler) FindByFileID(fileID uint) ([]model.Share, error) {
	var result []model.Share
	for _, s := range m.shares {
		if s.FileID == fileID {
			result = append(result, *s)
		}
	}
	return result, nil
}
func (m *mockShareRepoForHandler) FindByUserID(userID uint) ([]model.Share, error) {
	var result []model.Share
	for _, s := range m.shares {
		if s.UserID == userID {
			result = append(result, *s)
		}
	}
	return result, nil
}
func (m *mockShareRepoForHandler) Delete(id uint) error {
	delete(m.shares, id)
	return nil
}
func (m *mockShareRepoForHandler) DeleteByFileID(fileID uint) error {
	for id, s := range m.shares {
		if s.FileID == fileID {
			delete(m.shares, id)
		}
	}
	return nil
}

// mockShareFileRepoForHandler implements ShareFileRepository
type mockShareFileRepoForHandler struct {
	files map[uint]*model.File
}

func newMockShareFileRepoForHandler() *mockShareFileRepoForHandler {
	return &mockShareFileRepoForHandler{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "test.txt", Size: 1024},
		},
	}
}
func (m *mockShareFileRepoForHandler) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, errors.New("file not found")
}

// mockShareRepoListError implements ShareRepository with List error
type mockShareRepoListError struct{}

func (m *mockShareRepoListError) Create(share *model.Share) error { return nil }
func (m *mockShareRepoListError) FindByID(id uint) (*model.Share, error) {
	return nil, errors.New("not found")
}
func (m *mockShareRepoListError) FindByToken(token string) (*model.Share, error) {
	return nil, errors.New("not found")
}
func (m *mockShareRepoListError) FindByFileID(fileID uint) ([]model.Share, error) {
	return []model.Share{}, nil
}
func (m *mockShareRepoListError) FindByUserID(userID uint) ([]model.Share, error) {
	return nil, errors.New("db query failed")
}
func (m *mockShareRepoListError) Delete(id uint) error { return nil }
func (m *mockShareRepoListError) DeleteByFileID(fileID uint) error { return nil }

// mockShareFileRepoError implements ShareFileRepository with error
type mockShareFileRepoError struct{}
func (m *mockShareFileRepoError) FindByID(id uint) (*model.File, error) {
	return nil, errors.New("file not found")
}

// TestShareHandler_CreateShare_Success_Path 测试创建分享成功
func TestShareHandler_CreateShare_Success_Path(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	body := `{}`
	req := httptest.NewRequest("POST", "/shares/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestShareHandler_CreateShare_FileNotFound 测试文件不存在
func TestShareHandler_CreateShare_FileNotFound(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:id", handler.CreateShare)

	w := httptest.NewRecorder()
	body := `{}`
	req := httptest.NewRequest("POST", "/shares/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestShareHandler_GetShare_NotFound_Path 测试分享不存在
func TestShareHandler_GetShare_NotFound_Path(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/shares/:token", handler.GetShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shares/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestShareHandler_ListShares_Error 测试分享列表错误
func TestShareHandler_ListShares_Error(t *testing.T) {
	shareRepo := &mockShareRepoListError{}
	fileRepo := &mockShareFileRepoError{}
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/shares", handler.ListShares)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shares", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestShareHandler_DeleteShare_NotFound_Path 测试删除不存在的分享
func TestShareHandler_DeleteShare_NotFound_Path(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/shares/:id", handler.DeleteShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/shares/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestShareHandler_CreateShareViaBody_Success_Path 测试通过body创建分享成功
func TestShareHandler_CreateShareViaBody_Success_Path(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares", handler.CreateShareViaBody)

	w := httptest.NewRecorder()
	body := `{"file_id":1}`
	req := httptest.NewRequest("POST", "/shares", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestShareHandler_CreateShareViaBody_FileNotFound 测试文件不存在
func TestShareHandler_CreateShareViaBody_FileNotFound(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares", handler.CreateShareViaBody)

	w := httptest.NewRecorder()
	body := `{"file_id":999}`
	req := httptest.NewRequest("POST", "/shares", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestShareHandler_ValidateShare_NotFound 测试分享不存在
func TestShareHandler_ValidateShare_NotFound(t *testing.T) {
	shareRepo := newMockShareRepoForHandler()
	fileRepo := newMockShareFileRepoForHandler()
	svc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/:token/validate", handler.ValidateShare)

	w := httptest.NewRecorder()
	body := `{"password":""}`
	req := httptest.NewRequest("POST", "/shares/nonexistent/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
