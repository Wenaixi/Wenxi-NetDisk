package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupResponseTestRouter(handler gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", handler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	return w
}

func TestSuccess(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Success(c, map[string]string{"key": "value"})
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSuccess_WithNilData(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Success(c, nil)
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSuccess_WithStringData(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Success(c, "ok")
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSuccess_WithSliceData(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Success(c, []string{"a", "b"})
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestBadRequest(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		BadRequest(c, "invalid input")
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Unauthorized(c, "token expired")
	})

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		NotFound(c, "resource not found")
	})

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestInternalError(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		InternalError(c, "database error")
	})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestError_CustomHTTPStatus(t *testing.T) {
	w := setupResponseTestRouter(func(c *gin.Context) {
		Error(c, http.StatusTeapot, "I'm a teapot")
	})

	if w.Code != http.StatusTeapot {
		t.Errorf("expected status 418, got %d", w.Code)
	}
}
