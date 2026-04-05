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

// TestValidateToken_Expired 测试过期 token 被拒绝
func TestValidateToken_Expired(t *testing.T) {
	// Create a manager with 1ms expiry
	mgr := NewJWTManager("test-secret", 0)

	token, err := mgr.GenerateToken(uint(1), "expireuser")
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(100 * time.Millisecond)

	_, err = mgr.ValidateToken(token)
	assert.Error(t, err)
}

// TestValidateToken_WrongSecret 测试错误密钥
func TestValidateToken_WrongSecret(t *testing.T) {
	mgr1 := NewJWTManager("secret-one", 60)
	mgr2 := NewJWTManager("secret-two", 60)

	token, _ := mgr1.GenerateToken(uint(1), "user1")

	_, err := mgr2.ValidateToken(token)
	assert.Error(t, err)
}

// TestValidateToken_Malformed 测试畸形 token
func TestValidateToken_Malformed(t *testing.T) {
	mgr := newTestJWTManager()

	// Various malformed inputs
	malformed := []string{
		"not.a.token",
		"eyJhbGciOiJIUzI1NiJ9",        // Only header
		"eyJhbGciOiJIUzI1NiJ9.",       // Header + dot
		"eyJhbGciOiJIUzI1NiJ9.e30.",   // Header + payload + dot (no signature)
		"abc.def.ghi",
		"\x00\x01\x02",
	}

	for _, m := range malformed {
		_, err := mgr.ValidateToken(m)
		assert.Error(t, err, "expected error for malformed token: %s", m)
	}
}

// TestJWTManager_EmptySecret 测试空密钥
func TestJWTManager_EmptySecret(t *testing.T) {
	mgr := NewJWTManager("", 60)

	token, err := mgr.GenerateToken(uint(1), "user1")
	assert.NoError(t, err)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
}

// TestJWTManager_ShortExpire 测试极短过期时间
func TestJWTManager_ShortExpire(t *testing.T) {
	mgr := NewJWTManager("test-secret", 1) // 1 minute

	token, err := mgr.GenerateToken(uint(1), "user1")
	assert.NoError(t, err)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)

	// Check expiry is ~1 minute from now
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.ExpiresAt.Before(time.Now().Add(2*time.Minute)))
}

// TestJWTManager_LargeUserID 测试大用户ID
func TestJWTManager_LargeUserID(t *testing.T) {
	mgr := newTestJWTManager()

	token, err := mgr.GenerateToken(uint(1<<32-1), "user")
	assert.NoError(t, err)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(1<<32-1), claims.UserID)
}

// TestJWTManager_SpecialCharsInUsername 测试特殊字符用户名
func TestJWTManager_SpecialCharsInUsername(t *testing.T) {
	mgr := newTestJWTManager()

	token, err := mgr.GenerateToken(uint(1), "user@domain.com/<script>")
	assert.NoError(t, err)

	claims, err := mgr.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user@domain.com/<script>", claims.Username)
}
