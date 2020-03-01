package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUsersStorage(t *testing.T) {
	userModel := NewUsersStorage()
	assert.NotNil(t, userModel)
	trackModel := NewTrackStorage()
	assert.NotNil(t, trackModel)
}
