package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func setupTestRouterWithAuth(jwtManager *jwt.JWTManager, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthRequired(jwtManager))
	r.GET("/test", handler)
	return r
}

func TestAuthRequired_NoHeader(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440) // 24 hours in minutes

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidFormat(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthRequired_WrongScheme(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic sometoken")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken123")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	// Create a JWT manager with immediate expiry (0 duration means token expires right now)
	jwtManager := jwt.NewJWTManager("test-secret", 0)

	token, err := jwtManager.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Small delay to ensure token has expired
	time.Sleep(10 * time.Millisecond)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for expired token, got %d", w.Code)
	}
}

func TestAuthRequired_ValidToken(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	token, err := jwtManager.GenerateToken(42, "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var receivedUserID uint
	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		receivedUserID = userID.(uint)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if receivedUserID != 42 {
		t.Errorf("expected user_id 42, got %d", receivedUserID)
	}
}

func TestAuthRequired_DifferentSecret(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	token, err := jwtManager.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Create middleware with different secret
	wrongManager := jwt.NewJWTManager("wrong-secret", 1440)

	r := setupTestRouterWithAuth(wrongManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for wrong secret, got %d", w.Code)
	}
}

func TestAuthRequired_SetsUsername(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	token, err := jwtManager.GenerateToken(1, "alice")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var receivedUsername string
	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		username, _ := c.Get("username")
		receivedUsername = username.(string)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if receivedUsername != "alice" {
		t.Errorf("expected username 'alice', got '%s'", receivedUsername)
	}
}

func TestAuthRequired_EmptyBearer(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthRequired_BearerOnly(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 1440)

	r := setupTestRouterWithAuth(jwtManager, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
