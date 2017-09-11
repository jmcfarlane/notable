package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/jmcfarlane/notable/app"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	readSecondaryFile  = ioutil.ReadFile
	writeSecondaryFile = ioutil.WriteFile
)

var uidFromSecondaryPath = regexp.MustCompile(`secondary\.([-a-z0-9]+)`)

// Secondary backend
type Secondary struct {
	Path string
}

func uidFromSecondary(note app.Note, path string) (string, error) {
	if note.UID != "" {
		return note.UID, nil
	}
	match := uidFromSecondaryPath.FindStringSubmatch(path)
	if len(match) != 2 {
		return "", errors.Errorf("Unable to extract UID from secondary=%s", path)
	}
	return match[1], nil
}

func (db *Secondary) deleteByUID(uid string) error {
	path := fmt.Sprintf("%s.secondary.%s.%d", db.Path, uid, time.Now().UnixNano())
	err := writeSecondaryFile(path, make([]byte, 0), 0644)
	log.Infof("Persist (delete) via secondary path=%s err=%v", path, err)
	return err
}

func (db *Secondary) list() app.Notes {
	var notes app.Notes
	secondaryPaths, err := filepath.Glob(fmt.Sprintf("%s.secondary.*.*", db.Path))
	if err != nil {
		log.Errorf("Error trying to find secondary files err=%v", err)
	}
	for _, path := range secondaryPaths {
		var note app.Note
		v, err := readSecondaryFile(path)
		if err != nil {
			log.Errorf("Error trying to reads secondary file err=%v", err)
			continue
		}
		err = note.FromBytes(v)
		if err != nil {
			log.Errorf("Error trying to consume secondary file err=%v", err)
			continue
		}
		note.AheadOfPrimary = true
		note.SecondaryPath = path
		uid, err := uidFromSecondary(note, path)
		if err == nil {
			note.UID = uid
		}
		notes = append(notes, note)
	}
	// Notes are assumed to be unsorted. It's up to the caller to sort
	// according to their needs.
	return notes
}

func (db *Secondary) update(note app.Note) (app.Note, error) {
	note, err := app.Persistable(note)
	if err != nil {
		return note, err
	}
	b, err := note.ToBytes()
	if err != nil {
		return note, errors.Wrap(err, "Aborted prior to persist attempt")
	}
	path := fmt.Sprintf("%s.secondary.%s.%d", db.Path, note.UID, note.Time.UnixNano())
	err = writeSecondaryFile(path, b, 0644)
	log.Infof("Persist via secondary path=%s err=%v", path, err)
	return note, err
}

func reloadAsNeeded() {
	var last time.Time
	for _ = range time.NewTicker(time.Second * 2).C {
		fi, err := os.Stat(*dbPath)
		if err != nil {
			log.Errorf("Unable to stat path=%s err=%v", *dbPath, err)
			continue
		}
		mtime := fi.ModTime()
		if !last.IsZero() && mtime.After(last) {
			db, err = openBoltDB(*dbPath, *secondary)
			log.Infof("Database reloaded due to upstream change err=%v", err)
		}
		last = mtime
	}
}
