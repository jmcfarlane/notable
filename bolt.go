package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"

	log "github.com/sirupsen/logrus"
)

// BoltDB backend
type BoltDB struct {
	Engine      *bbolt.DB
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
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return nil, err
	}
	dbExisted := pathExists(path)
	engine, err := bbolt.Open(path, 0600, &bbolt.Options{
		ReadOnly: secondary,
		Timeout:  *boltTimeout,
	})
	if err != nil {
		return db, err
	}
	db.Engine = engine
	if !secondary && !dbExisted {
		db.createSchema()
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

func (db *BoltDB) create(note Note) (Note, error) {
	return db.update(note)
}

func (db *BoltDB) createSchema() {
	err := db.Engine.Update(func(tx *bbolt.Tx) error {
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
	err := db.Engine.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		return bucket.Delete([]byte(uid))
	})
	return err
}

func (db *BoltDB) getNoteByUID(uid string, password string) (Note, error) {
	if db.Secondary != nil {
		notes := db.Secondary.list()
		// Sort in reverse order, so the FIRST note wins.
		sort.Sort(sort.Reverse(TimeSorter(notes)))
		for _, note := range notes {
			if note.UID == uid {
				return decryptNote(note, password)
			}
		}
	}
	var note Note
	err := db.Engine.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.NotesBucket)
		v := b.Get([]byte(uid))
		return note.FromBytes(v)
	})
	if err != nil {
		return note, err
	}
	return decryptNote(note, password)
}

func (db *BoltDB) list() Notes {
	var notes Notes
	var updates map[string]Note
	if db.Secondary != nil {
		updatedNotes := db.Secondary.list()
		// Sort in ascending order, so the LAST note wins.
		sort.Sort(TimeSorter(updatedNotes))
		updates = Map(updatedNotes)
	}
	db.Engine.View(func(tx *bbolt.Tx) error {
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

	// Add any secondary notes not yet seen (at all) by the primary
	primary := Map(notes)
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

	sort.Sort(sort.Reverse(TimeSorter(notes)))
	return notes
}

func (db *BoltDB) update(note Note) (Note, error) {
	note, err := Persistable(note)
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
	err = db.Engine.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.NotesBucket)
		bucket.Put([]byte(note.UID), b)
		return nil
	})
	if err != nil {
		return note, err
	}
	return note, indexNote(note)
}
