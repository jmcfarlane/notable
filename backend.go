package main

import (
	"os"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
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
