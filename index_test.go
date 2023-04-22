package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnIndexError(t *testing.T) {
	assert.NotNil(t, unIndex(""))
}

func TestGetIndexError(t *testing.T) {
	idx, err := getIndex("")
	assert.Nil(t, idx)
	assert.NotNil(t, err)
}

func TestGetIndexNoteError(t *testing.T) {
	assert.NotNil(t, indexNote(Note{UID: ""}))
}

func TestReIndex(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	createTestNote(mock, "")
	count, err := reIndex(mock.db)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, 1, count, "Expected at least one note re-indexed")
}
