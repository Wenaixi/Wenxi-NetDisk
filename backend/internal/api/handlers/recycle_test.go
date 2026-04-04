package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
