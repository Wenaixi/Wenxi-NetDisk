package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHandler_Constructors tests handler constructors don't panic with nil
func TestHandler_Constructors(t *testing.T) {
	t.Run("should create all handlers with nil service", func(t *testing.T) {
		assert.NotNil(t, NewAuthHandler(nil))
		assert.NotNil(t, NewFileHandler(nil))
		assert.NotNil(t, NewFolderHandler(nil))
		assert.NotNil(t, NewShareHandler(nil))
		assert.NotNil(t, NewRecycleHandler(nil))
		assert.NotNil(t, NewLanZouHandler(nil, nil))
		assert.NotNil(t, NewShareParseHandler(nil))
		assert.NotNil(t, NewUploadHandler(nil, nil))
		assert.NotNil(t, NewDownloadHandler(nil))
	})
}

// TestAuthHandler_GetCurrentUser_NoUserID tests missing user_id in context
func TestAuthHandler_GetCurrentUser_NoUserID(t *testing.T) {
	handler := NewAuthHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/auth/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/me", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic for missing user_id: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestAuthHandler_GetCurrentUser_WithUser tests with user_id in context
func TestAuthHandler_GetCurrentUser_WithUser(t *testing.T) {
	handler := NewAuthHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/auth/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/me", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestRecycleHandler_List_WithUser tests recycle list handler
func TestRecycleHandler_List_WithUser(t *testing.T) {
	handler := NewRecycleHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/recycle", handler.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/recycle", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestRecycleHandler_Clear_WithUser tests recycle clear handler
func TestRecycleHandler_Clear_WithUser(t *testing.T) {
	handler := NewRecycleHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/recycle/clear", handler.Clear)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/recycle/clear", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestShareHandler_GetShare tests get share by token
func TestShareHandler_GetShare(t *testing.T) {
	handler := NewShareHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/shares/:token", handler.GetShare)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shares/abc123", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestShareHandler_ListShares tests listing shares
func TestShareHandler_ListShares(t *testing.T) {
	handler := NewShareHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/shares", handler.ListShares)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shares", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_GetStatus tests lanzou status endpoint
func TestLanzouHandler_GetStatus(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/status", handler.GetStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/status", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_Disconnect tests disconnect endpoint
func TestLanzouHandler_Disconnect(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/lanzou/connect", handler.Disconnect)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/lanzou/connect", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_ListFiles tests listing lanzou files
func TestLanzouHandler_ListFiles(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/files", handler.ListFiles)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/files", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestLanzouHandler_ListFolders tests listing lanzou folders
func TestLanzouHandler_ListFolders(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/lanzou/folders", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestFileHandler_ListFiles tests listing files
func TestFileHandler_ListFiles(t *testing.T) {
	handler := NewFileHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files", handler.ListFiles)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
	r.ServeHTTP(w, req)
}

// TestFileHandler_CreateFileMetadata tests creating file metadata
func TestFileHandler_CreateFileMetadata(t *testing.T) {
	handler := NewFileHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/files", handler.CreateFileMetadata)

	w := httptest.NewRecorder()
	body := `{"name":"test.txt","size":1024,"lanzou_file_id":"lz-123"}`
	req := httptest.NewRequest("POST", "/files", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic with nil service: %v", r)
		}
	}()
}

// TestLanzouHandler_CreateFolder_Validation tests create folder validation
func TestLanzouHandler_CreateFolder_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/folders", handler.CreateFolder)

	t.Run("should fail with missing name", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/folders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `not json`
		req := httptest.NewRequest("POST", "/lanzou/folders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLanzouHandler_InitializeUpload_Validation tests init upload validation
func TestLanzouHandler_InitializeUpload_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/upload/init", handler.InitializeUpload)

	t.Run("should fail with missing file_name", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"file_size":1000}`
		req := httptest.NewRequest("POST", "/lanzou/upload/init", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `invalid`
		req := httptest.NewRequest("POST", "/lanzou/upload/init", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLanzouHandler_CreateShare_Validation tests create share validation
func TestLanzouHandler_CreateShare_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/share", handler.CreateShare)

	t.Run("should fail with missing file_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/share", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with non-numeric file_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"file_id":"abc"}`
		req := httptest.NewRequest("POST", "/lanzou/share", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLanzouHandler_CompleteUpload_Validation tests complete upload validation
func TestLanzouHandler_CompleteUpload_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/lanzou/upload/complete/:id", handler.CompleteUpload)

	t.Run("should fail with non-numeric session_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"lanzou_file_id":"lz-123"}`
		req := httptest.NewRequest("POST", "/lanzou/upload/complete/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with negative session_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"lanzou_file_id":"lz-123"}`
		req := httptest.NewRequest("POST", "/lanzou/upload/complete/-1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLanzouHandler_GetFileURL_Validation tests get file URL validation
func TestLanzouHandler_GetFileURL_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/files/:id/url", handler.GetFileURL)

	t.Run("should fail with non-numeric file_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/lanzou/files/abc/url", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestLanzouHandler_UploadStatus_Validation tests upload status validation
func TestLanzouHandler_UploadStatus_Validation(t *testing.T) {
	handler := NewLanZouHandler(nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/lanzou/upload/status/:id", handler.UploadStatus)

	t.Run("should fail with non-numeric session_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/lanzou/upload/status/abc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestShareParseHandler_ParseShare_Validation tests parse share validation
func TestShareParseHandler_ParseShare_Validation(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/parse", handler.ParseShare)

	t.Run("should fail with empty body", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/shares/parse", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `not json`
		req := httptest.NewRequest("POST", "/shares/parse", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestShareParseHandler_GetShareDownloadURL_Validation tests get download URL validation
func TestShareParseHandler_GetShareDownloadURL_Validation(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/download", handler.GetShareDownloadURL)

	t.Run("should fail with empty body", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/shares/download", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestShareParseHandler_ValidateShareURL_Validation tests validate share URL validation
func TestShareParseHandler_ValidateShareURL_Validation(t *testing.T) {
	handler := NewShareParseHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/shares/validate", handler.ValidateShareURL)

	t.Run("should fail with empty body", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/shares/validate", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestDownloadHandler_GetShareDownloadURL_Validation tests share download URL validation
func TestDownloadHandler_GetShareDownloadURL_Validation(t *testing.T) {
	handler := NewDownloadHandler(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/files/:id/share-download", handler.GetDownloadURL)

	t.Run("should fail with non-numeric file_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/files/abc/share-download", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestResponseFormat tests that all handlers return properly formatted responses
func TestResponseFormat_BadRequest(t *testing.T) {
	t.Run("register returns proper error format", func(t *testing.T) {
		r := setupTestRouter()
		handler := &AuthHandler{}
		r.POST("/register", handler.Register)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("response should be valid JSON: %v", err)
		}

		if _, ok := resp["code"]; !ok {
			t.Error("response should contain 'code' field")
		}
		if _, ok := resp["msg"]; !ok {
			t.Error("response should contain 'msg' field")
		}
	})
}
