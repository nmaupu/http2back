package server

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Provider interface {
	Copy(in io.Reader, name string) string
	String() string
}

type Filesystem struct {
	DestDir string
}

func (f Filesystem) Copy(in io.Reader, name string) string {
	filename := fmt.Sprintf("%s/%s", f.DestDir, name)
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	defer out.Close()
	if err != nil {
		log.Println("Error: ", err)
		return ""
	}

	io.Copy(out, in)
	return filename
}

func (f Filesystem) String() string {
	return fmt.Sprintf("Filesystem (%s)", f.DestDir)
}
