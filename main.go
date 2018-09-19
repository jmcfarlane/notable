package main

//go:generate rice embed-go

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	rice "github.com/GeertJohan/go.rice"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

// This is the application itself
var booted = time.Now()
var db Backend
var idx bleve.Index
var box *rice.Box

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
	asset, err := box.String("templates/index.html")
	if err != nil {
		log.Panic("Unable to read file from box: ", err)
	}
	fmt.Fprint(w, string(asset))
}

func browserCmd(goos string) (string, []string) {
	url := fmt.Sprintf("http://%s:%d", *bind, *port)
	switch goos {
	case "darwin":
		return "open", []string{url}
	case "windows":
		return "cmd", []string{"/c", "start", url}
	}
	return "xdg-open", []string{url}
}

func openBrowser() error {
	name, args := browserCmd(runtime.GOOS)
	return exec.Command(name, args...).Run()
}

func start(router *httprouter.Router, service *messenger) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bind, *port))
	if err != nil {
		log.Fatal(err)
	}
	go func(listener net.Listener) {
		serviceCh := service.add()
		msg := <-serviceCh
		if msg == "" {
			if err := listener.Close(); err != nil {
				log.Fatalf("Failed to stop TCP listener: err=%v", err)
			}
			log.Info("TCP listener closed, goodbye!")
			return
		}
		log.Warnf("Restart requested msg=%s", msg)
		listener.Close()
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Start()
		log.Infof("Replacement started pid=%v", cmd.Process.Pid)
		os.Exit(0)
	}(listener)
	log.Infof("Listening on %s:%v pid=%d", *bind, *port, os.Getpid())
	http.Serve(listener, router)
}

func homeDirPath(path string) string {
	h, err := homedir.Expand(path)
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

func (m *messenger) empty() bool {
	m.Lock()
	defer m.Unlock()
	return len(m.clients) == 0

}

func (m *messenger) send(msg string) {
	m.Lock()
	defer m.Unlock()
	for _, client := range m.clients {
		client <- msg
	}
}

func getRouter(frontend, service *messenger) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/pid", pid)
	router.GET("/admin", adminHandler(frontend))
	router.GET("/api/notes/list", listHandler)
	router.GET("/api/notes/search", searchHandler)
	router.GET("/api/version", versionHandler)
	router.POST("/api/note/content/:uid", getContent)
	router.POST("/api/note/create", createNote)
	router.DELETE("/api/note/:uid", deleteNote)
	router.PUT("/api/note/:uid", updateNote)
	router.PUT("/api/restart", restartHandler(service))
	router.PUT("/api/stop", stopHandler(service))
	router.NotFound = withoutCaching(http.FileServer(box.HTTPBox()))
	return router
}

func persistSecondaryUpdate(note Note) error {
	if note.Deleted {
		return db.deleteByUID(note.UID)
	}
	note.AheadOfPrimary = false
	_, err := db.update(note)
	return err
}

func consumeSecondaries(db Backend, secondaries Secondary, frontend *messenger) {
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
		frontend.send("reload")
	}
}

func run(w io.Writer) {
	if *version {
		fmt.Fprintln(w, getVersionInfo())
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
		*dbPath = filepath.Join(homeDirPath("~/"), ".notable/notes.db")
	}
	db, err = openBoltDB(*dbPath, *secondary)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to open database"))
	}
	idx, err = getIndex(db.dbFilePath() + ".idx")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to establish search index"))
	}
	if *doReIndex {
		err = reIndex(db)
		if err != nil {
			log.Panic("Re-indexing failed:", err)
		}
	}
	backend := new(messenger)
	frontend := new(messenger)
	service := new(messenger)
	if *secondary {
		go reloadAsNeeded(db, frontend, backend)
	} else {
		secondaries := Secondary{Path: db.dbFilePath()}
		consumeSecondaries(db, secondaries, frontend)
		go func() {
			stopCh := backend.add()
			for range time.NewTicker(time.Second * 2).C {
				select {
				case <-stopCh:
					backend.close(stopCh)
					log.Info("Consumption of secondary files stopped, goodbye!")
					return
				case <-time.After(time.Millisecond):
					consumeSecondaries(db, secondaries, frontend)
				}
			}
		}()
	}
	defer db.close()
	log.Infof("Using backend %s", db)
	if *daemon {
		daemonize()
	} else {
		start(getRouter(frontend, service), service)
	}
	backend.send("stop")
	for !backend.empty() {
		time.Sleep(time.Millisecond * 10)
	}
	log.Info("Service fully stopped, goodbye!!")
}

func main() {
	flag.Parse()
	box = rice.MustFindBox("static")
	run(os.Stdout)
}
