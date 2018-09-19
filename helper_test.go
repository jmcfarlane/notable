package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	rice "github.com/GeertJohan/go.rice"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

const (
	applicationJSON = "application/json"
	urlEncoded      = "application/x-www-form-urlencoded"
)

// Mock structs for the purpose fo testing
type Mock struct {
	db        Backend
	secondary Backend
	server    *httptest.Server
}

func createTestNote(mock Mock, password string) (Note, Note, int, error) {
	expected := Note{
		Content:  "note body beer",
		Password: password,
		Subject:  "test",
		Tags:     "tag1 tag2",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(expected)
	resp, err := http.Post(mock.server.URL+"/api/note/create", applicationJSON, b)
	if err != nil {
		return expected, Note{}, 0, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	got := Note{}
	json.Unmarshal(content, &got)
	return expected, got, resp.StatusCode, err
}

func copyFile(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func setup(t *testing.T) Mock {
	box = rice.MustFindBox("static")
	tempDir, err := ioutil.TempDir(os.TempDir(), "notable-testing")
	if !assert.Nil(t, err, "Error creating temp dir") {
		return Mock{}
	}
	db, err = openBoltDB(filepath.Join(tempDir, "notes.db"), false)
	assert.Nil(t, err)
	idx, err = getIndex(db.dbFilePath() + ".idx")
	assert.Nil(t, err)

	// Because the secondary needs to be on separate filesystem that
	// has some sync mechanism (dropbox, syncthing, keybase, etc) we
	// need a copy of the file. Bolt won't allow a secondary readonly
	// against the _same_ file.
	secondaryPath := filepath.Join(tempDir, "secondary.db")
	err = copyFile(db.dbFilePath(), secondaryPath)
	assert.Nil(t, err)

	// Open the secondary (knowing it's name is different)
	dbSecondary, err := openBoltDB(secondaryPath, true)
	assert.Nil(t, err)

	// Now fake the secondary path, so it reads/writes via journal
	// files named like they would in the wild (prefixed by db.Path)
	dbSecondary.Secondary.Path = db.dbFilePath()

	return Mock{
		db:        db,
		secondary: dbSecondary,
		server:    httptest.NewServer(getRouter(new(messenger), new(messenger))),
	}
}

func tearDown(mock Mock) {
	defer mock.server.Close()
	mock.db.close()
	tempDir := filepath.Dir(mock.db.dbFilePath())
	log.Warnf("Deleted temp db dir path=%q err=%v", tempDir, os.RemoveAll(tempDir))
}
