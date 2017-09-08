package main

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/jmcfarlane/notable/app"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// BoltDB backend
type BoltDB struct {
	Engine      *bolt.DB
	Path        string
	Type        string
	NotesBucket []byte

	// Secondary nodes do not have direct write access to the
	// database. They write change files which are consumed by the
	// primary process.
	Secondary *Secondary
}

func openBoltDB(path string, secondary bool) (*BoltDB, error) {
	db := &BoltDB{
		NotesBucket: []byte("notes"),
		Path:        path,
		Type:        "BoltDB",
	}
	if secondary {
		db.Secondary = &Secondary{
			Path: db.Path,
		}
	}
	_, fileExisted := createParentDirs(path)
	engine, err := bolt.Open(path, 0600, &bolt.Options{
		ReadOnly: secondary,
		Timeout:  *boltTimeout,
	})
	if err != nil {
		return db, err
	}
	db.Engine = engine
	if !secondary && !fileExisted {
		db.createSchema()
		db.migrate()
	}
	return db, nil
}

func isNil(s *Secondary) bool {
	if s == nil {
		return true
	}
	return false
}

func (db *BoltDB) String() string {
	return fmt.Sprintf("type=%s path=%s secondary=%v", db.Type, db.Path, !isNil(db.Secondary))
}

func (db *BoltDB) close() {
	db.Engine.Close()
}

func (db *BoltDB) dbFilePath() string {
	return db.Path
}

func (db *BoltDB) create(note app.Note) (app.Note, error) {
	return db.update(note)
}

func (db *BoltDB) createSchema() {
	err := db.Engine.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket(db.NotesBucket)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (db *BoltDB) deleteByUID(uid string) error {
	if db.Secondary != nil {
		return db.Secondary.deleteByUID(uid)
	}
	if err := unIndex(uid); err != nil {
		return err
	}
	err := db.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		return bucket.Delete([]byte(uid))
	})
	return err
}

func (db *BoltDB) getNoteByUID(uid string, password string) (app.Note, error) {
	if db.Secondary != nil {
		notes := db.Secondary.list()
		// Sort in reverse order, so the FIRST note wins.
		sort.Sort(sort.Reverse(app.TimeSorter(notes)))
		for _, note := range notes {
			if note.UID == uid {
				return decryptNote(note, password)
			}
		}
	}
	var note app.Note
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
	oldDB, err := openSqlite3(oldDBPath)
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

func (db *BoltDB) list() app.Notes {
	var notes app.Notes
	var updates map[string]app.Note
	if db.Secondary != nil {
		updatedNotes := db.Secondary.list()
		// Sort in ascending order, so the LAST note wins.
		sort.Sort(app.TimeSorter(updatedNotes))
		updates = app.Map(updatedNotes)
	}
	db.Engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// Take any local updates over what the db has (ignoring deletes)
			if note, ok := updates[string(k)]; ok {
				if !note.Deleted {
					note.Content = ""
					notes = append(notes, note)
					continue
				}
			}
			// From the primary database
			var note app.Note
			err := note.FromBytes(v)
			if err != nil {
				log.Fatal(err)
			}
			note.Content = ""
			notes = append(notes, note)
		}
		return nil
	})

	// Add any secondary notes not yet seen (at all) by the primary
	primary := app.Map(notes)
	for k, note := range updates {
		if _, ok := primary[k]; !ok {
			note.Content = ""
			notes = append(notes, note)
			continue
		}
	}

	// Represent any note deletions not yet consumed by the primary
	for i, note := range notes {
		if secondaryNote, ok := updates[note.UID]; ok {
			notes[i].Deleted = secondaryNote.Deleted
		}
	}

	sort.Sort(sort.Reverse(app.TimeSorter(notes)))
	return notes
}

func (db *BoltDB) update(note app.Note) (app.Note, error) {
	note, err := app.Persistable(note)
	if err != nil {
		return note, err
	}
	if db.Secondary != nil {
		return db.Secondary.update(note)
	}
	b, err := note.ToBytes()
	if err != nil {
		return note, errors.Wrap(err, "Aborted prior to persist attempt")
	}
	err = db.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		bucket.Put([]byte(note.UID), b)
		return nil
	})
	if err != nil {
		return note, err
	}
	return note, indexNote(note)
}
