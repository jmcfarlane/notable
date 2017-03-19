package main

//go:generate go-bindata-assetfs -modtime=1257894000 static/...

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	homedir "github.com/mitchellh/go-homedir"

	log "github.com/Sirupsen/logrus"
)

// Program version information
var (
	buildarch     string
	buildcompiler string
	buildhash     string
	buildstamp    string
	builduser     string
	buildversion  string
)

// This is the application itself
var router = getRouter()
var booted = time.Now()
var db Backend
var err error

// Flags
var (
	bind   = flag.String("bind", "localhost", "Bind address")
	dbPath = flag.String("db", "", "File system path to db file")

	port = flag.Int("port", 8080, "Interface and port to listen on")

	browser = flag.Bool("browser", true, "Open a web browser")
	daemon  = flag.Bool("daemon", true, "Run as a daemon")
	restart = flag.Bool("restart", false, "Restart if already running")
	useBolt = flag.Bool("use.bolt", false, "Use the new BoltDB backend")
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
	h, err := homedir.Expand("~/")
	if err != nil {
		log.Panic("Unable to determine user home directory")
	}
	return h
}

func withoutCaching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func getRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/pid", pid)
	router.GET("/api/notes/list", searchHandler)
	router.GET("/api/version", versionHandler)
	router.POST("/api/note/content/:uid", getContent)
	router.POST("/api/note/create", createNote)
	router.DELETE("/api/note/:uid", deleteNote)
	router.PUT("/api/note/:uid", updateNote)
	router.NotFound = withoutCaching(http.FileServer(assetFS()))
	return router
}

func init() {
	flag.Parse()
	if *dbPath == "" {
		*dbPath = filepath.Join(homeDirPath(), ".notable/notes.db")
	}
}

func main() {
	if *version {
		fmt.Printf("Version:\t%s\n", buildversion)
		fmt.Printf("Build time:\t%s\n", buildstamp)
		fmt.Printf("Build user:\t%s@%s\n", builduser, buildhash)
		fmt.Printf("Compiler:\t%s\n", buildcompiler)
		fmt.Printf("Arch:\t\t%s\n", buildarch)
		return
	}
	if *useBolt || runtime.GOOS == "darwin" {
		db, err = NewBoltDB(*dbPath)
	} else {
		db, err = NewSqlite3(*dbPath)
	}
	defer db.close()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Using backend %s", db)
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
