package provider

import (
	"fmt"
	"io"
	"os"
)

var (
	_ Provider = Filesystem{}
)

type Filesystem struct {
	DestDir string
}

func (f Filesystem) String() string {
	return fmt.Sprintf("Filesystem (%s)", f.DestDir)
}

func (f Filesystem) Copy(in io.Reader, name string) string {
	filename := fmt.Sprintf("%s/%s", f.DestDir, name)
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
	defer out.Close()

	io.Copy(out, in)
	return filename
}
