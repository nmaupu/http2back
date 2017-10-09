package provider

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
)

var (
	_ Provider = Ftp{}
)

type Ftp struct {
	Addr, Username, Password, DestDir string
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

func (f Ftp) String() string {
	return fmt.Sprintf("Ftp (ftp://%s@%s/%s)", f.Username, f.Addr, f.DestDir)
}
