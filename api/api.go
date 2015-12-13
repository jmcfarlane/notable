package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmcfarlane/notable/database"
	"github.com/julienschmidt/httprouter"

	log "github.com/Sirupsen/logrus"
)

// CreateNote creates a new note
func CreateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := database.Note{}
	json.Unmarshal(payload, &note)
	note, err := database.Create(note)

	// Return the note (minus the content) in case the UI sees any
	// changes (like timestamps or ids for new notes)
	note.Content = ""
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	noteJSON, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		log.Error("Failed to parse note into json", err)
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, string(noteJSON))
}

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

// UpdateNote persists the updated note to storage
func UpdateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := database.Note{}
	json.Unmarshal(payload, &note)
	note, err := database.Update(note)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Return the note (minus the content) in case the UI sees any
	// changes (like timestamps or ids for new notes)
	note.Content = ""
	noteJSON, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		log.Error("Failed to parse note into json", err)
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, string(noteJSON))
}
