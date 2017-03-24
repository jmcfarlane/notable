package main

import (
	"bytes"
	"encoding/gob"
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

// timeSorter sorts notes lines by last updated
type timeSorter Notes

func (a timeSorter) Len() int           { return len(a) }
func (a timeSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timeSorter) Less(i, j int) bool { return a[i].Updated < a[j].Updated }

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

// FromBytes converts encoding.Gob bytes into a Note
func (note *Note) FromBytes(b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&note)
}

// ToBytes converts a raw note into encoding.Gob bytes
func (note Note) ToBytes() ([]byte, error) {
	// Never include the password in a byte representation
	note.Password = ""

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(note)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
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
		encrypted, err := encrypt(note.Content, note.Password)
		if err != nil {
			return note, err
		}
		note.Content = encrypted
		note.Encrypted = true
	} else {
		note.Encrypted = false
	}
	return note, nil
}
