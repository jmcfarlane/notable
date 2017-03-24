package main

import (
	"database/sql"
	"errors"
	"fmt"

	// Imported only for it's side effect

	_ "github.com/mattn/go-sqlite3"

	log "github.com/Sirupsen/logrus"
)

// Sqlite3 backend
type Sqlite3 struct {
	Engine *sql.DB
	Path   string
	Type   string
}

// NewSqlite3 engine instance
func NewSqlite3(path string) (*Sqlite3, error) {
	db := &Sqlite3{Path: path, Type: "Sqlite3"}
	_, fileExisted := createParentDirs(path)
	engine, err := sql.Open("sqlite3", path)
	if err != nil {
		return db, fmt.Errorf(connectionErrFmt, "sqlite3", path, err)
	}
	db.Engine = engine
	if !fileExisted {
		db.createSchema()
	}
	return db, nil
}

func (db *Sqlite3) String() string {
	return fmt.Sprintf("type=%s path=%s", db.Type, db.Path)
}

func (db *Sqlite3) close() {
	db.Engine.Close()
}

func (db *Sqlite3) dbFilePath() string {
	return db.Path
}

func (db *Sqlite3) createSchema() {
	stmt, err := db.Engine.Prepare(`
		CREATE TABLE notes (
			uid string,
			created string,
			updated string,
			tags string,
			content string,
			encrypted INTEGER DEFAULT 0,
			subject TEXT
		);`)
	if err != nil {
		if err.Error() != "table notes already exists" {
			log.Panicf("Unable to prepare schema path=%s, err=%v", db.Path, err)
		}
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Errorf("Unable to create schema err=%v", err)
	}
	log.Infof("Schema created path=%s", db.Path)
}

// Removes a note from storage
func (db *Sqlite3) deleteByUID(uid string) error {
	if uid == "" {
		return errors.New("Deletion uid must not be an empty string")
	}
	stmt, err := db.Engine.Prepare(`DELETE FROM notes WHERE uid=?`)
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
func (db *Sqlite3) create(note Note) (Note, error) {
	note, err := persistable(note)
	if err != nil {
		panic(err)
	}
	// No sql injection please
	stmt, err := db.Engine.Prepare(`
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
		note.Created,
		note.Encrypted,
		note.Subject,
		note.Tags,
		note.UID,
		note.Updated)
	if err != nil {
		panic(err)
	}
	log.Infof("Completed DB insert uid=%s", note.UID)

	return note, nil
}

func (db *Sqlite3) getNoteByUID(uid string, password string) (Note, error) {
	rows, _ := db.Engine.Query("SELECT content FROM notes WHERE uid=?", uid)
	var notes Notes
	for rows.Next() {
		var note Note
		rows.Scan(&note.Content)
		notes = append(notes, note)
	}
	if len(notes) != 1 {
		return Note{}, fmt.Errorf("No note found")
	}
	return decryptNote(notes[0], password)
}

func (db *Sqlite3) fetchAll() Notes {
	rows, _ := db.Engine.Query(`
		SELECT
			content, created, encrypted, subject, tags, uid, updated
		FROM
			notes
		`)
	defer rows.Close()

	var notes Notes
	for rows.Next() {
		var note Note
		rows.Scan(
			&note.Content,
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

func (db *Sqlite3) search(query string) Notes {
	rows, _ := db.Engine.Query(`
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

func (db *Sqlite3) update(note Note) (Note, error) {
	note, err := persistable(note)
	if err != nil {
		panic(err)
	}
	// No sql injection please
	stmt, err := db.Engine.Prepare(`
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
