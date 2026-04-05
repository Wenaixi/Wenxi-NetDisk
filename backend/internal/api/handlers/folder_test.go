package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
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

// mockFolderRepoForHandlerTest 简单mock文件夹仓库
type mockFolderRepoForHandlerTest struct {
	folders map[uint]*model.Folder
}

func newMockFolderRepoForHandlerTest() *mockFolderRepoForHandlerTest {
	return &mockFolderRepoForHandlerTest{
		folders: map[uint]*model.Folder{
			1: {ID: 1, UserID: 1, ParentID: nil, Name: "root-folder"},
			2: {ID: 2, UserID: 1, ParentID: uintPtr(1), Name: "sub-folder"},
			3: {ID: 3, UserID: 2, ParentID: nil, Name: "other-folder"},
		},
	}
}

func uintPtr(v uint) *uint { return &v }

func (m *mockFolderRepoForHandlerTest) Create(folder *model.Folder) error {
	folder.ID = uint(len(m.folders)) + 1
	m.folders[folder.ID] = folder
	return nil
}
func (m *mockFolderRepoForHandlerTest) FindByID(id uint) (*model.Folder, error) {
	if f, ok := m.folders[id]; ok {
		return f, nil
	}
	return nil, errors.New("folder not found")
}
func (m *mockFolderRepoForHandlerTest) FindByUserID(userID uint) ([]model.Folder, error) {
	var result []model.Folder
	for _, f := range m.folders {
		if f.UserID == userID {
			result = append(result, *f)
		}
	}
	return result, nil
}
func (m *mockFolderRepoForHandlerTest) FindByParentID(userID uint, parentID *uint) ([]model.Folder, error) {
	var result []model.Folder
	for _, f := range m.folders {
		if f.UserID == userID {
			if parentID == nil && f.ParentID == nil {
				result = append(result, *f)
			} else if parentID != nil && f.ParentID != nil && *f.ParentID == *parentID {
				result = append(result, *f)
			}
		}
	}
	return result, nil
}
func (m *mockFolderRepoForHandlerTest) Update(folder *model.Folder) error {
	m.folders[folder.ID] = folder
	return nil
}
func (m *mockFolderRepoForHandlerTest) Delete(id uint) error {
	delete(m.folders, id)
	return nil
}
func (m *mockFolderRepoForHandlerTest) DeleteByUserID(userID uint) error {
	for id, f := range m.folders {
		if f.UserID == userID {
			delete(m.folders, id)
		}
	}
	return nil
}
func (m *mockFolderRepoForHandlerTest) DeleteByParentID(parentID uint) error {
	for id, f := range m.folders {
		if f.ParentID != nil && *f.ParentID == parentID {
			delete(m.folders, id)
		}
	}
	return nil
}

// TestFolderHandler_CreateFolder_Success 测试创建文件夹成功
func TestFolderHandler_CreateFolder_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"new-folder"}`
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_CreateFolder_WithParent 测试创建子文件夹
func TestFolderHandler_CreateFolder_WithParent(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"sub-folder","parent_id":1}`
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_CreateFolder_ParentNotFound 测试父文件夹不存在
func TestFolderHandler_CreateFolder_ParentNotFound(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.POST("/folders", handler.CreateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"sub-folder","parent_id":999}`
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFolderHandler_ListFolders_Success 测试获取文件夹列表成功
func TestFolderHandler_ListFolders_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_ListFolders_WithParentID 测试带parent_id的列表查询
func TestFolderHandler_ListFolders_WithParentID_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders?parent_id=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_GetFolder_Success 测试获取文件夹详情成功
func TestFolderHandler_GetFolder_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_GetFolder_NotFound 测试获取不存在的文件夹
func TestFolderHandler_GetFolder_NotFound_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestFolderHandler_GetFolder_AccessDenied 测试获取他人文件夹
func TestFolderHandler_GetFolder_AccessDenied_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders/:id", handler.GetFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders/3", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestFolderHandler_DeleteFolder_Success 测试删除文件夹成功
func TestFolderHandler_DeleteFolder_Success_Path_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_DeleteFolder_NotFound 测试删除不存在的文件夹
func TestFolderHandler_DeleteFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.DELETE("/folders/:id", handler.DeleteFolder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/folders/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFolderHandler_UpdateFolder_Success 测试更新文件夹成功
func TestFolderHandler_UpdateFolder_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"renamed-folder"}`
	req := httptest.NewRequest("PUT", "/folders/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_UpdateFolder_Description 测试更新文件夹描述
func TestFolderHandler_UpdateFolder_Description_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	body := `{"description":"my folder"}`
	req := httptest.NewRequest("PUT", "/folders/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_UpdateFolder_NotFound 测试更新不存在的文件夹
func TestFolderHandler_UpdateFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id", handler.UpdateFolder)

	w := httptest.NewRecorder()
	body := `{"name":"renamed"}`
	req := httptest.NewRequest("PUT", "/folders/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFolderHandler_MoveFolder_Success 测试移动文件夹成功
func TestFolderHandler_MoveFolder_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	body := `{"parent_id":1}`
	req := httptest.NewRequest("PUT", "/folders/2/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_MoveFolder_NoParent 测试移动到根目录
func TestFolderHandler_MoveFolder_NoParent(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	body := `{}`
	req := httptest.NewRequest("PUT", "/folders/2/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_MoveFolder_NotFound 测试移动不存在的文件夹
func TestFolderHandler_MoveFolder_NotFound(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	body := `{"parent_id":1}`
	req := httptest.NewRequest("PUT", "/folders/999/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFolderHandler_MoveFolder_Self 测试移动到自身
func TestFolderHandler_MoveFolder_Self(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/move", handler.MoveFolder)

	w := httptest.NewRecorder()
	body := `{"parent_id":1}`
	req := httptest.NewRequest("PUT", "/folders/1/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestFolderHandler_UpdateFolderDescription_Success 测试更新描述成功
func TestFolderHandler_UpdateFolderDescription_Success_Path(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/description", handler.UpdateFolderDescription)

	w := httptest.NewRecorder()
	body := `{"description":"updated description"}`
	req := httptest.NewRequest("PUT", "/folders/1/description", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFolderHandler_UpdateFolderDescription_NotFound 测试更新不存在的文件夹描述
func TestFolderHandler_UpdateFolderDescription_NotFound(t *testing.T) {
	repo := newMockFolderRepoForHandlerTest()
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.PUT("/folders/:id/description", handler.UpdateFolderDescription)

	w := httptest.NewRecorder()
	body := `{"description":"test"}`
	req := httptest.NewRequest("PUT", "/folders/999/description", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// mockFolderRepoNotFoundError 模拟FindByID返回错误的仓库
type mockFolderRepoNotFoundError struct{}

func (m *mockFolderRepoNotFoundError) Create(folder *model.Folder) error { return nil }
func (m *mockFolderRepoNotFoundError) FindByID(id uint) (*model.Folder, error) {
	return nil, errors.New("folder not found")
}
func (m *mockFolderRepoNotFoundError) FindByUserID(userID uint) ([]model.Folder, error) {
	return []model.Folder{}, nil
}
func (m *mockFolderRepoNotFoundError) FindByParentID(userID uint, parentID *uint) ([]model.Folder, error) {
	return []model.Folder{}, nil
}
func (m *mockFolderRepoNotFoundError) Update(folder *model.Folder) error { return errors.New("update failed") }
func (m *mockFolderRepoNotFoundError) Delete(id uint) error { return errors.New("delete failed") }
func (m *mockFolderRepoNotFoundError) DeleteByUserID(userID uint) error { return nil }
func (m *mockFolderRepoNotFoundError) DeleteByParentID(parentID uint) error { return nil }

// TestFolderHandler_ListFolders_Error 测试列表查询错误
func TestFolderHandler_ListFolders_Error(t *testing.T) {
	repo := &mockFolderRepoNotFoundError{}
	svc := service.NewFolderService(repo)
	handler := NewFolderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.GET("/folders", handler.ListFolders)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/folders", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
