package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenBoltDBWithError(t *testing.T) {
	db, err := openBoltDB("/proc/test/notes.db", false)
	assert.Nil(t, db)
	assert.NotNil(t, err)
}

func TestSecondaryIsNil(t *testing.T) {
	assert.True(t, isNil(nil))
	assert.False(t, isNil(&Secondary{Path: "bla"}))
}

func TestBoltDbStringer(t *testing.T) {
	db := BoltDB{Path: "primary"}
	dbString := db.String()
	assert.Contains(t, dbString, "path=primary")
	assert.Contains(t, dbString, "secondary=false")
}

func TestBoltDbStringerWithSecondary(t *testing.T) {
	db := BoltDB{
		Path:      "secondary",
		Secondary: &Secondary{Path: "bla"},
	}
	dbString := db.String()
	assert.Contains(t, dbString, "secondary=true")
}
