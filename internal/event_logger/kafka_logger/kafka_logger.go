package kafka_logger

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"strconv"
)

type KafkaLogger struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewKafkaLogger(producer sarama.AsyncProducer, topic string) *KafkaLogger {
	return &KafkaLogger{producer: producer, topic: topic}
}

func (s *KafkaLogger) Send(event event_logger.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic:     s.topic,
		Key:       sarama.StringEncoder(strconv.FormatInt(event.ID, 10)),
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: event.Timestamp,
	}

	s.producer.Input() <- message

	return nil
}
