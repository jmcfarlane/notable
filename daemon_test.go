package main

import (
	"net"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunning(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	u, err := url.Parse(mock.server.URL)
	assert.Nil(t, err)
	h, p, err := net.SplitHostPort(u.Host)
	assert.Nil(t, err)
	pInt, err := strconv.Atoi(p)
	assert.Nil(t, err)
	*port = pInt
	*bind = h
	assert.True(t, running())
}

func TestDaemonizeCmd(t *testing.T) {
	p := "/usr/local/bin/notable"
	args := []string{p}
	name, out := daemonizeCmd(args)
	assert.Equal(t, p, name)
	assert.Equal(t, []string{"-browser=false", "-daemon=false"}, out)
}

func TestDaemonizeCmdWithBrowser(t *testing.T) {
	p := "/usr/local/bin/notable"
	args := []string{p, "-browser=true"}
	name, out := daemonizeCmd(args)
	assert.Equal(t, p, name)
	assert.Equal(t, []string{"-browser=false", "-daemon=false"}, out)
}

func TestDaemonizeCmdWithPort(t *testing.T) {
	p := "/usr/local/bin/notable"
	args := []string{p, "-port=8000"}
	name, out := daemonizeCmd(args)
	assert.Equal(t, p, name)
	assert.Equal(t, []string{"-port=8000", "-browser=false", "-daemon=false"}, out)
}

func TestDaemonizeCmdWantingToDaemonize(t *testing.T) {
	p := "/usr/local/bin/notable"
	args := []string{p, "-daemon=true"}
	name, out := daemonizeCmd(args)
	assert.Equal(t, p, name)
	assert.Equal(t, []string{"-browser=false", "-daemon=false"}, out)

	args = []string{p, "-daemon=true", "-browser=true"}
	name, out = daemonizeCmd(args)
	assert.Equal(t, p, name)
	assert.Equal(t, []string{"-browser=false", "-daemon=false"}, out)
}

func TestDaemonizeCannotRecurse(t *testing.T) {
	p := "/usr/local/bin/notable"
	// If the args to the running process asked to run in the
	// foreground, nothing should even attempt to daemonize it. If
	// something accidentally does... panic.
	args := []string{p, "-daemon=false"}
	assert.Panics(t, func() { daemonizeCmd(args) })
}
