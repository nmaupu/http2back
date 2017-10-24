package notifier

import (
	"fmt"
	"os"
	"time"
)

var (
	_ Notifier = File{}
)

type File struct {
	Dest *string
}

func (f File) Notify(event *Event) error {
	out, err := os.OpenFile(*f.Dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	message := fmt.Sprintf("%s %s %s\n", time.Now().Format("20060102_150405"), event.Title, event.Message)
	_, err = out.Write([]byte(message))
	return err
}

func (f File) String() string {
	return fmt.Sprintf("File notifier - Dest: %s", *f.Dest)
}
