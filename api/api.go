package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmcfarlane/notable/database"
	"github.com/julienschmidt/httprouter"
)

// GetContent fetches note content from the database by it's uid.
func GetContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	password := r.PostFormValue("password")
	body, err := database.GetContentByUID(ps.ByName("uid"), password)
	if err != nil {
		fmt.Fprintf(w, "ERROR")
	}
	fmt.Fprintf(w, body)
}

// Search for notes based on an optional querystring parameter
func Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	notes := database.Search("")
	thing, err := json.MarshalIndent(notes, "", "\t")
	if err != nil {
		fmt.Fprintf(w, "[]")
	}
	fmt.Fprintf(w, string(thing))
}
