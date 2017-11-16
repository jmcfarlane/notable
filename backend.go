package main

import (
	"os"
	"path/filepath"
)

const connectionErrFmt = "Unable to connect driver=%s path=%s err=%v"

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

// Backend system
type Backend interface {
	close()
	create(Note) (Note, error)
	createSchema()
	dbFilePath() string
	deleteByUID(string) error
	getNoteByUID(string, password string) (Note, error)
	list() Notes
	String() string
	update(Note) (Note, error)
}

func decryptNote(note Note, password string) (Note, error) {
	if password != "" {
		clearText, err := Decrypt(note, password)
		if err != nil {
			return note, err
		}
		note.Content = clearText
	}
	return note, nil
}

func getContentByUID(b Backend, uid string, password string) (string, error) {
	note, err := b.getNoteByUID(uid, password)
	if err != nil {
		return "", err
	}
	return note.Content, nil
}
