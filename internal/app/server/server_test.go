package server

import (
	"github.com/stretchr/testify/assert"
	"no_homomorphism"

	"testing"
)

func TestKek(t *testing.T) {
	api := main.InitStorages()
	assert.NotNil(t, api)
}
