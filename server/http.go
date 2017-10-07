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
func push(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "PUT" || r.Method == "POST" {
		r.ParseMultipartForm(32 << 20) // 32768
		in, handler, err := r.FormFile("file")
		defer in.Close()
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		t := time.Now()
		// Unique filename
		name := fmt.Sprintf("%s-%s.%06d", handler.Filename, t.Format("20060102_150405"), rand.Intn(100000))
		ret := provider.Copy(in, name)
		log.Printf("Called : %s %s - file %s -> %s", r.Method, r.URL.Path, handler.Filename, ret)
	}
}

// Starting the HTTP server
func Start(port *int, bind_addr *string, p Provider) {
	provider = p

	log.Printf("Starting http server on %s:%d using provider %s", *bind_addr, *port, p)

	http.HandleFunc("/push", push)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bind_addr, *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
