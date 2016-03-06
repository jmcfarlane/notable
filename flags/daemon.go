package flags

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Daemonize (please use something like upstart, daemontools, systemd)
func Daemonize() bool {
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
	args = append(args, "-daemon=false")
	cmd := exec.Command(os.Args[0], args...)
	cmd.Start()
	log.Infof("Started pid=%v", cmd.Process.Pid)
	os.Exit(0)
	return false
}

func running() bool {
	resp, err := http.Get("http://localhost:" + strconv.Itoa(*Port) + "/pid")
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
	return true
}
