package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionInfo(t *testing.T) {
	buildArch = "ba"
	buildBranch = "bb"
	buildCompiler = "bc"
	buildDate = "bd"
	buildHash = "bh"
	buildUser = "bu"
	buildVersion = "bv"
	v := getVersionInfo()
	assert.Equal(t, "ba", v.Arch)
	assert.Equal(t, "bb", v.Branch)
	assert.Equal(t, "bc", v.Compiler)
	assert.Equal(t, "bd", v.Date)
	assert.Equal(t, "bh", v.Hash)
	assert.Equal(t, "bu", v.User)
	assert.Equal(t, "bv", v.Version)
}

func TestGetVersionInfoStringer(t *testing.T) {
	buildArch = "ba"
	buildBranch = "bb"
	buildCompiler = "bc"
	buildDate = "bd"
	buildHash = "bh"
	buildUser = "bu"
	buildVersion = "bv"
	s := getVersionInfo().String()
	assert.Regexp(t, `Version:\s+bv`, s)
	assert.Regexp(t, `Build time:\s+bd`, s)
	assert.Regexp(t, `Build:\s+bu@bb:bh`, s)
	assert.Regexp(t, `Compiler:\s+bc`, s)
	assert.Regexp(t, `Arch:\s+ba`, s)

}
