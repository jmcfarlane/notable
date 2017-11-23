package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Force the db path to be set to a safe place, we don't wanna
	// blow away data just cuz we're testing right?
	*dbPath = filepath.Join(os.TempDir(), "notable-testing/notes.db")
	os.Exit(m.Run())
}

func TestMainFlagVersion(t *testing.T) {
	v := *version
	bv := buildVersion
	defer func() {
		buildVersion = bv
		*version = v
	}()
	*version = true
	buildVersion = fmt.Sprintf("test verison %s", time.Now())
	b := new(bytes.Buffer)
	run(b)
	assert.Contains(t, b.String(), buildVersion)
}

func TestRunWhenAlreadyRunning(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	b := new(bytes.Buffer)
	u, err := url.Parse(mock.server.URL)
	assert.Nil(t, err)
	p, err := strconv.Atoi(u.Port())
	*port = p
	*browser = false
	*daemon = false
	start := time.Now()
	run(b)
	assert.True(t, time.Since(start).Seconds() < 0.1)
	assert.Equal(t, "", b.String())
}

func TestMainStartStop(t *testing.T) {
	b := *browser
	d := *daemon
	p := *dbPath
	defer func() {
		*browser = b
		*daemon = d
		*dbPath = p
	}()
	*daemon = false
	*dbPath = "/tmp/notable-testing/notes.db"
	*browser = false
	url := fmt.Sprintf("http://localhost:%d/api/stop", *port)
	go func() {
		time.Sleep(time.Second * 2)
		req, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			log.Errorf("Unable to create http PUT request url=%v err=%s", url, err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("Unable to request stop url=%v err=%s", url, err)
			return
		}
		resp.Body.Close()
	}()
	run(new(bytes.Buffer))
	t.Logf("^^ foreground service unblocked, so PUT to url=%q worked!", url)
}

func TestHomeDir(t *testing.T) {
	assert.NotEqual(t, "", homeDirPath("~/"))
}

func TestHomeDirError(t *testing.T) {
	assert.Panics(t, func() { homeDirPath("~proc") })
}

func TestBrowserCmdDarwin(t *testing.T) {
	name, args := browserCmd("darwin")
	assert.Equal(t, "open", name)
	assert.Equal(t, []string{fmt.Sprintf("http://%s:%d", *bind, *port)}, args)
}

func TestBrowserCmdDefautl(t *testing.T) {
	name, args := browserCmd("?")
	assert.Equal(t, "xdg-open", name)
	assert.Equal(t, []string{fmt.Sprintf("http://%s:%d", *bind, *port)}, args)
}

func TestBrowserCmdLinux(t *testing.T) {
	name, args := browserCmd("linux")
	assert.Equal(t, "xdg-open", name)
	assert.Equal(t, []string{fmt.Sprintf("http://%s:%d", *bind, *port)}, args)
}

func TestBrowserCmdWindows(t *testing.T) {
	name, args := browserCmd("windows")
	assert.Equal(t, "cmd", name)
	assert.Equal(t, []string{`/c`, "start", fmt.Sprintf("http://%s:%d", *bind, *port)}, args)
}
