package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
)

func NewAsyncProducer(conf Config, opts ...Option) (sarama.AsyncProducer, error) {
	config := PrepareConfig(opts...)

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewAsyncProducer failed: %w", err)
	}

	return asyncProducer, nil
}

type ProducerWrapper struct {
	producer sarama.AsyncProducer
}

func NewProducerWrapper(producer sarama.AsyncProducer) *ProducerWrapper {
	return &ProducerWrapper{
		producer: producer,
	}
}

func (p *ProducerWrapper) Send(message Message) error {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic:     message.Topic,
		Key:       sarama.StringEncoder(message.Key),
		Value:     sarama.ByteEncoder(message.Value),
		Timestamp: message.Timestamp,
	}

	return nil
}
