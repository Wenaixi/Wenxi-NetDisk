package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFileHandler_GetFile_Validation 测试获取文件详情验证
func TestFileHandler_GetFile_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.GET("/files/:id", handler.GetFile)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/files/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFileHandler_DeleteFile_Validation 测试删除文件验证
func TestFileHandler_DeleteFile_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.DELETE("/files/:id", handler.DeleteFile)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/files/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFileHandler_MoveFile_Validation 测试移动文件验证
func TestFileHandler_MoveFile_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.PUT("/files/:id/move", handler.MoveFile)

		w := httptest.NewRecorder()
		body := `{"folder_id":2}`
		req := httptest.NewRequest("PUT", "/files/abc/move", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFileHandler_UpdateFileDescription_Validation 测试更新文件描述验证
func TestFileHandler_UpdateFileDescription_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.PUT("/files/:id/description", handler.UpdateFileDescription)

		w := httptest.NewRecorder()
		body := `{"description":"test"}`
		req := httptest.NewRequest("PUT", "/files/abc/description", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

// TestFileHandler_RenameFile_Validation 测试重命名文件验证
func TestFileHandler_RenameFile_Validation(t *testing.T) {
	t.Run("should fail with invalid file id", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.PUT("/files/:id", handler.RenameFile)

		w := httptest.NewRecorder()
		body := `{"name":"test.txt"}`
		req := httptest.NewRequest("PUT", "/files/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("should fail with missing name", func(t *testing.T) {
		r := setupTestRouterWithUser()
		handler := &FileHandler{}
		r.PUT("/files/1", handler.RenameFile)

		w := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest("PUT", "/files/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 for missing name, got %d", w.Code)
		}
	})
}

// TestFileHandler_CreateFileMetadata_Validation 测试创建文件元数据验证
func TestFileHandler_CreateFileMetadata_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{
			"missing required fields",
			`{}`,
			http.StatusBadRequest,
		},
		{
			"missing name",
			`{"size":100,"lanzou_file_id":"abc","encryption_key":"key","encryption_nonce":"nonce"}`,
			http.StatusBadRequest,
		},
		{
			"missing size",
			`{"name":"test.txt","lanzou_file_id":"abc","encryption_key":"key","encryption_nonce":"nonce"}`,
			http.StatusBadRequest,
		},
		{
			"invalid JSON",
			`invalid json`,
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouterWithUser()
			handler := &FileHandler{}
			r.POST("/files", handler.CreateFileMetadata)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/files", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectCode, w.Code, w.Body.String())
			}
		})
	}
}
