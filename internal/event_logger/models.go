package event_logger

type EventLogger interface {
	Send(event Event) error
}

type EventFactory interface {
	Create(eventType EventType, eventDetails string) (Event, error)
}
