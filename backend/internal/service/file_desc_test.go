package service

import (
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestLanZouService_GetFileDescription_NotConnected(t *testing.T) {
	emptyRepo := &mockAccessLanzouTokenRepo{}
	svc := NewLanZouService(emptyRepo)

	_, err := svc.GetFileDescription(1, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestLanZouService_SetFileDescription_NotConnected(t *testing.T) {
	emptyRepo := &mockAccessLanzouTokenRepo{}
	svc := NewLanZouService(emptyRepo)

	_, err := svc.SetFileDescription(1, 123, "test desc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestLanZouService_GetFileDescription_ExpiredToken(t *testing.T) {
	expiredToken := &model.LanZouToken{
		UserID:    1,
		Cookie:    "test",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	repo := &mockAccessLanzouTokenRepo{token: expiredToken}
	svc := NewLanZouService(repo)

	_, err := svc.GetFileDescription(1, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestLanZouService_SetFileDescription_ExpiredToken(t *testing.T) {
	expiredToken := &model.LanZouToken{
		UserID:    1,
		Cookie:    "test",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	repo := &mockAccessLanzouTokenRepo{token: expiredToken}
	svc := NewLanZouService(repo)

	_, err := svc.SetFileDescription(1, 123, "test desc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}
