package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmcfarlane/notable/api"
	"github.com/julienschmidt/httprouter"
)

// Index returns the landing page html (the application only has one
// page).
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	f, _ := ioutil.ReadFile("../notable/static/templates/index.html")
	fmt.Fprint(w, string(f))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/notes/list", api.Search)
	router.POST("/api/note/content/:uid", api.GetContent)
	router.NotFound = http.FileServer(http.Dir("../notable"))
	log.Fatal(http.ListenAndServe(":8080", router))
}
