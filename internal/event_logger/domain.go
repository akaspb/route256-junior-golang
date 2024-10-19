package event_logger

import "time"

type EventType string

type Event struct {
	ID        int64     `json:"id"`
	Type      EventType `json:"type"`
	Details   any       `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}
