package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKek(t *testing.T) {
	api := InitStorages()
	assert.NotNil(t, api)
}
