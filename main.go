package main

import (
	"github.com/nmaupu/http2back/server"
)

func main() {
	addr := "127.0.0.1"
	port := 8080
	provider := server.Filesystem{"/tmp/test"}
	server.Start(&port, &addr, provider)
}
