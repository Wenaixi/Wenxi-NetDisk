package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLanZouHandler_Connect_Validation 测试蓝奏云连接验证
func TestLanZouHandler_Connect_Validation(t *testing.T) {
	t.Run("should fail with missing cookie", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/connect", handler.Connect)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/connect", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/connect", handler.Connect)

		w := httptest.NewRecorder()
		body := `invalid`
		req := httptest.NewRequest("POST", "/lanzou/connect", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_CreateFolder_Validation 测试创建蓝奏云文件夹验证
func TestLanZouHandler_CreateFolder_Validation(t *testing.T) {
	t.Run("should fail with missing name", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/folders", handler.CreateFolder)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/folders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_InitializeUpload_Validation 测试初始化上传验证
func TestLanZouHandler_InitializeUpload_Validation(t *testing.T) {
	t.Run("should fail with empty body", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/upload/init", handler.InitializeUpload)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/upload/init", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// file_name and file_size are required, empty body should return 400
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should fail with missing file_name", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/upload/init", handler.InitializeUpload)

		w := httptest.NewRecorder()
		body := `{"file_size":1000}`
		req := httptest.NewRequest("POST", "/lanzou/upload/init", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_CompleteUpload_Validation 测试完成上传验证
func TestLanZouHandler_CompleteUpload_Validation(t *testing.T) {
	t.Run("should fail with invalid session id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/upload/complete/abc", handler.CompleteUpload)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/upload/complete/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_UploadStatus_Validation 测试上传状态验证
func TestLanZouHandler_UploadStatus_Validation(t *testing.T) {
	t.Run("should fail with invalid session id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.GET("/lanzou/upload/status/abc", handler.UploadStatus)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/lanzou/upload/status/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_CreateShare_Validation 测试蓝奏云分享验证
func TestLanZouHandler_CreateShare_Validation(t *testing.T) {
	t.Run("should fail with missing file_id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.POST("/lanzou/share", handler.CreateShare)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("POST", "/lanzou/share", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestLanZouHandler_GetFileURL_Validation 测试获取文件URL验证
func TestLanZouHandler_GetFileURL_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &LanZouHandler{}
		r.GET("/lanzou/files/abc/url", handler.GetFileURL)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/lanzou/files/abc/url", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
