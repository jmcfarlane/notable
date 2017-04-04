package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	log "github.com/Sirupsen/logrus"
)

// CreateNote creates a new note
func createNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := Note{}
	json.Unmarshal(payload, &note)
	note, err := db.create(note)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	noteJSON, err := note.ToJSON()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, noteJSON)
}

// The current process id
func pid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, strconv.Itoa(os.Getpid()))
}

// Remove a note from storage
func deleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	outcome := apiResponse{}
	err := db.deleteByUID(ps.ByName("uid"))
	if err != nil {
		outcome.Success = false
		outcome.Message = err.Error()
		log.Error(err)
	}
	outcome.Success = true
	outcomeJSON, _ := json.Marshal(outcome)
	fmt.Fprintf(w, string(outcomeJSON))
}

// Fetch note content from the database by it's uid.
func getContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	password := r.PostFormValue("password")
	content, err := getContentByUID(db, ps.ByName("uid"), password)
	if smellsEncrypted(content) == true {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Nope, try again")
		return
	}
	if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "ERROR")
	}
	w.Write([]byte(content))
}

func listHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	encoder := json.NewEncoder(w)
	encoder.Encode(db.list())
}

func restartHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	restartChan <- r.URL.Query().Get("msg")
	w.Write([]byte("ok"))
}

func searchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uids, err := searchIndex(r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uids)
}

// Persist the updated note to storage
func updateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := Note{}
	json.Unmarshal(payload, &note)
	note, err := db.update(note)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	noteJSON, err := note.ToJSON()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, noteJSON)
}

func versionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(struct {
		Arch     string
		Compiler string
		Date     string
		Hash     string
		Uptime   string
		User     string
		Version  string
	}{
		buildarch,
		buildcompiler,
		buildstamp,
		buildhash,
		time.Since(booted).String(),
		builduser,
		buildversion,
	})
}
