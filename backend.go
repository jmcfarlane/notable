package main

import (
	"os"
	"path/filepath"

	"github.com/jmcfarlane/notable/app"
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
	create(app.Note) (app.Note, error)
	createSchema()
	dbFilePath() string
	deleteByUID(string) error
	getNoteByUID(string, password string) (app.Note, error)
	list() app.Notes
	String() string
	update(app.Note) (app.Note, error)
}

func decryptNote(note app.Note, password string) (app.Note, error) {
	if password != "" {
		decrypted, err := app.Decrypt(note.Content, password)
		if err != nil {
			return note, err
		}
		note.Content = decrypted
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
