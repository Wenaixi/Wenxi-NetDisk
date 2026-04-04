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
