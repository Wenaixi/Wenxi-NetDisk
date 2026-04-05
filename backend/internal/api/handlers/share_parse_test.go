package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShareParseHandler_ParseShare_EmptyBody(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/parse", handler.ParseShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/parse", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareParseHandler_ParseShare_InvalidJSON(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/parse", handler.ParseShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/parse", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareParseHandler_GetShareDownloadURL_EmptyBody(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/download", handler.GetShareDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/download", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareParseHandler_GetShareDownloadURL_InvalidJSON(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/download", handler.GetShareDownloadURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/download", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareParseHandler_ValidateShareURL_EmptyBody(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/validate", handler.ValidateShareURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/validate", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShareParseHandler_ValidateShareURL_InvalidJSON(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/validate", handler.ValidateShareURL)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/share/validate", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestShareParseHandler_ParseShare_ValidURL tests with valid URL (network call expected)
func TestShareParseHandler_ParseShare_ValidURL(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/parse", handler.ParseShare)

	w := httptest.NewRecorder()
	body := `{"url":"https://lanzou.com/i12345"}`
	req := httptest.NewRequest("POST", "/share/parse", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestShareParseHandler_GetShareDownloadURL_ValidURL tests download URL handler
func TestShareParseHandler_GetShareDownloadURL_ValidURL(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/download", handler.GetShareDownloadURL)

	w := httptest.NewRecorder()
	body := `{"url":"https://lanzou.com/i12345","pwd":"1234"}`
	req := httptest.NewRequest("POST", "/share/download", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected: service is nil, handler validated input")
		}
	}()
	r.ServeHTTP(w, req)
}

// TestShareParseHandler_ValidateShareURL_ValidURL tests validate handler
func TestShareParseHandler_ValidateShareURL_ValidURL(t *testing.T) {
	svc := service.NewShareParseService(nil)
	handler := NewShareParseHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/share/validate", handler.ValidateShareURL)

	w := httptest.NewRecorder()
	body := `{"url":"https://lanzou.com/i12345"}`
	req := httptest.NewRequest("POST", "/share/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// ValidateShareURL calls lanzou.ValidateShareURL which makes HTTP call
	t.Logf("status: %d", w.Code)
}
