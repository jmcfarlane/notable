package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/websocket"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// CreateNote creates a new note
func createNote(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)
	note := Note{}
	json.Unmarshal(payload, &note)
	note, err := db.create(note)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	noteJSON, err := note.ToJSON()
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	fmt.Fprint(w, noteJSON)
}

// The current process id
func pid(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, strconv.Itoa(os.Getpid()))
}

// Remove a note from storage
func deleteNote(w http.ResponseWriter, r *http.Request) {
	outcome := APIResponse{}
	err := db.deleteByUID(chi.URLParam(r, "uid"))
	if err != nil {
		outcome.Success = false
		outcome.Message = err.Error()
		log.Error(err)
	}
	outcome.Success = true
	outcomeJSON, _ := json.Marshal(outcome)
	fmt.Fprint(w, string(outcomeJSON))
}

// Fetch note content from the database by it's uid.
func getContent(w http.ResponseWriter, r *http.Request) {
	password := r.PostFormValue("password")
	content, err := getContentByUID(db, chi.URLParam(r, "uid"), password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Nope, try again")
		return
	}
	w.Write([]byte(content))
}

func adminHandler(m *messenger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			ch := m.add()
			log.Infof("Registered client websocket=%v", ws)
			go func(m *messenger) {
				for {
					data := make([]byte, 16)
					_, err := ws.Read(data)
					if err != nil {
						m.close(ch)
						return
					}
					log.Warnf("Websocket unexpectedly sent data=%s", string(bytes.Trim(data, "\x00")))
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

func listHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(db.list())
}

func restartHandler(service *messenger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			http.Error(w, "msg required", http.StatusBadRequest)
			return
		}
		service.send(msg)
		w.Write([]byte("ok"))
	}
}

func stopHandler(service *messenger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.send("")
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("goodbye"))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	uids, err := searchIndex(r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uids)
}

// Persist the updated note to storage
func updateNote(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	note := Note{}
	if err := json.Unmarshal(payload, &note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	note, err = db.update(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func reIndexHandler(w http.ResponseWriter, r *http.Request) {
	i, err := reIndex(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(struct {
		Err   error
		Count int
	}{
		Err:   err,
		Count: i,
	})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	vi := getVersionInfo()
	vi.Pid = os.Getpid()
	vi.Uptime = time.Since(booted).String()
	json.NewEncoder(w).Encode(vi)
}
