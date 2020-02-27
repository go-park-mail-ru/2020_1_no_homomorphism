package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUsersStorage(t *testing.T) {
	userModel := NewUsersStorage()
	assert.NotNil(t, userModel)
	trackModel := NewTrackStorage()
	assert.NotNil(t, trackModel)
}
