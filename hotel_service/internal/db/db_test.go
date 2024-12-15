package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase()

	assert.NoError(t, err)
	assert.NotNil(t, db)
}
