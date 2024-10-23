package kafka_logger

import (
	"encoding/json"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/kafka"
	"strconv"
)

type KafkaLogger struct {
	producer kafka.Producer
	topic    string
}

func NewKafkaLogger(producer kafka.Producer, topic string) *KafkaLogger {
	return &KafkaLogger{producer: producer, topic: topic}
}

func (s *KafkaLogger) Send(event event_logger.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.producer.Send(kafka.Message{
		Topic:     s.topic,
		Key:       strconv.FormatInt(event.ID, 10),
		Value:     bytes,
		Timestamp: event.Timestamp,
	})
}
