package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestIndexHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, _ := http.Get(mock.server.URL + "/")
	body, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "Notable"))
	assert.True(t, strings.Contains(string(body), "/lib/requirejs/require.js"))
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response code != 200")
}

func TestStaticHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/js/main.js")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "The notable client side application"))
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response code != 200")
}

func TestPidHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/pid")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	_, err = strconv.Atoi(string(body))
	assert.Nil(t, err, fmt.Sprintf("Not a valid pid: %s", string(body)))
}

func TestVersionHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/api/version")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	// The version handling is a bit all over the place, stub for now
	assert.Contains(t, string(body), "Uptime")
}

func TestSearchHandlerWithSyntaxErrorInQuery(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/api/notes/search?q=-=+abc")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "syntax error", strings.TrimSuffix(string(body), "\n"))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateHandlerWithNonJSONInput(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	client := &http.Client{}
	p := strings.NewReader("this is not json")
	req, err := http.NewRequest("PUT", mock.server.URL+"/api/note/abc123", p)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(body), "invalid character")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateHandlerWithInvalidInput(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	client := &http.Client{}
	p := strings.NewReader(`{"subject": 1234}`)
	req, err := http.NewRequest("PUT", mock.server.URL+"/api/note/abc123", p)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(body), "into Go struct field Note.subject of type string")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAdminHandler(t *testing.T) {
	frontend := new(messenger)
	mux := httprouter.New()
	mux.GET("/admin", adminHandler(frontend))
	server := httptest.NewServer(mux)
	defer server.Close()
	msgs := []string{"test", "foo", "bar", "stop"}
	u, err := url.Parse(server.URL + "/admin")
	u.Scheme = "ws"
	assert.Nil(t, err)
	ws, err := websocket.Dial(u.String(), "tcp", server.URL)
	assert.Nil(t, err)
	go func() {
		time.Sleep(time.Second * 2)
		for _, msg := range msgs {
			frontend.send(msg)
		}
	}()
	for _, msg := range msgs {
		var data = make([]byte, 16)
		b, err := ws.Read(data)
		assert.Nil(t, err)
		assert.True(t, b > 0)
		assert.Equal(t, msg, string(bytes.Trim(data, "\x00")))
	}
}

func TestAdminHandlerHandlesWrite(t *testing.T) {
	frontend := new(messenger)
	mux := httprouter.New()
	mux.GET("/admin", adminHandler(frontend))
	server := httptest.NewServer(mux)
	defer server.Close()
	msgs := []string{"test", "foo", "bar", "stop"}
	u, err := url.Parse(server.URL + "/admin")
	u.Scheme = "ws"
	assert.Nil(t, err)
	ws, err := websocket.Dial(u.String(), "tcp", server.URL)
	assert.Nil(t, err)
	// Note: Current implementation makes no attempt to parse the 16
	// byte frames. The example payload here just happens to land in
	// two frames ;)
	b, err := ws.Write([]byte("<testing client><write>"))
	assert.Nil(t, err)
	assert.True(t, b > 0)
	go func() {
		time.Sleep(time.Second * 2)
		for _, msg := range msgs {
			frontend.send(msg)
		}
	}()
	for _, msg := range msgs {
		var data = make([]byte, 16)
		b, err := ws.Read(data)
		assert.Nil(t, err)
		assert.True(t, b > 0)
		assert.Equal(t, msg, string(bytes.Trim(data, "\x00")))
	}
}
