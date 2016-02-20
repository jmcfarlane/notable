package main

//go:generate go-bindata-assetfs static/...

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jmcfarlane/notable/api"
	"github.com/julienschmidt/httprouter"
)

var (
	listen = flag.String("listen", ":8080", "Interface and port to listen on")
)

// Index the landing page html (the application only has one page.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	asset, err := Asset("static/templates/index.html")
	if err != nil {
		log.Panic("Unable to read file from bindata: ", err)
	}
	fmt.Fprint(w, string(asset))
}

func main() {
	flag.Parse()
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/notes/list", api.Search)
	router.POST("/api/note/content/:uid", api.GetContent)
	router.POST("/api/note/create", api.CreateNote)
	router.DELETE("/api/note/:uid", api.DeleteNote)
	router.PUT("/api/note/:uid", api.UpdateNote)
	router.NotFound = http.FileServer(assetFS())
	fmt.Println("Listening on ", *listen)
	log.Fatal(http.ListenAndServe(*listen, router))
}
