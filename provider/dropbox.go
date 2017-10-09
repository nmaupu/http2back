package provider

import (
	"fmt"
	"github.com/tj/go-dropbox"
	"io"
)

var (
	_ Provider = Dropbox{}
)

type Dropbox struct {
	AccessToken, Dest string
}

func (d Dropbox) Copy(in io.Reader, name string) string {
	client := dropbox.New(dropbox.NewConfig(d.AccessToken))

	out, err := client.Files.Upload(&dropbox.UploadInput{
		Path:       fmt.Sprintf("%s/%s", d.Dest, name),
		Mode:       dropbox.WriteModeAdd,
		AutoRename: false,
		Mute:       false,
		Reader:     in,
	})
	if err != nil {
		panic(fmt.Sprintf("Unable to upload file - %s", err))
	}

	return out.PathDisplay
}

func (d Dropbox) String() string {
	return fmt.Sprintf("Dropbox (%s)", d.Dest)
}
