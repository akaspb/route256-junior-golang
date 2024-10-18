package event_logger

type EventLogger interface {
	Send(event Event) error
}

type EventFactory interface {
	Create(eventType EventType, event []byte) (Event, error)
}
