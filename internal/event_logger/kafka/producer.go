package kafka

import (
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/event"
)

func NewAsyncProducer(conf Config, opts ...Option) (sarama.AsyncProducer, error) {
	config := PrepareConfig(opts...)

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return asyncProducer, nil
}

type TopicSender struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewTopicSender(producer sarama.AsyncProducer, topic string) *TopicSender {
	return &TopicSender{producer: producer, topic: topic}
}

func (s *TopicSender) Send(event event.Event) error {
	msg := &sarama.ProducerMessage{
		Topic:     s.topic,
		Key:       sarama.StringEncoder(strconv.FormatInt(event.ID, 10)),
		Value:     sarama.ByteEncoder(event.Operation),
		Timestamp: event.Timestamp,
	}

	s.producer.Input() <- msg

	return nil
}
