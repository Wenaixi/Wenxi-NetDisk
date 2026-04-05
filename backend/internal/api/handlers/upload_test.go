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

func TestUploadHandler_ListVersions_InvalidID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/abc/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_ListVersions_NonNumericID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/notanumber/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_ListVersions_FloatingPointID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/12.5/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_ListVersions_NegativeID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/-1/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_InvalidFileID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/abc/versions/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_InvalidVersionID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/123/versions/abc/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_NonNumericFileID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/notanumber/versions/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_NonNumericVersionID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/123/versions/notanumber/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_FloatingPointFileID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/12.5/versions/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_RestoreVersion_FloatingPointVersionID(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/123/versions/4.5/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_GetUploadURL_InvalidJSON(t *testing.T) {
	fileSvc := service.NewFileService(nil)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload-url", handler.GetUploadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/upload-url", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// mockFileVersionRepoForHandler implements FileVersionRepository
type mockFileVersionRepoForHandler struct {
	versions map[uint]*model.FileVersion
}

func newMockFileVersionRepoForHandler() *mockFileVersionRepoForHandler {
	return &mockFileVersionRepoForHandler{
		versions: map[uint]*model.FileVersion{
			1: {ID: 1, FileID: 1, UserID: 1, LanZouFileID: "lz1", Size: 1024, EncryptionKey: "k1", EncryptionNonce: "n1"},
		},
	}
}

func (m *mockFileVersionRepoForHandler) Create(version *model.FileVersion) error {
	version.ID = uint(len(m.versions)) + 1
	m.versions[version.ID] = version
	return nil
}
func (m *mockFileVersionRepoForHandler) FindByFileID(fileID uint) ([]model.FileVersion, error) {
	var result []model.FileVersion
	for _, v := range m.versions {
		if v.FileID == fileID {
			result = append(result, *v)
		}
	}
	return result, nil
}
func (m *mockFileVersionRepoForHandler) FindByID(id uint) (*model.FileVersion, error) {
	if v, ok := m.versions[id]; ok {
		return v, nil
	}
	return nil, errors.New("version not found")
}
func (m *mockFileVersionRepoForHandler) Delete(id uint) error {
	delete(m.versions, id)
	return nil
}
func (m *mockFileVersionRepoForHandler) DeleteByFileID(fileID uint) error {
	for id, v := range m.versions {
		if v.FileID == fileID {
			delete(m.versions, id)
		}
	}
	return nil
}

// mockVersionFileRepoForHandler implements VersionFileRepository
type mockVersionFileRepoForHandler struct {
	files map[uint]*model.File
}

func newMockVersionFileRepoForHandler() *mockVersionFileRepoForHandler {
	return &mockVersionFileRepoForHandler{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "test.txt", Size: 512, LanZouFileID: "lz_old"},
		},
	}
}

func (m *mockVersionFileRepoForHandler) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, errors.New("file not found")
}
func (m *mockVersionFileRepoForHandler) Update(file *model.File) error {
	m.files[file.ID] = file
	return nil
}

// mockFileRepoForUpload implements FileRepository for upload tests
type mockFileRepoForUpload struct {
	files map[uint]*model.File
}

func newMockFileRepoForUpload() *mockFileRepoForUpload {
	return &mockFileRepoForUpload{
		files: make(map[uint]*model.File),
	}
}

func (m *mockFileRepoForUpload) Create(file *model.File) error {
	file.ID = uint(len(m.files)) + 1
	m.files[file.ID] = file
	return nil
}
func (m *mockFileRepoForUpload) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, errors.New("file not found")
}
func (m *mockFileRepoForUpload) FindByUserID(userID uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}
func (m *mockFileRepoForUpload) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return m.FindByUserID(userID)
}
func (m *mockFileRepoForUpload) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	for _, f := range m.files {
		if f.LanZouFileID == lanzouFileID {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockFileRepoForUpload) Delete(id uint) error {
	delete(m.files, id)
	return nil
}
func (m *mockFileRepoForUpload) Update(file *model.File) error {
	m.files[file.ID] = file
	return nil
}

// TestUploadHandler_ListVersions_Success_Path 测试获取版本列表成功
func TestUploadHandler_ListVersions_Success_Path(t *testing.T) {
	versionRepo := newMockFileVersionRepoForHandler()
	fileRepo := newMockVersionFileRepoForHandler()
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)
	fileSvc := service.NewFileService(nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/1/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUploadHandler_ListVersions_FileNotFound 测试文件不存在
func TestUploadHandler_ListVersions_FileNotFound(t *testing.T) {
	versionRepo := newMockFileVersionRepoForHandler()
	fileRepo := newMockVersionFileRepoForHandler()
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)
	fileSvc := service.NewFileService(nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/versions", handler.ListVersions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/999/versions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestUploadHandler_RestoreVersion_Success_Path 测试恢复版本成功
func TestUploadHandler_RestoreVersion_Success_Path(t *testing.T) {
	versionRepo := newMockFileVersionRepoForHandler()
	fileRepo := newMockVersionFileRepoForHandler()
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)
	fileSvc := service.NewFileService(nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/1/versions/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUploadHandler_RestoreVersion_FileNotFound 测试文件不存在
func TestUploadHandler_RestoreVersion_FileNotFound(t *testing.T) {
	versionRepo := newMockFileVersionRepoForHandler()
	fileRepo := newMockVersionFileRepoForHandler()
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)
	fileSvc := service.NewFileService(nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/999/versions/1/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestUploadHandler_RestoreVersion_VersionNotFound 测试版本不存在
func TestUploadHandler_RestoreVersion_VersionNotFound(t *testing.T) {
	versionRepo := newMockFileVersionRepoForHandler()
	fileRepo := newMockVersionFileRepoForHandler()
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)
	fileSvc := service.NewFileService(nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/versions/:version_id/restore", handler.RestoreVersion)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/1/versions/999/restore", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestUploadHandler_GetUploadURL_Success_Path 测试获取上传URL成功
func TestUploadHandler_GetUploadURL_Success_Path(t *testing.T) {
	repo := newMockFileRepoForUpload()
	fileSvc := service.NewFileService(repo)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload-url", handler.GetUploadURL)

	w := httptest.NewRecorder()
	body := `{"file_name":"test.txt","size":100}`
	req := httptest.NewRequest("POST", "/files/upload-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUploadHandler_GetUploadURL_MissingFields 测试缺少必填字段
func TestUploadHandler_GetUploadURL_MissingFields(t *testing.T) {
	repo := newMockFileRepoForUpload()
	fileSvc := service.NewFileService(repo)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload-url", handler.GetUploadURL)

	w := httptest.NewRecorder()
	body := `{"file_name":"test.txt"}`
	req := httptest.NewRequest("POST", "/files/upload-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUploadHandler_GetUploadURL_WithFolderID 测试带文件夹ID
func TestUploadHandler_GetUploadURL_WithFolderID_Path(t *testing.T) {
	repo := newMockFileRepoForUpload()
	fileSvc := service.NewFileService(repo)
	versionSvc := service.NewFileVersionService(nil, nil)
	handler := NewUploadHandler(fileSvc, versionSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload-url", handler.GetUploadURL)

	w := httptest.NewRecorder()
	body := `{"file_name":"test.txt","size":100,"folder_id":1}`
	req := httptest.NewRequest("POST", "/files/upload-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
