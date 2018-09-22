package main

import (
	"errors"
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
	// Encrypted notes saved prior to 0.1.2 (2017-12-10) don't have a
	// good way to determine if decryption was successful. This crappy
	// mechanism only allowed encrypted notes to use latin1 characters
	// TODO: Remove this check some time past 2019 :)
	if password != "" && note.CipherType == "" && SmellsEncrypted(note.Content) {
		msg := "Decryption of old note encrypted with AES-CBC failed"
		return "", errors.New(msg)
	}
	return note.Content, nil
}
