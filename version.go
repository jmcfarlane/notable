package main

import (
	"bytes"
	"fmt"
)

// Program version information
var (
	buildArch     string
	buildBranch   string
	buildCompiler string
	buildDate     string
	buildHash     string
	buildUser     string
	buildVersion  string
)

type versionInfo struct {
	Arch     string
	Branch   string
	Compiler string
	Date     string
	Hash     string
	User     string
	Version  string

	// Non static
	Pid    int
	Uptime string
}

func (vi versionInfo) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Version:\t%s\n", vi.Version))
	buffer.WriteString(fmt.Sprintf("Build time:\t%s\n", vi.Date))
	buffer.WriteString(fmt.Sprintf("Build:\t\t%s@%s:%s\n", vi.User, vi.Branch, vi.Hash))
	buffer.WriteString(fmt.Sprintf("Compiler:\t%s\n", vi.Compiler))
	buffer.WriteString(fmt.Sprintf("Arch:\t\t%s\n", vi.Arch))
	return buffer.String()
}

func getVersionInfo() versionInfo {
	return versionInfo{
		Arch:     buildArch,
		Branch:   buildBranch,
		Compiler: buildCompiler,
		Hash:     buildHash,
		Date:     buildDate,
		User:     buildUser,
		Version:  buildVersion,
	}
}
