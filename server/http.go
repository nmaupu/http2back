package server

import (
	"encoding/json"
	"fmt"
	"github.com/nmaupu/http2back/provider"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Globals
var getProv func() provider.Provider = nil

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
		r.ParseMultipartForm(32 << 20)
		in, handler, err := r.FormFile("file")
		if err != nil {
			panic(fmt.Sprintf("Error: %s", err))
		}
		defer in.Close()

		// Unique filename
		t := time.Now()
		name := fmt.Sprintf("%s-%s.%06d", handler.Filename, t.Format("20060102_150405"), rand.Intn(100000))
		ret := getProv().Copy(in, name)

		// Send result
		log.Printf("Called : %s %s - file %s -> %s", r.Method, r.URL.Path, handler.Filename, ret)
		jr := JsonResult{ret}
		j, _ := json.Marshal(jr)
		w.Write(j)

	} else {
		panic(fmt.Sprintf("%s is an unsupported method for %s", r.Method, r.URL.Path))
	}
}

// Starting the HTTP server
func Start(port *int, bind_addr *string, getProvider func() provider.Provider) {
	getProv = getProvider

	log.Printf("Starting http server on %s:%d using provider %s", *bind_addr, *port, getProv())

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bind_addr, *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
