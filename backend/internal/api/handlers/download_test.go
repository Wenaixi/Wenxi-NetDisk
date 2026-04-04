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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/abc", nil)
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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/notanumber", nil)
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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/-1", nil)
	r.ServeHTTP(w, req)

	// strconv.ParseUint fails on negative numbers
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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/12.5", nil)
	r.ServeHTTP(w, req)

	// Float is not a valid uint
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
	r.GET("/download/:id", handler.GetDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/download/99999999999999999999999", nil)
	r.ServeHTTP(w, req)

	// Value too large for uint32 (ParseUint with 32 bits)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
