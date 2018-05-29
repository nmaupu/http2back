package notifier

import (
	"errors"
	"fmt"
	"github.com/zorkian/go-datadog-api"
)

var (
	_    Notifier = Datadog{}
	tags          = []string{
		"origin:http2back",
	}
)

type Datadog struct {
	ApiKey, AppKey, ExtraTag *string
}

func (d Datadog) Notify(event *Event) error {
	if d.ApiKey == nil || *d.ApiKey == "" {
		return errors.New("Cannot send event to Datadog, api key is not defined")
	}

	tags = append(tags, *d.ExtraTag)

	client := datadog.NewClient(*d.ApiKey, *d.AppKey)
	_, err := client.PostEvent(&datadog.Event{
		Title:     &event.Title,
		Text:      &event.Message,
		AlertType: datadog.String("info"),
		Tags:      tags,
	})

	return err
}

func (d Datadog) String() string {
	return fmt.Sprintf("Datadog notifier - ApiKey: %s, AppKey: %s", *d.ApiKey, *d.AppKey)
}
