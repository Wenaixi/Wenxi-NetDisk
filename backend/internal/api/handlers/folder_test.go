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

func TestFolderHandler_GetFolder_InvalidID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_GetFolder_NonNumericID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_DeleteFolder_InvalidID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_DeleteFolder_NonNumericID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/notanumber", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolder_InvalidID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/abc", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolder_InvalidJSON(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/123", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolder_EmptyBody(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/123", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_MoveFolder_InvalidID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/abc/move", strings.NewReader(`{"parent_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_MoveFolder_InvalidJSON(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/123/move", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_CreateFolder_EmptyBody(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_CreateFolder_InvalidJSON(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_ListFolders_InvalidParentID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders?parent_id=abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolderDescription_InvalidID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/description", handler.UpdateFolderDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/abc/description", strings.NewReader(`{"description":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolderDescription_InvalidJSON(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/description", handler.UpdateFolderDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/123/description", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_GetFolder_FloatingPointID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_DeleteFolder_FloatingPointID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/12.5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolder_NegativeID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/-1", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_MoveFolder_NegativeID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/-1/move", strings.NewReader(`{"parent_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_GetFolder_NegativeID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_DeleteFolder_NegativeID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFolderHandler_UpdateFolderDescription_NegativeID(t *testing.T) {
	svc := service.NewFolderService(nil)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/description", handler.UpdateFolderDescription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/folders/-1/description", strings.NewReader(`{"description":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
