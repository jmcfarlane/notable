package flags

import (
	"flag"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

var (
	// Port to listen on
	Port = flag.Int("port", 8080, "Interface and port to listen on")
	// DBPath is assumed to be in your home directory unless specified
	DBPath = flag.String("db", "", "File system path to db file")
	// Browser shoudl be opened
	Browser = flag.Bool("browser", true, "Open a web browser")
)

func homeDirPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Panic("Unable to determine user home directory")
	}
	return usr.HomeDir
}

func init() {
	flag.Parse()
	if *DBPath == "" {
		*DBPath = filepath.Join(homeDirPath(), ".notable/notes.sqlite3")
	}
	log.Infof("Database path=%s", *DBPath)
}
