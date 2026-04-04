package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
