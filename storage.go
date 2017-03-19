package main

import (
	"os"
	"path/filepath"
)

const connectionErrFmt = "Unable to connect driver=%s path=%s err=%v"

// Backend system
type Backend interface {
	close()
	create(Note) (Note, error)
	createSchema()
	dbFilePath() string
	deleteByUID(string) error
	getContentByUID(string, password string) (string, error)
	search(string) Notes
	String() string
	update(Note) (Note, error)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createParentDirs(path string) (bool, bool) {
	d := filepath.Dir(path)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		os.MkdirAll(d, 0777)
		return true, false
	}
	return false, fileExists(path)
}
