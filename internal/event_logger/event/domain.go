package event

import "time"

type EventType string

const (
	EventReceiveOrder EventType = "order-created"
	EventGiveOrders   EventType = "order-canceled"
	EventReturnOrder  EventType = "order-canceled"
)

type Event struct {
	ID        int64     `json:"id"`
	EventType EventType `json:"event"`
	Operation []byte    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}
