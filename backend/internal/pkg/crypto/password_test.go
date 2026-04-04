package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHash(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, len(hash) > 50) // bcrypt hashes are typically 60 chars

	// Test correct password
	result := CheckPassword(password, hash)
	assert.True(t, result)

	// Test wrong password
	result = CheckPassword("wrongpassword", hash)
	assert.False(t, result)
}

func TestPasswordHashEmpty(t *testing.T) {
	hash, err := HashPassword("")
	assert.NoError(t, err) // bcrypt accepts empty passwords
	assert.NotEmpty(t, hash)
	assert.True(t, CheckPassword("", hash))
}

func TestPasswordHashShort(t *testing.T) {
	hash, err := HashPassword("short")
	assert.NoError(t, err) // bcrypt accepts short passwords
	assert.NotEmpty(t, hash)
	assert.True(t, CheckPassword("short", hash))
}

func TestCheckPasswordInvalidHash(t *testing.T) {
	result := CheckPassword("password", "invalid-hash")
	assert.False(t, result)
}

func TestCheckPasswordEmptyVsNonEmpty(t *testing.T) {
	hash1, _ := HashPassword("password1")
	hash2, _ := HashPassword("password2")
	assert.NotEqual(t, hash1, hash2, "Same password should produce different hashes (salted)")

	// Verify original passwords work
	assert.True(t, CheckPassword("password1", hash1))
	assert.True(t, CheckPassword("password2", hash2))
	assert.False(t, CheckPassword("password1", hash2))
}
