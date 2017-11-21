package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func daemonizeCmd(args []string) (string, []string) {
	name, args := args[0], args[1:]
	filtered := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-browser") {
			continue
		}
		if strings.HasPrefix(arg, "-daemon") {
			if arg == "-daemon=false" {
				panic("Refusing to daemonize as proc asked for foreground")
			}
			continue
		}
		filtered = append(filtered, arg)
	}
	return name, append(filtered, "-browser=false", "-daemon=false")
}

// Daemonize (please use something like upstart, daemontools, systemd)
func daemonize() {
	name, args := daemonizeCmd(os.Args)
	cmd := exec.Command(name, args...)
	cmd.Start()
	log.Infof("Started pid=%v", cmd.Process.Pid)
	os.Exit(0)
}

func running() bool {
	resp, err := http.Get("http://" + *bind + ":" + strconv.Itoa(*port) + "/pid")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	pid, err := strconv.Atoi(string(contents))
	if err != nil {
		return false
	}
	log.Infof("Already running pid=%v", pid)
	if *restart {
		process := os.Process{Pid: pid}
		process.Signal(syscall.SIGINT)
		log.Warnf("Requested graceful shutdown pid=%v", pid)
		return false
	}
	return true
}
