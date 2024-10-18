package event_logger

import "time"

type EventType string

type Event struct {
	ID        int64     `json:"id"`
	EventType EventType `json:"event_type"`
	EventData []byte    `json:"event_data"`
	Timestamp time.Time `json:"timestamp"`
}
