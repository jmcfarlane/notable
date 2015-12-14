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

// APIResponse envelope to communicate details to the frontent
type APIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CreateNote creates a new note
func CreateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := database.Note{}
	json.Unmarshal(payload, &note)
	_, err := database.Create(note)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	noteJSON, err := note.ToJSON()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, noteJSON)
}

// DeleteNote removes a note from storage
func DeleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	outcome := APIResponse{}
	err := database.DeleteByUID(ps.ByName("uid"))
	if err != nil {
		outcome.Success = false
		outcome.Message = err.Error()
		log.Error(err)
	}
	outcome.Success = true
	outcomeJSON, _ := json.Marshal(outcome)
	fmt.Fprintf(w, string(outcomeJSON))
}

// GetContent fetches note content from the database by it's uid.
func GetContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	password := r.PostFormValue("password")
	content, err := database.GetContentByUID(ps.ByName("uid"), password)
	if database.SmellsEncrypted(content) == true {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Nope, try again")
		return
	}
	if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "ERROR")
	}
	fmt.Fprintf(w, content)
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
	noteJSON, err := note.ToJSON()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, noteJSON)
}
