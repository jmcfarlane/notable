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
	"sync"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/jmcfarlane/notable/app"
	"github.com/julienschmidt/httprouter"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// This is the application itself
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

	browser     = flag.Bool("browser", true, "Open a web browser")
	daemon      = flag.Bool("daemon", true, "Run as a daemon")
	doReIndex   = flag.Bool("reindex", false, "Re-index all notes on startup")
	restart     = flag.Bool("restart", false, "Restart if already running")
	secondary   = flag.Bool("secondary", false, "Run program as secondary, not primary")
	version     = flag.Bool("version", false, "Print program version information")
	boltTimeout = flag.Duration("bolt.timeout", time.Duration(time.Second*2), "Boltdb open timeout")
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

type messenger struct {
	sync.Mutex
	clients []chan string
}

func (m *messenger) add() chan string {
	reloadChan := make(chan string, 1)
	m.Lock()
	defer m.Unlock()
	m.clients = append(m.clients, reloadChan)
	return reloadChan
}

func (m *messenger) close(ch chan string) {
	m.Lock()
	defer m.Unlock()
	for i, client := range m.clients {
		if client != ch {
			continue
		}
		m.clients = append(m.clients[:i], m.clients[i+1:]...)
		close(ch)
	}
}

func (m *messenger) remove(i int) {
	m.Lock()
	defer m.Unlock()
	m.clients = append(m.clients[:i], m.clients[i+1:]...)
}

func (m *messenger) send(msg string) {
	m.Lock()
	defer m.Unlock()
	for _, client := range m.clients {
		client <- msg
	}
}

func getRouter(m *messenger) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/pid", pid)
	router.GET("/admin", adminHandler(m))
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

func persistSecondaryUpdate(note app.Note) error {
	if note.Deleted {
		return db.deleteByUID(note.UID)
	}
	note.AheadOfPrimary = false
	_, err := db.update(note)
	return err
}

func consumeUpdatesFromSecondaries(db Backend, secondaries Secondary, m *messenger) {
	clientReload := false
	for _, note := range secondaries.list() {
		if err := persistSecondaryUpdate(note); err != nil {
			log.Errorf("Unable to recover note=%v err=%v", note, err)
			continue
		}
		if err := os.Remove(note.SecondaryPath); err != nil {
			log.Errorf("Unable to delete secondary note=%v err=%v", note, err)
			continue
		}
		log.Infof("Successfully recovered uid=%s", note.UID)
		clientReload = true
	}
	if clientReload {
		m.send("reload")
	}
}

func main() {
	flag.Parse()
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
			log.Fatal(errors.Wrap(err, "Failed to open a browser"))
		}
	}
	if running() {
		return
	}
	var err error
	if *dbPath == "" {
		*dbPath = filepath.Join(homeDirPath(), ".notable/notes.db")
	}
	db, err = openBoltDB(*dbPath, *secondary)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to open database"))
	}
	idx, err = getIndex(*dbPath + ".idx")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to establish search index"))
	}
	if *doReIndex {
		err = reIndex(db)
		if err != nil {
			log.Panic("Re-indexing failed:", err)
		}
	}
	m := new(messenger)
	if *secondary {
		go reloadAsNeeded(m)
	} else {
		secondaries := Secondary{Path: *dbPath}
		consumeUpdatesFromSecondaries(db, secondaries, m)
		go func() {
			for _ = range time.NewTicker(time.Second * 2).C {
				consumeUpdatesFromSecondaries(db, secondaries, m)
			}
		}()
	}
	defer db.close()
	log.Infof("Using backend %s", db)
	if *daemon {
		daemonize()
	} else {
		start(getRouter(m))
	}
}
