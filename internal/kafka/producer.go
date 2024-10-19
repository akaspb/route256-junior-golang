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
