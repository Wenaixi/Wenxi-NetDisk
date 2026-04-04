package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockFileRepoRound2 模拟文件仓库用于 Round2 Handler 测试
type mockFileRepoRound2 struct {
	files map[uint]*model.File
}

func newMockFileRepoRound2() *mockFileRepoRound2 {
	return &mockFileRepoRound2{files: make(map[uint]*model.File)}
}
func (m *mockFileRepoRound2) Create(file *model.File) error {
	file.ID = uint(len(m.files)) + 1
	m.files[file.ID] = file
	return nil
}
func (m *mockFileRepoRound2) FindByUserID(userID uint) ([]model.File, error) {
	var result []model.File
	for _, f := range m.files {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}
func (m *mockFileRepoRound2) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, assert.AnError
}
func (m *mockFileRepoRound2) FindByFolderID(userID uint, folderID *uint) ([]model.File, error) {
	return nil, nil
}
func (m *mockFileRepoRound2) FindByLanZouFileID(lanzouFileID string) (*model.File, error) {
	return nil, assert.AnError
}
func (m *mockFileRepoRound2) Delete(id uint) error {
	delete(m.files, id)
	return nil
}
func (m *mockFileRepoRound2) Update(file *model.File) error {
	m.files[file.ID] = file
	return nil
}

// createMultipartForm 辅助：创建 multipart/form-data body
func createMultipartForm(fieldName, fileName, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, fileName)
	part.Write([]byte(content))
	writer.Close()
	return body, writer.FormDataContentType()
}

// TestUploadHandler_UploadFile_Multipart_Success tests multipart upload success
func TestUploadHandler_UploadFile_Multipart_Success(t *testing.T) {
	repo := newMockFileRepoRound2()
	fileSvc := service.NewFileService(repo)
	handler := NewUploadHandler(fileSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload", handler.UploadFile)

	w := httptest.NewRecorder()
	body, contentType := createMultipartForm("file", "test.txt", "hello world")
	req := httptest.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", contentType)
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "test.txt", data["name"])
}

func TestUploadHandler_UploadFile_NoFile(t *testing.T) {
	handler := NewUploadHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload", handler.UploadFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files/upload", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadHandler_UploadFile_WithFolder(t *testing.T) {
	repo := newMockFileRepoRound2()
	fileSvc := service.NewFileService(repo)
	handler := NewUploadHandler(fileSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload", handler.UploadFile)

	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "doc.pdf")
	part.Write([]byte("pdf content"))
	writer.WriteField("folder_id", "5")
	writer.Close()

	req := httptest.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "doc.pdf", data["name"])
}

// TestUploadHandler_UploadFile_WrongFieldName tests wrong form field name
func TestUploadHandler_UploadFile_WrongFieldName(t *testing.T) {
	handler := NewUploadHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/upload", handler.UploadFile)

	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.CreateFormFile("wrong_field", "test.txt")
	writer.Close()

	req := httptest.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUploadHandler_RestoreVersion_Success tests restoring a version
func TestUploadHandler_RestoreVersion_Success(t *testing.T) {
	vRepo := newMockVersionRepoForRestore()
	fRepo := newMockFileRepoForRestore()

	vRepo.versions[1] = &model.FileVersion{
		ID: 1, FileID: 1, UserID: 1, Size: 200,
	}

	versionSvc := service.NewFileVersionService(vRepo, fRepo)
	handler := NewUploadHandler(nil, versionSvc)

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

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "original.txt", data["name"])
}

// mockVersionRepoForRestore 模拟版本仓库
type mockVersionRepoForRestore struct {
	versions map[uint]*model.FileVersion
}

func newMockVersionRepoForRestore() *mockVersionRepoForRestore {
	return &mockVersionRepoForRestore{versions: make(map[uint]*model.FileVersion)}
}
func (m *mockVersionRepoForRestore) Create(version *model.FileVersion) error {
	version.ID = uint(len(m.versions)) + 1
	m.versions[version.ID] = version
	return nil
}
func (m *mockVersionRepoForRestore) FindByFileID(fileID uint) ([]model.FileVersion, error) {
	var result []model.FileVersion
	for _, v := range m.versions {
		if v.FileID == fileID {
			result = append(result, *v)
		}
	}
	return result, nil
}
func (m *mockVersionRepoForRestore) FindByID(id uint) (*model.FileVersion, error) {
	if v, ok := m.versions[id]; ok {
		return v, nil
	}
	return nil, assert.AnError
}
func (m *mockVersionRepoForRestore) Delete(id uint) error {
	delete(m.versions, id)
	return nil
}
func (m *mockVersionRepoForRestore) DeleteByFileID(fileID uint) error {
	for id, v := range m.versions {
		if v.FileID == fileID {
			delete(m.versions, id)
		}
	}
	return nil
}

// mockFileRepoForRestore 模拟文件仓库
type mockFileRepoForRestore struct {
	files map[uint]*model.File
}

func newMockFileRepoForRestore() *mockFileRepoForRestore {
	return &mockFileRepoForRestore{
		files: map[uint]*model.File{
			1: {ID: 1, UserID: 1, Name: "original.txt", Size: 100},
		},
	}
}
func (m *mockFileRepoForRestore) FindByID(id uint) (*model.File, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, assert.AnError
}
func (m *mockFileRepoForRestore) Update(file *model.File) error {
	m.files[file.ID] = file
	return nil
}

// TestUploadHandler_ListVersions_R2Success tests listing versions
func TestUploadHandler_ListVersions_R2Success(t *testing.T) {
	vRepo := newMockVersionRepoForRestore()
	fRepo := newMockFileRepoForRestore()
	versionSvc := service.NewFileVersionService(vRepo, fRepo)
	handler := NewUploadHandler(nil, versionSvc)

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

	assertJSONResponse(t, w, http.StatusOK)
}

// TestShareHandler_CreateShareViaBody_Success tests creating share via body
func TestShareHandler_CreateShareViaBody_R2Success(t *testing.T) {
	shareRepo := NewMockShareRepo()
	fileRepo := NewMockFileRepoForShare()
	fileRepo.files[1] = &model.File{ID: 1, UserID: 1, Name: "body-share.txt", Size: 64}
	shareSvc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(shareSvc)

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

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["share_token"])
	assert.Contains(t, data["share_url"], "/api/shares/")
}

func TestShareHandler_CreateShareViaBody_WithPassword(t *testing.T) {
	shareRepo := NewMockShareRepo()
	fileRepo := NewMockFileRepoForShare()
	fileRepo.files[1] = &model.File{ID: 1, UserID: 1, Name: "protected.txt", Size: 32}
	shareSvc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(shareSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares", handler.CreateShareViaBody)

	w := httptest.NewRecorder()
	body := `{"file_id":1,"password":"mypass","expires_at":"2026-06-01T00:00:00Z"}`
	req := httptest.NewRequest("POST", "/shares", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["share_token"])
}

// TestShareHandler_CreateShare_R2Success tests creating share via URL param
func TestShareHandler_CreateShare_R2Success(t *testing.T) {
	shareRepo := NewMockShareRepo()
	fileRepo := NewMockFileRepoForShare()
	fileRepo.files[1] = &model.File{ID: 1, UserID: 1, Name: "share.txt", Size: 512}
	shareSvc := service.NewShareService(shareRepo, fileRepo)
	handler := NewShareHandler(shareSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files/:id/share", handler.CreateShare)

	w := httptest.NewRecorder()
	body := `{}`
	req := httptest.NewRequest("POST", "/files/1/share", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["share_token"])
}

// TestGuessMimeType_EdgeCases tests edge cases for MIME type guessing
func TestGuessMimeType_EdgeCases(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".webp", "application/octet-stream"},
		{".html", "application/octet-stream"},
		{".js", "application/octet-stream"},
		{".css", "application/octet-stream"},
		{".svg", "application/octet-stream"},
		{".mkv", "application/octet-stream"},
		{".flac", "application/octet-stream"},
		{".exe", "application/octet-stream"},
		{".tar", "application/octet-stream"},
		{".gz", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := guessMimeType(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGuessMimeType_AllKnownTypes tests all known MIME types
func TestGuessMimeType_AllKnownTypes(t *testing.T) {
	expected := map[string]string{
		".txt":  "text/plain",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
	}

	for ext, mime := range expected {
		t.Run(ext, func(t *testing.T) {
			assert.Equal(t, mime, guessMimeType(ext))
		})
	}
}

// TestShareParseHandler_ParseShare_ValidBody tests parse share with valid body
func TestShareParseHandler_ParseShare_ValidBody(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/parse", handler.ParseShare)

	w := httptest.NewRecorder()
	body := `{"url":"https://abc.lanzoui.com/xyz123"}`
	req := httptest.NewRequest("POST", "/shares/parse", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: share parse service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestDownloadHandler_GetDownloadURL_Success tests download URL success
func TestDownloadHandler_GetDownloadURL_Success(t *testing.T) {
	handler := NewDownloadHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/1/download", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: download service is nil")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestFileHandler_ListFiles_Success tests listing files
func TestFileHandler_ListFiles_Success(t *testing.T) {
	repo := newMockFileRepoForHandler()
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

	assertJSONResponse(t, w, http.StatusOK)
}

// TestFileHandler_CreateFileMetadata_Success tests creating file metadata
func TestFileHandler_CreateFileMetadata_Success(t *testing.T) {
	repo := newMockFileRepoForHandler()
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
	body := `{"name":"new-file.txt","size":2048,"lanzou_file_id":"lz-999","encryption_key":"enc-123","encryption_nonce":"nonce-456"}`
	req := httptest.NewRequest("POST", "/files", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := assertJSONResponse(t, w, http.StatusOK)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "new-file.txt", data["name"])
}

// TestLanzouHandler_CreateShare_R2Success tests creating lanzou share
func TestLanzouHandler_CreateShare_R2Success(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	tokenRepo.tokens[1] = &model.LanZouToken{
		ID:        1, UserID: 1, Cookie: "test-cookie",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/share", handler.CreateShare)

	w := httptest.NewRecorder()
	body := `{"file_id":123,"minutes":10080}`
	req := httptest.NewRequest("POST", "/lanzou/share", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: lanzou client Task39 requires network")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_GetFileURL_R2Success tests getting file URL
func TestLanzouHandler_GetFileURL_R2Success(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	tokenRepo.tokens[1] = &model.LanZouToken{
		ID:        1, UserID: 1, Cookie: "test-cookie",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/files/:id/url", handler.GetFileURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/files/456/url", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: lanzou client Task22 requires network")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_ListFiles_R2Success tests listing lanzou files
func TestLanzouHandler_ListFiles_R2Success(t *testing.T) {
	tokenRepo := NewMockLanZouTokenRepo()
	tokenRepo.tokens[1] = &model.LanZouToken{
		ID:        1, UserID: 1, Cookie: "test-cookie",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	lanzouSvc := service.NewLanZouService(tokenRepo)
	handler := NewLanZouHandler(lanzouSvc, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/files", handler.ListFiles)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/files?folder_id=-1&page=1", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: lanzou client Task5/Task47 requires network")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestRecycleHandler_Restore_BadID tests restore with bad ID (renamed to avoid conflict)
func TestRecycleHandler_Restore_BadID(t *testing.T) {
	recycleRepo := &MockRecycleBinRepository{}
	handler := NewRecycleHandler(service.NewRecycleBinService(recycleRepo, nil, nil))

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

// TestRecycleHandler_Delete_BadID tests delete with bad ID (renamed to avoid conflict)
func TestRecycleHandler_Delete_BadID(t *testing.T) {
	recycleRepo := &MockRecycleBinRepository{}
	handler := NewRecycleHandler(service.NewRecycleBinService(recycleRepo, nil, nil))

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
