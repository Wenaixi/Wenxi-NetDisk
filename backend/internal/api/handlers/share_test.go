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
