package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	*browser = false
	b := new(bytes.Buffer)
	x := b
	run(b)
	assert.NotEmpty(t, b)
	assert.Empty(t, x)
	a := b
	run(b)
	assert.Equal(t, a, b)
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
