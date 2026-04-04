package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDownloadHandler_GetDownloadURL_InvalidID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/abc/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadHandler_GetDownloadURL_NonNumericID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/notanumber/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadHandler_GetDownloadURL_EmptyID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files//download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadHandler_GetDownloadURL_NegativeID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/-1/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadHandler_GetDownloadURL_FloatingPointID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/12.5/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadHandler_GetDownloadURL_VeryLongID(t *testing.T) {
	svc := service.NewDownloadService(nil, nil)
	handler := NewDownloadHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/download", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/99999999999999999999999/download", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
