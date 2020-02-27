package models

import (
	"github.com/stretchr/testify/assert"
	. "no_homomorphism"
	"testing"
)

func TestKek(t *testing.T) {
	api := InitStorages()
	assert.NotNil(t, api)
}
