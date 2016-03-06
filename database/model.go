package database

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/twinj/uuid"
)

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

// Now is the current timestamp in time.RFC3339 format
func Now() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

// Persistable prepares a note for being persisted to storage
func Persistable(note Note) (Note, error) {
	note.Updated = Now()
	// Generate a uuid if necessary
	if note.UID == "" {
		note.UID = uuid.NewV4().String()
	}
	// Make sure the contents are encrypted if a password is set
	if note.Password != "" {
		note.Content = Encrypt(note.Content, note.Password)
		note.Encrypted = true
	} else {
		note.Encrypted = false
	}
	return note, nil
}
