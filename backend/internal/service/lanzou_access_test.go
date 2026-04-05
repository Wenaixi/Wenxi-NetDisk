package service

import (
	"testing"
	"time"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockAccessLanzouTokenRepo struct {
	token *model.LanZouToken
}

func (m *mockAccessLanzouTokenRepo) Upsert(token *model.LanZouToken) error {
	m.token = token
	return nil
}
func (m *mockAccessLanzouTokenRepo) FindByUserID(userID uint) (*model.LanZouToken, error) {
	return m.token, nil
}
func (m *mockAccessLanzouTokenRepo) DeleteByUserID(userID uint) error {
	return nil
}

func TestLanZouService_SetFileAccess_NotConnected(t *testing.T) {
	emptyRepo := &mockAccessLanzouTokenRepo{}
	svc := NewLanZouService(emptyRepo)

	_, err := svc.SetFileAccess(1, 123, 2, "pwd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestLanZouService_SetFolderAccess_NotConnected(t *testing.T) {
	emptyRepo := &mockAccessLanzouTokenRepo{}
	svc := NewLanZouService(emptyRepo)

	_, err := svc.SetFolderAccess(1, 456, 2, "pwd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestLanZouService_SetFileAccess_ExpiredToken(t *testing.T) {
	expiredToken := &model.LanZouToken{
		UserID:    1,
		Cookie:    "test",
		ExpiresAt: time.Now().Add(-24 * time.Hour), // expired
	}
	repo := &mockAccessLanzouTokenRepo{token: expiredToken}
	svc := NewLanZouService(repo)

	_, err := svc.SetFileAccess(1, 123, 1, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}
