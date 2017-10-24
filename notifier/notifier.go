package notifier

type Notifier interface {
	Notify(event *Event) error
}

// Simple Event type used to Notify
type Event struct {
	Title, Message string
}
