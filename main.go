package main

//go:generate go-bindata-assetfs -modtime=1257894000 static/...

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/julienschmidt/httprouter"
	homedir "github.com/mitchellh/go-homedir"

	log "github.com/Sirupsen/logrus"
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

// This is the application itself
var router = getRouter()
var booted = time.Now()
var db Backend
var idx bleve.Index

// Support restarts
var restartChan = make(chan string, 1)

// Flags
var (
	bind   = flag.String("bind", "localhost", "Bind address")
	dbPath = flag.String("db", "", "File system path to db file")

	port = flag.Int("port", 8080, "Interface and port to listen on")

	browser   = flag.Bool("browser", true, "Open a web browser")
	daemon    = flag.Bool("daemon", true, "Run as a daemon")
	doReIndex = flag.Bool("reindex", false, "Re-index all notes on startup")
	restart   = flag.Bool("restart", false, "Restart if already running")
	useBolt   = flag.Bool("use.bolt", true, "Use the new BoltDB backend")
	version   = flag.Bool("version", false, "Print program version information")
)

// Index the landing page html (the application only has one page.
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	asset, err := Asset("static/templates/index.html")
	if err != nil {
		log.Panic("Unable to read file from bindata: ", err)
	}
	fmt.Fprint(w, string(asset))
}

func openBrowser() error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	url := "http://" + *bind + ":" + strconv.Itoa(*port)
	return exec.Command(args[0], append(args[1:], url)...).Run()
}

func start(router *httprouter.Router) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bind, *port))
	log.Infof("Listening on %s:%v pid=%d", *bind, *port, os.Getpid())
	if err != nil {
		log.Fatal(err)
	}
	go func(listener net.Listener) {
		log.Warnf("Restart requested msg=%s", <-restartChan)
		listener.Close()
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Start()
		log.Infof("Replacement started pid=%v", cmd.Process.Pid)
		os.Exit(0)
	}(listener)
	http.Serve(listener, router)
	time.Sleep(time.Second * 5)
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
	router.GET("/api/notes/list", listHandler)
	router.GET("/api/notes/search", searchHandler)
	router.GET("/api/version", versionHandler)
	router.POST("/api/note/content/:uid", getContent)
	router.POST("/api/note/create", createNote)
	router.DELETE("/api/note/:uid", deleteNote)
	router.PUT("/api/note/:uid", updateNote)
	router.PUT("/api/restart", restartHandler)
	router.NotFound = withoutCaching(http.FileServer(assetFS()))
	return router
}

func init() {
	flag.Parse()
	if *dbPath == "" {
		*dbPath = filepath.Join(homeDirPath(), ".notable/notes.db")
	}
	var err error
	idx, err = getIndex(*dbPath + ".idx")
	if err != nil {
		log.Fatalf("Unable to establish search index err=%v", err)
	}
}

func main() {
	if *version {
		fmt.Printf("Version:\t%s\n", buildVersion)
		fmt.Printf("Build time:\t%s\n", buildStamp)
		fmt.Printf("Build:\t\t%s@%s:%s\n", buildUser, buildBranch, buildHash)
		fmt.Printf("Compiler:\t%s\n", buildCompiler)
		fmt.Printf("Arch:\t\t%s\n", buildArch)
		return
	}
	if *browser {
		err := openBrowser()
		if err != nil {
			log.Fatalf("Failed to open a browser err=%v", err)
		}
	}
	if running() {
		return
	}
	var err error
	if *useBolt || runtime.GOOS == "darwin" {
		db, err = NewBoltDB(*dbPath)
	} else {
		db, err = NewSqlite3(*dbPath)
	}
	if err != nil {
		log.Fatal(err)
	}
	if *doReIndex {
		err = reIndex(db)
		if err != nil {
			log.Panic("Re-indexing failed:", err)
		}
	}
	defer db.close()
	log.Infof("Using backend %s", db)
	if *daemon {
		daemonize()
	} else {
		start(router)
	}
}
