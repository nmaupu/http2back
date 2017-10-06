package server

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func push(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Calling %s on %s", r.Method, r.URL.Path)
	if r.Method == "PUT" || r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file_in, handler, err := r.FormFile("file")
		defer file_in.Close()
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		t := time.Now()
		// Unique filename
		filename := fmt.Sprintf("/tmp/%s-%s.%06d", handler.Filename, t.Format("20060102_150405"), rand.Intn(100000))
		file_out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
		defer file_out.Close()
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		io.Copy(file_out, file_in)
	}
}

func Start(port *int, bind_addr *string) {
	http.HandleFunc("/push", push)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bind_addr, *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
