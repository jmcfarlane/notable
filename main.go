package main

//go:generate go-bindata-assetfs static/...

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/jmcfarlane/notable/api"
	"github.com/jmcfarlane/notable/flags"
	"github.com/julienschmidt/httprouter"

	log "github.com/Sirupsen/logrus"
)

// Index the landing page html (the application only has one page.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	asset, err := Asset("static/templates/index.html")
	if err != nil {
		log.Panic("Unable to read file from bindata: ", err)
	}
	fmt.Fprint(w, string(asset))
}

func browser() {
	if *flags.Browser {
		cmd := "open"
		if runtime.GOOS == "linux" {
			cmd = "xdg-open"
		}
		err := exec.Command(cmd, "http://localhost:"+strconv.Itoa(*flags.Port)).Run()
		if err != nil {
			log.Errorf("Error spawning web browser err=%v")
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/notes/list", api.Search)
	router.POST("/api/note/content/:uid", api.GetContent)
	router.POST("/api/note/create", api.CreateNote)
	router.DELETE("/api/note/:uid", api.DeleteNote)
	router.PUT("/api/note/:uid", api.UpdateNote)
	router.NotFound = http.FileServer(assetFS())
	browser()
	log.Infof("Listening on localhost:%v", *flags.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*flags.Port), router))
}
