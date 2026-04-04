package service

import (
	"errors"
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

// mockLanzouTokenRepo 模拟蓝奏云token仓库
type mockLanzouTokenRepo struct {
	tokens map[uint]*model.LanZouToken
}

func newMockLanzouTokenRepo() *mockLanzouTokenRepo {
	return &mockLanzouTokenRepo{
		tokens: make(map[uint]*model.LanZouToken),
	}
}

func (m *mockLanzouTokenRepo) Upsert(token *model.LanZouToken) error {
	m.tokens[token.UserID] = token
	return nil
}

func (m *mockLanzouTokenRepo) FindByUserID(userID uint) (*model.LanZouToken, error) {
	token, ok := m.tokens[userID]
	if !ok {
		return nil, errors.New("token not found")
	}
	return token, nil
}

func (m *mockLanzouTokenRepo) DeleteByUserID(userID uint) error {
	delete(m.tokens, userID)
	return nil
}

// TestLanZouService_SaveToken tests token saving
func TestLanZouService_SaveToken(t *testing.T) {
	t.Run("should save token with expiration", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		err := svc.SaveToken(1, "test-cookie", "test-token")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		token, err := repo.FindByUserID(1)
		if err != nil {
			t.Errorf("expected token, got error: %v", err)
			return
		}
		if token.Cookie != "test-cookie" {
			t.Errorf("expected cookie 'test-cookie', got '%s'", token.Cookie)
		}
		if token.TokenValue != "test-token" {
			t.Errorf("expected token value 'test-token', got '%s'", token.TokenValue)
		}
		if !token.ExpiresAt.After(time.Now().Add(29*24*time.Hour)) {
			t.Error("expected token to expire in ~30 days")
		}
	})

	t.Run("should update existing token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		// Save first token
		_ = svc.SaveToken(1, "old-cookie", "old-token")

		// Update with new token
		err := svc.SaveToken(1, "new-cookie", "new-token")
		if err != nil {
			t.Errorf("expected no error on update, got %v", err)
		}

		token, _ := repo.FindByUserID(1)
		if token.Cookie != "new-cookie" {
			t.Errorf("expected updated cookie 'new-cookie', got '%s'", token.Cookie)
		}
		if token.TokenValue != "new-token" {
			t.Errorf("expected updated token value 'new-token', got '%s'", token.TokenValue)
		}
	})
}

// TestLanZouService_GetToken tests token retrieval
func TestLanZouService_GetToken(t *testing.T) {
	t.Run("should return existing token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		_ = svc.SaveToken(1, "test-cookie", "test-token")

		token, err := svc.GetToken(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == nil {
			t.Error("expected token, got nil")
			return
		}
		if token.Cookie != "test-cookie" {
			t.Errorf("expected cookie 'test-cookie', got '%s'", token.Cookie)
		}
	})

	t.Run("should return error for non-existent token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		token, err := svc.GetToken(999)

		if err == nil {
			t.Error("expected error for non-existent token")
		}
		if token != nil {
			t.Error("expected nil token for non-existent user")
		}
	})
}

// TestLanZouService_DeleteToken tests token deletion
func TestLanZouService_DeleteToken(t *testing.T) {
	t.Run("should delete existing token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		_ = svc.SaveToken(1, "test-cookie", "test-token")

		err := svc.DeleteToken(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		_, err = repo.FindByUserID(1)
		if err == nil {
			t.Error("expected error after deletion")
		}
	})

	t.Run("should succeed for non-existent token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		err := svc.DeleteToken(999)

		// DeleteByUserID on empty map should not error
		if err != nil {
			t.Errorf("expected no error for non-existent token deletion, got %v", err)
		}
	})
}

// TestLanZouService_IsConnected tests connection status
func TestLanZouService_IsConnected(t *testing.T) {
	t.Run("should return true for valid non-expired token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		_ = svc.SaveToken(1, "test-cookie", "test-token")

		if !svc.IsConnected(1) {
			t.Error("expected user to be connected")
		}
	})

	t.Run("should return false for non-existent token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		if svc.IsConnected(999) {
			t.Error("expected user to not be connected")
		}
	})

	t.Run("should return false for expired token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		// Manually insert an expired token
		repo.tokens[1] = &model.LanZouToken{
			UserID:     1,
			Cookie:     "expired-cookie",
			TokenValue: "expired-token",
			ExpiresAt:  time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		}

		if svc.IsConnected(1) {
			t.Error("expected user to not be connected with expired token")
		}
	})
}

// TestLanZouService_GetClient tests client creation
func TestLanZouService_GetClient(t *testing.T) {
	t.Run("should create client for valid token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		_ = svc.SaveToken(1, "test-cookie", "test-token")

		client, err := svc.GetClient(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if client == nil {
			t.Error("expected client, got nil")
		}
	})

	t.Run("should return error for non-existent token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		client, err := svc.GetClient(999)

		if err == nil {
			t.Error("expected error for non-existent token")
		}
		if client != nil {
			t.Error("expected nil client for non-existent user")
		}
	})

	t.Run("should return error for expired token", func(t *testing.T) {
		repo := newMockLanzouTokenRepo()
		svc := NewLanZouService(repo)

		// Manually insert an expired token
		repo.tokens[1] = &model.LanZouToken{
			UserID:     1,
			Cookie:     "expired-cookie",
			TokenValue: "expired-token",
			ExpiresAt:  time.Now().Add(-1 * time.Hour),
		}

		client, err := svc.GetClient(1)

		if err == nil {
			t.Error("expected error for expired token")
		}
		if client != nil {
			t.Error("expected nil client for expired token")
		}
	})
}

// mockLanzouClientProvider 模拟蓝奏云客户端提供者
type mockLanzouClientProvider struct {
	connected bool
}

func (m *mockLanzouClientProvider) IsConnected(userID uint) bool {
	return m.connected
}

func (m *mockLanzouClientProvider) GetClient(userID uint) (*lanzou.Client, error) {
	if !m.connected {
		return nil, errors.New("lanzou not connected")
	}
	return lanzou.NewClient("test-cookie"), nil
}
