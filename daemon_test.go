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
