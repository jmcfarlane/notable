package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmcfarlane/notable/api"
	"github.com/julienschmidt/httprouter"
)

// Index the landing page html (the application only has one page.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	f, _ := ioutil.ReadFile("notable/static/templates/index.html")
	fmt.Fprint(w, string(f))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/notes/list", api.Search)
	router.POST("/api/note/content/:uid", api.GetContent)
	router.POST("/api/note/create", api.CreateNote)
	router.PUT("/api/note/:uid", api.UpdateNote)
	router.NotFound = http.FileServer(http.Dir("notable"))
	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
