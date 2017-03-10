package main

import (
	"encoding/json"
	"time"

	"github.com/twinj/uuid"

	log "github.com/Sirupsen/logrus"
)

// Envelope to communicate details to the frontent
type apiResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Note represents a single note stored by the application
type Note struct {
	Content   string `json:"content"`
	Created   string `json:"created"`
	Encrypted bool   `json:"encrypted"`
	Password  string `json:"password"`
	Subject   string `json:"subject"`
	Tags      string `json:"tags"`
	UID       string `json:"uid"`
	Updated   string `json:"updated"`
}

// Notes is a collection of Note objects
type Notes []Note

// ToJSON converts a (filtered) note into json fields filtered:
// - Content: For performance reasons
// - Password: For security reasons
func (note Note) ToJSON() (string, error) {
	note.Content = ""
	note.Password = ""
	noteJSON, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		log.Error("Failed to parse note into json", err)
		return "", err
	}
	return string(noteJSON), err
}

// The current timestamp in time.RFC3339 format
func now() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

// Prepare a note for being persisted to storage
func persistable(note Note) (Note, error) {
	note.Updated = now()
	// Generate a uuid if necessary
	if note.UID == "" {
		note.UID = uuid.NewV4().String()
	}
	// Make sure the contents are encrypted if a password is set
	if note.Password != "" {
		note.Content = encrypt(note.Content, note.Password)
		note.Encrypted = true
	} else {
		note.Encrypted = false
	}
	return note, nil
}
