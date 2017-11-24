package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeSorter(t *testing.T) {
	notes := []Note{
		Note{Updated: "2006-01-02T15:04:05Z07:00"},
		Note{Updated: "2007-01-02T15:04:05Z07:00"},
		Note{Updated: "2005-01-02T15:04:05Z07:00"},
	}
	sort.Sort(TimeSorter(notes))
	assert.True(t, strings.HasPrefix(notes[0].Updated, "2005"))
	assert.True(t, strings.HasPrefix(notes[1].Updated, "2006"))
	assert.True(t, strings.HasPrefix(notes[2].Updated, "2007"))
}

func TestTimeSorterReverse(t *testing.T) {
	notes := []Note{
		Note{Updated: "2006-01-02T15:04:05Z07:00"},
		Note{Updated: "2007-01-02T15:04:05Z07:00"},
		Note{Updated: "2005-01-02T15:04:05Z07:00"},
	}
	sort.Sort(sort.Reverse(TimeSorter(notes)))
	assert.True(t, strings.HasPrefix(notes[0].Updated, "2007"))
	assert.True(t, strings.HasPrefix(notes[1].Updated, "2006"))
	assert.True(t, strings.HasPrefix(notes[2].Updated, "2005"))
}
