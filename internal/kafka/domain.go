package kafka

import "time"

type Message struct {
	Topic     string
	Key       string
	Value     []byte
	Timestamp time.Time
}
