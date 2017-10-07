package main

import (
	"github.com/nmaupu/http2back/server"
)

func getProviderFilesystem() server.Provider {
	return server.Filesystem{"/tmp/test"}
}

func main() {
	addr := "127.0.0.1"
	port := 8080

	server.Start(&port, &addr, getProviderFilesystem)
}
