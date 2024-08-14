package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase()
	assert.NotNil(t, db)
}
