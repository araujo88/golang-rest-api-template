package cache

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisClient := NewMockCache(ctrl)
	assert.NotNil(t, redisClient)
}
