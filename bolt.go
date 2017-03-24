package main

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/boltdb/bolt"

	log "github.com/Sirupsen/logrus"
)

// BoltDB backend
type BoltDB struct {
	Engine      *bolt.DB
	Path        string
	Type        string
	NotesBucket []byte
}

// NewBoltDB engine instance
func NewBoltDB(path string) (*BoltDB, error) {
	db := &BoltDB{
		NotesBucket: []byte("notes"),
		Path:        path,
		Type:        "BoltDB",
	}
	_, fileExisted := createParentDirs(path)
	engine, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return db, err
	}
	db.Engine = engine
	if !fileExisted {
		db.createSchema()
		db.migrate()
	}
	return db, nil
}

func (db *BoltDB) String() string {
	return fmt.Sprintf("type=%s path=%s", db.Type, db.Path)
}

func (db *BoltDB) close() {
	db.Engine.Close()
}

func (db *BoltDB) dbFilePath() string {
	return db.Path
}

func (db *BoltDB) create(note Note) (Note, error) {
	return db.update(note)
}

func (db *BoltDB) createSchema() {
	err = db.Engine.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket(db.NotesBucket)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (db *BoltDB) deleteByUID(uid string) error {
	err = db.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		return bucket.Delete([]byte(uid))
	})
	return err
}

func (db *BoltDB) getNoteByUID(uid string, password string) (Note, error) {
	var note Note
	err := db.Engine.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.NotesBucket)
		v := b.Get([]byte(uid))
		return note.FromBytes(v)
	})
	if err != nil {
		return note, err
	}
	return decryptNote(note, password)
}

func (db *BoltDB) migrate() {
	oldDBPath := filepath.Join(filepath.Dir(*dbPath), "notes.sqlite3")
	if !fileExists(oldDBPath) {
		return
	}
	oldDB, err := NewSqlite3(oldDBPath)
	if err != nil {
		log.Panic(err)
	}
	notes := oldDB.fetchAll()
	err = db.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		for _, note := range notes {
			b, err := note.ToBytes()
			if err != nil {
				log.Fatal(err)
			}
			bucket.Put([]byte(note.UID), b)
			fmt.Println("Migrated to BoltDB:", note.Subject, bucket)
		}
		return nil
	})
	log.Infof("Migration complete err=%v", err)
}

func (db *BoltDB) search(query string) Notes {
	var notes Notes
	db.Engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var note Note
			err := note.FromBytes(v)
			if err != nil {
				log.Fatal(err)
			}
			note.Content = ""
			notes = append(notes, note)
		}
		return nil
	})
	sort.Sort(sort.Reverse(timeSorter(notes)))
	return notes
}

func (db *BoltDB) update(note Note) (Note, error) {
	note, err := persistable(note)
	if err != nil {
		return note, err
	}
	err = db.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		b, err := note.ToBytes()
		if err != nil {
			return err
		}
		bucket.Put([]byte(note.UID), b)
		return nil
	})
	return note, err
}
