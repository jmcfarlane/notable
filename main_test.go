package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jmcfarlane/notable/app"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

const (
	applicationJSON = "application/json"
	urlEncoded      = "application/x-www-form-urlencoded"
)

type Mock struct {
	db     Backend
	server *httptest.Server
}

func createTestNote(mock Mock, password string) (app.Note, app.Note, int, error) {
	expected := app.Note{
		Content:  "note body beer",
		Password: password,
		Subject:  "test",
		Tags:     "tag1 tag2",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(expected)
	resp, err := http.Post(mock.server.URL+"/api/note/create", applicationJSON, b)
	if err != nil {
		return expected, app.Note{}, 0, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	got := app.Note{}
	json.Unmarshal(content, &got)
	return expected, got, resp.StatusCode, err
}

func setup(t *testing.T) Mock {
	file, err := ioutil.TempFile(os.TempDir(), "notable-testing")
	if !assert.Nil(t, err, "Error creating temp file") {
		return Mock{}
	}

	// Set the `db` global to a test backend
	backend := os.Getenv("BACKEND")
	switch backend {
	case "sqlite3":
		db, err = openSqlite3(file.Name())
	default:
		backend = "boltdb"
		db, err = openBoltDB(file.Name(), false)
	}
	idx, err = getIndex(file.Name() + ".idx")
	fmt.Println("TESTING BACKEND:", backend)
	db.createSchema()
	return Mock{
		db:     db,
		server: httptest.NewServer(router),
	}
}

func tearDown(mock Mock) {
	defer mock.server.Close()
	mock.db.close()
	os.Remove(mock.db.dbFilePath())
	log.Warnf("Deleted temp db path=%s", mock.db.dbFilePath())
}

func TestIndexHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, _ := http.Get(mock.server.URL + "/")
	body, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "Notable"))
	assert.True(t, strings.Contains(string(body), "/lib/requirejs/require.js"))
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response code != 200")
}

func TestNoteCreation(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	expected, got, code, err := createTestNote(mock, "")
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, "", got.Content, "Content should be empty, is lazy loaded")
	assert.Equal(t, expected.Subject, got.Subject, "Subject should match")
	assert.Equal(t, expected.Tags, got.Tags, "Tags should match")
	assert.False(t, got.Encrypted, "Should not be encrypted, no password")
	assert.Equal(t, http.StatusOK, code, "Response code != 200")
}

func TestNoteCreationContentFetch(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	expected, got, code, err := createTestNote(mock, "")
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, http.StatusOK, code, "Response code != 200")
	resp, err := http.Post(mock.server.URL+"/api/note/content/"+got.UID, "", nil)
	assert.Nil(t, err, "Should be no http error")
	content, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expected.Content, string(content), "Did not get the content back")
}

func TestNoteCreationContentFetchGet(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/api/note/content/abc123")
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Expected: Method Not Allowed")
}

func TestEncryptedNoteCreationContentFetch(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	password := "fancy-password"
	expected, got, code, err := createTestNote(mock, password)
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, http.StatusOK, code, "Response code != 200")
	form := url.Values{}
	form.Add("password", password)
	body := strings.NewReader(form.Encode())
	resp, err := http.Post(mock.server.URL+"/api/note/content/"+got.UID, urlEncoded, body)
	assert.Nil(t, err, "Should be no http error")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response code != 200")
	assert.Equal(t, expected.Content, string(content), "Did not get the content back")
}

func TestEncryptedNoteCreationExcludesPassword(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	password := "fancy-password"
	_, got, _, _ := createTestNote(mock, password)
	persistedNote, err := db.getNoteByUID(got.UID, password)
	assert.Nil(t, err, "Should be no error fetching the note")
	assert.Equal(t, "", persistedNote.Password, "Passwords should not be persisted!")
}

func TestEncryptedNoteCreationContentFetchWithWrongPassword(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	expected, got, code, err := createTestNote(mock, "foobar")
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, http.StatusOK, code, "Response code != 200")
	form := url.Values{}
	form.Add("password", "wrong password!")
	body := strings.NewReader(form.Encode())
	resp, err := http.Post(mock.server.URL+"/api/note/content/"+got.UID, urlEncoded, body)
	assert.Nil(t, err, "Should be no http error")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Request should be forbidden!")
	assert.NotEqual(t, expected.Content, string(content), "Got content back?")
}

func TestNoteDeletion(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, got, code, err := createTestNote(mock, "")
	assert.Nil(t, err, "Should be no http error")
	assert.Equal(t, http.StatusOK, code, "Response code != 200")
	req, err := http.NewRequest("DELETE", mock.server.URL+"/api/note/"+got.UID, nil)
	assert.Nil(t, err, "Should be no error creating a new request")
	resp, err := http.DefaultClient.Do(req)
	content, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be no error reading the response")
	outcome := app.APIResponse{}
	err = json.Unmarshal(content, &outcome)
	assert.Nil(t, err, "Should be no error unmarshaling the response")
	assert.Equal(t, true, outcome.Success)
	resp, err = http.Get(mock.server.URL + "/api/notes/list")
	assert.Nil(t, err, "Should be no http error")
	content, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be not be a response error")
	notes := []app.Note{}
	json.Unmarshal(content, &notes)
	assert.Equal(t, 0, len(notes), "Should be no notes as the only one was deleted")
}

func TestNoteListing(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	expected, _, _, _ := createTestNote(mock, "")
	resp, err := http.Get(mock.server.URL + "/api/notes/list")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be no http error")
	notes := []app.Note{}
	json.Unmarshal(content, &notes)
	assert.Equal(t, expected.Subject, notes[0].Subject, "Listing miissing our note")
}

func TestNoteFullTextSearch(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, got, _, _ := createTestNote(mock, "")
	resp, err := http.Get(mock.server.URL + "/api/notes/search?q=beer")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be no http error")
	uids := []string{}
	json.Unmarshal(content, &uids)
	if assert.Equal(t, 1, len(uids), "Full text should find 1 note") {
		assert.Equal(t, got.UID, uids[0], "Full text found the note")
	}
}

func TestNoteFullTextSearchDeletion(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, got, _, _ := createTestNote(mock, "")
	req, err := http.NewRequest("DELETE", mock.server.URL+"/api/note/"+got.UID, nil)
	assert.Nil(t, err, "Should be no error creating a new request")
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err, "Should be no error requesting a note deletion")
	resp, err = http.Get(mock.server.URL + "/api/notes/search?q=beer")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be no http error")
	uids := []string{}
	json.Unmarshal(content, &uids)
	assert.Equal(t, 0, len(uids), "Full text should find 0 notes")
}

func TestNoteFullTextSearchUpdate(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, got, _, _ := createTestNote(mock, "")
	got.Content = "ipa only pls"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(got)
	req, err := http.NewRequest("PUT", mock.server.URL+"/api/note/"+got.UID, b)
	assert.Nil(t, err, "Should be no error creating the http request")
	_, err = http.DefaultClient.Do(req)
	assert.Nil(t, err, "Should be no error requesting a note update")
	resp, err := http.Get(mock.server.URL + "/api/notes/search?q=ipa")
	content, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be no http error")
	uids := []string{}
	json.Unmarshal(content, &uids)
	if assert.Equal(t, 1, len(uids), "Full text should find 1 note") {
		assert.Equal(t, got.UID, uids[0], "Full text found the note")
	}
}
