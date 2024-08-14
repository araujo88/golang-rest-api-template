package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	redisClient := NewRedisClient()
	assert.NotNil(t, redisClient)
}
