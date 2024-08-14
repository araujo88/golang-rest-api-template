package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "1234"
	hashedPassword, err := HashPassword(password)
	assert.Nil(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestGenerateToken(t *testing.T) {
	user := "chud"
	token, err := GenerateToken(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRandomKey(t *testing.T) {
	randomKey := GenerateRandomKey()
	assert.NotEmpty(t, randomKey)
	assert.Len(t, randomKey, 44)
}
