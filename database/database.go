package database

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

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

	// Don't pass this one around unless really necessary
	Password string `json:"-"`
}

// Notes is a collection of Note objects
type Notes []Note

func connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/home/jmcfarlane/Desktop/notes.sqlite3")
	if err != nil {
		panic(err)
	}
	return db, err
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

// Now is a current timestamp in string format
func Now() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

// Search fetches all notes as filtered by the provided query
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
	// Make sure the timestamp reflects the update
	note.Updated = Now()
	// Make sure the contents are encrypted if a password is set
	if note.Password != "" {
		note.Content = Encrypt(note.Content, note.Password)
	}
	db, _ := connection()
	defer db.Close()

	// No sql injection please
	stmt, err := db.Prepare(`
      UPDATE notes SET
        tags = ?,
        subject = ?,
        content = ?,
        encrypted = ?,
        updated = ?
      WHERE uid = ?
	`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(
		note.Tags,
		note.Subject,
		note.Content,
		note.Encrypted,
		note.Updated,
		note.UID)
	if err != nil {
		panic(err)
	}
	affect, _ := res.RowsAffected()
	log.Infof("Completed DB update uid=%s, affected=%d", note.UID, affect)

	return note, nil
}
