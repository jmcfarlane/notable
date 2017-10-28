// https://git.tcc.li/projects/OE/repos/go-quickstart

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
	buildHash     string
	buildStamp    string
	buildUser     string
	buildVersion  string
)

type versionInfo struct {
	buildArch     string
	buildBranch   string
	buildCompiler string
	buildHash     string
	buildStamp    string
	buildUser     string
	buildVersion  string
}

func (vi versionInfo) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Version:\t%s\n", vi.buildVersion))
	buffer.WriteString(fmt.Sprintf("Build time:\t%s\n", vi.buildStamp))
	buffer.WriteString(fmt.Sprintf("Build:\t\t%s@%s:%s\n", vi.buildUser, vi.buildBranch, vi.buildHash))
	buffer.WriteString(fmt.Sprintf("Compiler:\t%s\n", vi.buildCompiler))
	buffer.WriteString(fmt.Sprintf("Arch:\t\t%s\n", vi.buildArch))
	return buffer.String()
}

func getVersionInfo() versionInfo {
	return versionInfo{
		buildArch:     buildArch,
		buildBranch:   buildBranch,
		buildCompiler: buildCompiler,
		buildHash:     buildHash,
		buildStamp:    buildStamp,
		buildUser:     buildUser,
		buildVersion:  buildVersion,
	}
}
