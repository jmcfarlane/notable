package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPidHandler(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	resp, err := http.Get(mock.server.URL + "/pid")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	pid, err := strconv.Atoi(string(body))
	assert.Nil(t, err, fmt.Sprintf("Not a valid pid: %s", string(body)))
	assert.True(t, pid > 1024)
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
