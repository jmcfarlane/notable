package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/twinj/uuid"

	// Imported only for it's side effect
	_ "github.com/mattn/go-sqlite3"
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

// connection to a sqlite database (currently hard coded for testing)
func connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/home/jmcfarlane/Desktop/notes.sqlite3")
	if err != nil {
		panic(err)
	}
	return db, err
}

// persistable prepares a note for being persisted to storage
func persistable(note Note) (Note, error) {
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

// DeleteByUID removes a note from storage
func DeleteByUID(uid string) error {
	if uid == "" {
		return errors.New("Deletion uid must not be an empty string")
	}
	db, _ := connection()
	defer db.Close()
	stmt, err := db.Prepare(`DELETE FROM notes WHERE uid=?`)
	if err != nil {
		log.Error(err)
		return err
	}
	res, err := stmt.Exec(uid)
	if err != nil {
		log.Error(err)
		return err
	}
	affected, _ := res.RowsAffected()
	log.Infof("Completed DB delete uid=%s, affected=%d", uid, affected)
	return nil
}

// Create a note
func Create(note Note) (Note, error) {
	note, err := persistable(note)
	if err != nil {
		panic(err)
	}
	db, _ := connection()
	defer db.Close()

	// No sql injection please
	stmt, err := db.Prepare(`
      INSERT INTO notes
        (content, created, encrypted, subject, tags, uid, updated)
	  VALUES
	    (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(
		note.Content,
		note.Encrypted,
		note.Subject,
		note.Tags,
		note.UID,
		note.UID,
		note.Updated)
	if err != nil {
		panic(err)
	}
	log.Infof("Completed DB insert uid=%s", note.UID)

	return note, nil
}

// GetContentByUID fetches content (which might be encrypted) by uid.
func GetContentByUID(uid string, password string) (string, error) {
	db, _ := connection()
	defer db.Close()
	rows, _ := db.Query("SELECT content FROM notes WHERE uid=?", uid)
	var notes Notes
	for rows.Next() {
		var note Note
		rows.Scan(&note.Content)
		notes = append(notes, note)
	}
	if len(notes) == 1 {
		content := notes[0].Content
		if password == "" {
			return content, nil
		}
		return Decrypt(content, password)
	}

	return "", fmt.Errorf("No note found")
}

// Now is the current timestamp in time.RFC3339 format
func Now() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

// Search all notes as filtered by the provided query
func Search(query string) Notes {
	db, _ := connection()
	defer db.Close()
	rows, _ := db.Query(`
		SELECT
			created, encrypted, subject, tags, uid, updated
		FROM
			notes
		ORDER BY updated DESC`)
	defer rows.Close()

	var notes Notes
	for rows.Next() {
		var note Note
		rows.Scan(
			&note.Created,
			&note.Encrypted,
			&note.Subject,
			&note.Tags,
			&note.UID,
			&note.Updated)
		notes = append(notes, note)
	}

	return notes
}

// Update a note
func Update(note Note) (Note, error) {
	note, err := persistable(note)
	if err != nil {
		panic(err)
	}
	db, _ := connection()
	defer db.Close()

	// No sql injection please
	stmt, err := db.Prepare(`
      UPDATE notes SET
        content = ?,
        encrypted = ?,
        subject = ?,
        tags = ?,
        updated = ?
      WHERE uid = ?
	`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(
		note.Content,
		note.Encrypted,
		note.Subject,
		note.Tags,
		note.Updated,
		note.UID)
	if err != nil {
		panic(err)
	}
	affected, _ := res.RowsAffected()
	log.Infof("Completed DB update uid=%s, affected=%d", note.UID, affected)

	return note, nil
}
