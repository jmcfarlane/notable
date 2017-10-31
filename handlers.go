package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmcfarlane/notable/app"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"

	log "github.com/sirupsen/logrus"
)

// CreateNote creates a new note
func createNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := app.Note{}
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
	outcome := app.APIResponse{}
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
	if err != nil || app.SmellsEncrypted(content) == true {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Nope, try again")
		return
	}
	w.Write([]byte(content))
}

func adminHandler(m *messenger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			ch := m.add()
			log.Infof("Registered client websocket=%v", ws)
			go func(m *messenger) {
				for {
					var data []byte
					_, err := ws.Read(data)
					if err != nil {
						m.close(ch)
						return
					}
					log.Warnf("Websocket unexpectedly sent data=%s", string(data))
				}
			}(m)
			for msg := range ch {
				log.Infof("Sending frontend push msg=%s", msg)
				ws.Write([]byte(msg))
			}
			log.Info("Unregistered client websocket")
		}).ServeHTTP(w, r)
	}
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
	note := app.Note{}
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
	type v struct {
		Arch     string
		Branch   string
		Compiler string
		Date     string
		Hash     string
		Uptime   string
		User     string
		Version  string
	}
	json.NewEncoder(w).Encode(v{
		Arch:     buildArch,
		Branch:   buildBranch,
		Compiler: buildCompiler,
		Hash:     buildHash,
		Date:     buildStamp,
		User:     buildUser,
		Version:  buildVersion,
		Uptime:   time.Since(booted).String(),
	})
}
