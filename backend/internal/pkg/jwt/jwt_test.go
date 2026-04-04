package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestJWTManager() *JWTManager {
	return NewJWTManager("test-secret-key-1234567890", 60) // 60 minutes
}

func TestGenerateAndParseToken(t *testing.T) {
	mgr := newTestJWTManager()
	userID := uint(123)
	username := "testuser"

	token, err := mgr.GenerateToken(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}

func TestParseTokenInvalid(t *testing.T) {
	mgr := newTestJWTManager()
	_, err := mgr.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestParseTokenEmpty(t *testing.T) {
	mgr := newTestJWTManager()
	_, err := mgr.ValidateToken("")
	assert.Error(t, err)
}

func TestTokenExpiry(t *testing.T) {
	mgr := newTestJWTManager()
	userID := uint(1)
	username := "testuser"

	token, err := mgr.GenerateToken(userID, username)
	assert.NoError(t, err)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)

	// Check expiry is set correctly (within 24 hours)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.ExpiresAt.Before(time.Now().Add(61*time.Minute)))
}

func TestTokenContainsIssuedAt(t *testing.T) {
	mgr := newTestJWTManager()
	userID := uint(456)
	username := "issueduser"

	before := time.Now()
	token, _ := mgr.GenerateToken(userID, username)
	after := time.Now()

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)
	assert.True(t, claims.IssuedAt.After(before.Add(-time.Second)))
	assert.True(t, claims.IssuedAt.Before(after.Add(time.Second)))
}

func TestDifferentUsersGetDifferentTokens(t *testing.T) {
	mgr := newTestJWTManager()

	token1, _ := mgr.GenerateToken(uint(1), "user1")
	token2, _ := mgr.GenerateToken(uint(2), "user2")

	assert.NotEqual(t, token1, token2)
}

func TestSameUserGetsDifferentTokensDueToTime(t *testing.T) {
	mgr := newTestJWTManager()

	token1, _ := mgr.GenerateToken(uint(1), "user1")
	time.Sleep(time.Second)
	token2, _ := mgr.GenerateToken(uint(1), "user1")

	// Tokens should be different due to different IssuedAt times
	assert.NotEqual(t, token1, token2)

	// But both should be valid for same user
	claims1, _ := mgr.ValidateToken(token1)
	claims2, _ := mgr.ValidateToken(token2)
	assert.Equal(t, claims1.UserID, claims2.UserID)
	assert.Equal(t, claims1.Username, claims2.Username)
}
