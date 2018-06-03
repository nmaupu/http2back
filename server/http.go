package server

import (
	"encoding/json"
	"fmt"
	"github.com/nmaupu/http2back/notifier"
	"github.com/nmaupu/http2back/provider"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Globals
var (
	getProv         func() provider.Provider
	getNotifs       []func() notifier.Notifier
	maxMemoryBuffer int64
)

// HTTP API
func handleRequest(w http.ResponseWriter, r *http.Request) {
	type JsonResult struct {
		Result string `json:"result"`
	}
	type JsonError struct {
		Error string `json:"error"`
	}
	defer func() {
		if r := recover(); r != nil {
			je := JsonError{fmt.Sprint(r)}
			log.Println(je.Error)
			j, _ := json.Marshal(je)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(j)
		}
	}()

	err := r.ParseForm()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}

	if r.Method == "PUT" || r.Method == "POST" {
		r.ParseMultipartForm(maxMemoryBuffer)

		// Getting extradir parameter
		extradir := r.Form.Get("extradir")
		if extradir != "" {
			extradir = fmt.Sprintf("%s/", extradir)
		}

		// Get file
		in, handler, err := r.FormFile("file")
		if err != nil {
			panic(fmt.Sprintf("Error: %s", err))
		}
		defer in.Close()

		// Unique filename
		t := time.Now()
		name := fmt.Sprintf("%s%s-%06d.%s", extradir, t.Format("20060102_150405"), rand.Intn(100000), handler.Filename)
		ret := getProv().Copy(in, name)

		// Send result
		log.Printf("Called : %s %s - file %s -> %s", r.Method, r.URL.Path, handler.Filename, ret)
		jr := JsonResult{ret}
		j, _ := json.Marshal(jr)
		w.Write(j)

		// Send notifications to all notifiers
		for _, f := range getNotifs {
			notif := f()
			if notif != nil {
				err = notif.Notify(&notifier.Event{
					Title:   fmt.Sprintf("http2back - new file available: %s", ret),
					Message: ret,
				})
				if err != nil {
					log.Printf("Notification has not been sent: %v", err)
				}
			}
		}

	} else {
		panic(fmt.Sprintf("%s is an unsupported method for %s", r.Method, r.URL.Path))
	}
}

// Starting the HTTP server
func Start(port *int, bind_addr *string, maxMemMB *int, getProvider func() provider.Provider, getNotifiers []func() notifier.Notifier) {
	maxMemoryBuffer = int64((*maxMemMB) << 20)
	getProv = getProvider
	getNotifs = getNotifiers

	log.Printf("Starting http server on %s:%d using provider %s - Memory per buffer per file: %d MiB\n", *bind_addr, *port, getProv(), *maxMemMB)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bind_addr, *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
