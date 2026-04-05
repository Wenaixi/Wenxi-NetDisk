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

func TestFileHandler_GetFile_InvalidID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_GetFile_NonNumericID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_GetFile_FloatingPointID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_GetFile_NegativeID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_DeleteFile_InvalidID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/files/:id", handler.DeleteFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/files/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_DeleteFile_NonNumericID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/files/:id", handler.DeleteFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/files/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_DeleteFile_FloatingPointID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/files/:id", handler.DeleteFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/files/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_DeleteFile_NegativeID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/files/:id", handler.DeleteFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/files/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_MoveFile_InvalidID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/move", handler.MoveFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/abc/move", strings.NewReader(`{"folder_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_MoveFile_FloatingPointID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/move", handler.MoveFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/12.5/move", strings.NewReader(`{"folder_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_MoveFile_NegativeID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/move", handler.MoveFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/-1/move", strings.NewReader(`{"folder_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_MoveFile_InvalidJSON(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/1/move", handler.MoveFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/move", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_UpdateFileDescription_InvalidID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/description", handler.UpdateFileDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/abc/description", strings.NewReader(`{"description":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_UpdateFileDescription_FloatingPointID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/description", handler.UpdateFileDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/12.5/description", strings.NewReader(`{"description":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_UpdateFileDescription_InvalidJSON(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/1/description", handler.UpdateFileDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/description", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_RenameFile_InvalidID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id", handler.RenameFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/abc", strings.NewReader(`{"name":"test.txt"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_RenameFile_MissingName(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/1", handler.RenameFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_RenameFile_FloatingPointID(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id", handler.RenameFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/12.5", strings.NewReader(`{"name":"test.txt"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileHandler_RenameFile_InvalidJSON(t *testing.T) {
	svc := service.NewFileService(nil)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/1", handler.RenameFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// mockFileRepoUpdateError 模拟Update返回错误的仓库
type mockFileRepoUpdateError struct {
	files map[uint]*model.File
}

func newMockFileRepoUpdateError() *mockFileRepoUpdateError {
	return &mockFileRepoUpdateError{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "test.txt", Size: 1024},
		},
	}
}
func (m *mockFileRepoUpdateError) Create(file *model.File) error { return nil }
func (m *mockFileRepoUpdateError) FindByUserID(userID uint) ([]model.File, error) {
	return []model.File{*m.files[1]}, nil
}
func (m *mockFileRepoUpdateError) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, errors.New("not found")
}
func (m *mockFileRepoUpdateError) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return m.FindByUserID(userID)
}
func (m *mockFileRepoUpdateError) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return nil, errors.New("not found")
}
func (m *mockFileRepoUpdateError) Delete(id uint) error { return nil }
func (m *mockFileRepoUpdateError) Update(file *model.File) error {
	return errors.New("update failed")
}

// TestFileHandler_UpdateFileDescription_Error 测试更新描述失败
func TestFileHandler_UpdateFileDescription_Error(t *testing.T) {
	repo := newMockFileRepoUpdateError()
	svc := service.NewFileService(repo)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/description", handler.UpdateFileDescription)

	w := httptest.NewRecorder()
	body := `{"description":"new desc"}`
	req := httptest.NewRequest("PUT", "/files/1/description", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFileHandler_RenameFile_Error 测试重命名失败
func TestFileHandler_RenameFile_Error(t *testing.T) {
	repo := newMockFileRepoUpdateError()
	svc := service.NewFileService(repo)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id", handler.RenameFile)

	w := httptest.NewRecorder()
	body := `{"name":"new.txt"}`
	req := httptest.NewRequest("PUT", "/files/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFileHandler_MoveFile_Error 测试移动失败
func TestFileHandler_MoveFile_Error(t *testing.T) {
	repo := newMockFileRepoUpdateError()
	svc := service.NewFileService(repo)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/files/:id/move", handler.MoveFile)

	w := httptest.NewRecorder()
	body := `{"folder_id":1}`
	req := httptest.NewRequest("PUT", "/files/1/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// mockFileRepoCreateError 模拟Create返回错误的仓库
type mockFileRepoCreateError struct {
	files map[uint]*model.File
}

func newMockFileRepoCreateError() *mockFileRepoCreateError {
	return &mockFileRepoCreateError{
		files: map[uint]*model.File{},
	}
}
func (m *mockFileRepoCreateError) Create(file *model.File) error {
	return errors.New("db write failed")
}
func (m *mockFileRepoCreateError) FindByUserID(userID uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoCreateError) FindByID(id uint) (*model.File, error) {
	return nil, errors.New("not found")
}
func (m *mockFileRepoCreateError) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoCreateError) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return nil, errors.New("not found")
}
func (m *mockFileRepoCreateError) Delete(id uint) error { return nil }
func (m *mockFileRepoCreateError) Update(file *model.File) error { return nil }

// TestFileHandler_CreateFileMetadata_Error 测试创建元数据失败
func TestFileHandler_CreateFileMetadata_Error(t *testing.T) {
	repo := newMockFileRepoCreateError()
	svc := service.NewFileService(repo)
	handler := NewFileHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files", handler.CreateFileMetadata)

	w := httptest.NewRecorder()
	body := `{"name":"test.txt","size":100,"lanzou_file_id":"lz1","encryption_key":"k","encryption_nonce":"n"}`
	req := httptest.NewRequest("POST", "/files", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
type mockFileRepoNotFoundError struct{}

func (m *mockFileRepoNotFoundError) Create(file *model.File) error { return nil }
func (m *mockFileRepoNotFoundError) FindByUserID(userID uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoNotFoundError) FindByID(id uint) (*model.File, error) {
	return nil, errors.New("file not found")
}
func (m *mockFileRepoNotFoundError) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoNotFoundError) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return nil, errors.New("not found")
}
func (m *mockFileRepoNotFoundError) Delete(id uint) error { return errors.New("not found") }
func (m *mockFileRepoNotFoundError) Update(file *model.File) error { return errors.New("not found") }

// TestFileHandler_GetFile_NotFound 测试获取不存在的文件
func TestFileHandler_GetFile_NotFound(t *testing.T) {
	repo := &mockFileRepoNotFoundError{}
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// TestFileHandler_DeleteFile_AccessDenied 测试删除他人文件
func TestFileHandler_DeleteFile_AccessDenied_Path(t *testing.T) {
	repo := &mockFileRepoAccessDenied{}
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/files/:id", handler.DeleteFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/files/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// TestFileHandler_GetFile_Success 测试获取文件成功
func TestFileHandler_GetFile_Success_Path(t *testing.T) {
	repo := newMockFileRepoForHandlerTest()
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id", handler.GetFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TestFileHandler_CreateFileMetadata_Success 测试创建文件元数据成功
func TestFileHandler_CreateFileMetadata_Success_Path(t *testing.T) {
	repo := newMockFileRepoForHandlerTest()
	fileSvc := service.NewFileService(repo)
	handler := NewFileHandler(fileSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files", handler.CreateFileMetadata)

	w := httptest.NewRecorder()
	body := `{"name":"new.txt","size":100,"lanzou_file_id":"lz1","encryption_key":"k","encryption_nonce":"n"}`
	req := httptest.NewRequest("POST", "/files", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TestFileHandler_ListFiles_Success 测试获取文件列表成功
func TestFileHandler_ListFiles_Success_Path(t *testing.T) {
	repo := newMockFileRepoForHandlerTest()
	fileSvc := service.NewFileService(repo)
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

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// mockFileRepoForHandlerTest 简单mock
type mockFileRepoForHandlerTest struct {
	files map[uint]*model.File
}

func newMockFileRepoForHandlerTest() *mockFileRepoForHandlerTest {
	return &mockFileRepoForHandlerTest{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "test.txt", Size: 1024},
		},
	}
}
func (m *mockFileRepoForHandlerTest) Create(file *model.File) error {
	file.ID = uint(len(m.files)) + 1
	m.files[file.ID] = file
	return nil
}
func (m *mockFileRepoForHandlerTest) FindByUserID(userID uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}
func (m *mockFileRepoForHandlerTest) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, errors.New("not found")
}
func (m *mockFileRepoForHandlerTest) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return m.FindByUserID(userID)
}
func (m *mockFileRepoForHandlerTest) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockFileRepoForHandlerTest) Delete(id uint) error {
	delete(m.files, id)
	return nil
}
func (m *mockFileRepoForHandlerTest) Update(file *model.File) error {
	m.files[file.ID] = file
	return nil
}

// mockFileRepoAccessDenied 模拟访问他人文件的仓库
type mockFileRepoAccessDenied struct{}

func (m *mockFileRepoAccessDenied) Create(file *model.File) error { return nil }
func (m *mockFileRepoAccessDenied) FindByUserID(userID uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoAccessDenied) FindByID(id uint) (*model.File, error) {
	return &model.File{ID: 1, UserID: 2, Name: "other-file.txt", Size: 100}, nil
}
func (m *mockFileRepoAccessDenied) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return []model.File{}, nil
}
func (m *mockFileRepoAccessDenied) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return nil, errors.New("not found")
}
func (m *mockFileRepoAccessDenied) Delete(id uint) error { return nil }
func (m *mockFileRepoAccessDenied) Update(file *model.File) error { return nil }
