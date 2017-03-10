package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

// Daemonize (please use something like upstart, daemontools, systemd)
func daemonize() bool {
	if running() {
		return true
	}
	args := os.Args[1:]
	i := 0
	for ; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-daemon") {
			args = append(args[:i], args[i+1:]...)
			break
		}
	}
	args = append(args, "-browser=false")
	args = append(args, "-daemon=false")
	cmd := exec.Command(os.Args[0], args...)
	cmd.Start()
	log.Infof("Started pid=%v", cmd.Process.Pid)
	os.Exit(0)
	return false
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
