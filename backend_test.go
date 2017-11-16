package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	assert.True(t, pathExists("/bin"))
	assert.False(t, pathExists("/bin-never-named-this"))
}
