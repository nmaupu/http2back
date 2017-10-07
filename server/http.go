package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Globals
var provider Provider = nil

// HTTP API
func handleRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprintf("Error : %s", r))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 - %s", r)))
		}
	}()

	err := r.ParseForm()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}

	if r.Method == "PUT" || r.Method == "POST" {
		r.ParseMultipartForm(32 << 20) // 32768
		in, handler, err := r.FormFile("file")
		defer in.Close()
		if err != nil {
			panic(fmt.Sprintf("Error: %s", err))
		}

		// Unique filename
		t := time.Now()
		name := fmt.Sprintf("%s-%s.%06d", handler.Filename, t.Format("20060102_150405"), rand.Intn(100000))
		ret := provider.Copy(in, name)

		defer log.Printf("Called : %s %s - file %s -> %s", r.Method, r.URL.Path, handler.Filename, ret)
		defer w.Write([]byte(fmt.Sprintf("Saved ok - %s", ret)))
	} else {
		panic(fmt.Sprintf("%s is an unsupported method for %s", r.Method, r.URL.Path))
	}
}

// Starting the HTTP server
func Start(port *int, bind_addr *string, p Provider) {
	provider = p

	log.Printf("Starting http server on %s:%d using provider %s", *bind_addr, *port, p)

	http.HandleFunc("/", handleRequest)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bind_addr, *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
