package main

//go:generate go-bindata-assetfs static/...

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/julienschmidt/httprouter"

	log "github.com/Sirupsen/logrus"
)

// Program version information
var (
	buildarch     string
	buildcompiler string
	buildhash     string
	buildstamp    string
	builduser     string
)

// Flags
var (
	port    = flag.Int("port", 8080, "Interface and port to listen on")
	dbPath  = flag.String("db", "", "File system path to db file")
	browser = flag.Bool("browser", true, "Open a web browser")
	daemon  = flag.Bool("daemon", true, "Run as a daemon")
	restart = flag.Bool("restart", false, "Restart if already running")
	bind    = flag.String("bind", "localhost", "Bind address")
	version = flag.Bool("version", false, "Print program version information")
)

// Index the landing page html (the application only has one page.
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	asset, err := Asset("static/templates/index.html")
	if err != nil {
		log.Panic("Unable to read file from bindata: ", err)
	}
	fmt.Fprint(w, string(asset))
}

func openBrowser() {
	cmd := "open"
	if runtime.GOOS == "linux" {
		cmd = "xdg-open"
	}
	err := exec.Command(cmd, "http://"+*bind+":"+strconv.Itoa(*port)).Run()
	if err != nil {
		log.Errorf("Error spawning web browser err=%v", err)
	}
}

func start(router *httprouter.Router) {
	log.Infof("Listening on %s:%v", *bind, *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}

func homeDirPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Panic("Unable to determine user home directory")
	}
	return usr.HomeDir
}

func init() {
	flag.Parse()
	if *dbPath == "" {
		*dbPath = filepath.Join(homeDirPath(), ".notable/notes.sqlite3")
	}
	log.Infof("Database path=%s", *dbPath)
	createSchema()
}

func main() {
	if *version {
		fmt.Printf("Build time:\t%s\n", buildstamp)
		fmt.Printf("Build user:\t%s@%s\n", builduser, buildhash)
		fmt.Printf("Compiler:\t%s\n", buildcompiler)
		fmt.Printf("Arch:\t\t%s\n", buildarch)
		return
	}
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/pid", pid)
	router.GET("/api/notes/list", searchHandler)
	router.POST("/api/note/content/:uid", getContent)
	router.POST("/api/note/create", createNote)
	router.DELETE("/api/note/:uid", deleteNote)
	router.PUT("/api/note/:uid", updateNote)
	router.NotFound = http.FileServer(assetFS())
	if *browser {
		openBrowser()
	}
	if *daemon {
		if !daemonize() {
			start(router)
		}
	} else {
		start(router)
	}
}
