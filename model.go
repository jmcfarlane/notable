package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/twinj/uuid"

	log "github.com/sirupsen/logrus"
)

// APIResponse - Envelope to communicate details to the frontent
type APIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Note represents a single note stored by the application
type Note struct {
	Autosave  bool   `json:"autosave"`
	Content   string `json:"content"`
	Created   string `json:"created"`
	Encrypted bool   `json:"encrypted"`
	Password  string `json:"password"`
	Subject   string `json:"subject"`
	Tags      string `json:"tags"`
	UID       string `json:"uid"`
	Updated   string `json:"updated"`

	// Fields largely to support the `-secondary` feature
	AheadOfPrimary bool      `json:"ahead_of_primary"`
	Deleted        bool      `json:"deleted"`
	Time           time.Time `json:"time"`

	// Private fields
	CipherType    string `json:"-"`
	SecondaryPath string `json:"-"`
}

// Notes is a collection of Note objects
type Notes []Note

// Map of notes key'd by uid
func Map(notes Notes) map[string]Note {
	var m = make(map[string]Note)
	for _, note := range notes {
		m[note.UID] = note
	}
	return m
}

// TimeSorter sorts notes lines by last updated
type TimeSorter Notes

func (a TimeSorter) Len() int           { return len(a) }
func (a TimeSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeSorter) Less(i, j int) bool { return a[i].Updated < a[j].Updated }

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
	// Zero byte file represents a logical delete
	if len(b) == 0 {
		note.Deleted = true
		return nil
	}
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

// Persistable prepares a note for being persisted to storage
func Persistable(note Note) (Note, error) {
	note.Time = time.Now()
	note.Updated = note.Time.Format(time.RFC3339)
	// Generate a uuid if necessary
	if note.UID == "" {
		note.UID = uuid.NewV4().String()
	}
	// Make sure the contents are encrypted if a password is set
	if note.Password != "" {
		ciperText, cipherType, err := Encrypt(note)
		if err != nil {
			return note, err
		}
		note.CipherType = cipherType
		note.Content = ciperText
		note.Encrypted = true
	} else {
		note.Encrypted = false
	}
	return note, nil
}
