package database

import (
	"database/sql"
	"fmt"

	// Imported only for it's side effect
	_ "github.com/mattn/go-sqlite3"
)

// Note represents a single note stored by the application
type Note struct {
	Content   string `json:"content"`
	Created   string `json:"created"`
	Encrypted bool   `json:"encrypted"`
	Subject   string `json:"subject"`
	Tags      string `json:"tags"`
	UID       string `json:"uid"`
	Updated   string `json:"updated"`
}

// Notes is a collection of Note objects
type Notes []Note

func connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/home/jmcfarlane/Sync/Notes/notes.sqlite3")
	if err != nil {
		panic(err)
	}
	return db, err
}

// DecryptContent returns the passed content decrypted by the password
func DecryptContent(content string, password string) (string, error) {
	return content, nil
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
		return Decrypt(notes[0].Content, password)
	}

	return "", fmt.Errorf("No note found")
}

// Search fetches all notes as filtered by the provided query
func Search(query string) Notes {
	db, _ := connection()
	defer db.Close()
	rows, _ := db.Query("SELECT created, encrypted, subject, tags, uid, updated FROM notes ORDER BY updated DESC")
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
