package server

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
)

type Provider interface {
	Copy(in io.Reader, name string) string
	String() string
}

/* FTP provider */
type Ftp struct {
	Addr, Username, Password, DestDir string
}

func (f Ftp) String() string {
	return fmt.Sprintf("Ftp (ftp://%s@%s/%s)", f.Username, f.Addr, f.DestDir)
}

func (f Ftp) Copy(in io.Reader, name string) string {
	var ftpConn *ftp.ServerConn
	var err error

	log.Printf(fmt.Sprintf("Connecting to ftp %s@%s", f.Username, f.Addr))
	if ftpConn, err = ftp.Connect(f.Addr); err != nil {
		panic(fmt.Sprintf("Unable to connect to Ftp %s - %s", f.Addr, err))
	}
	defer ftpConn.Quit()

	if err = ftpConn.Login(f.Username, f.Password); err != nil {
		panic(fmt.Sprintf("Unable to login with %s - %s", f.Username, err))
	}

	// Upload file
	remoteFilename := fmt.Sprintf("%s%s", f.DestDir, name)
	if err = ftpConn.Stor(remoteFilename, in); err != nil {
		panic(fmt.Sprintf("Unable to upload file - %s", err))
	}

	return remoteFilename
}

/* Filesystem provider */
type Filesystem struct {
	DestDir string
}

func (f Filesystem) String() string {
	return fmt.Sprintf("Filesystem (%s)", f.DestDir)
}

func (f Filesystem) Copy(in io.Reader, name string) string {
	filename := fmt.Sprintf("%s/%s", f.DestDir, name)
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	defer out.Close()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}

	io.Copy(out, in)
	return filename
}
