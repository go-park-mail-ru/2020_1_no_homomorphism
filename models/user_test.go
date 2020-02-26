package models

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)



func Test_NewUsersStorage(t *testing.T) {
	mu := &sync.Mutex{}
	_, err := NewUsersStorage(nil)
	assert.NotNil(t, err)
	_, err = NewUsersStorage(mu)
	assert.Nil(t, err)
}
