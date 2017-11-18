package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

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
