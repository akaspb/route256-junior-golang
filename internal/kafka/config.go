package kafka

import (
	"github.com/IBM/sarama"
	"time"
)

type Config struct {
	Brokers []string
}

func PrepareConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	// Default options:
	WithProducerPartitioner(sarama.NewHashPartitioner).Apply(c)
	WithRequiredAcks(sarama.WaitForAll).Apply(c)
	WithMaxRetries(100).Apply(c)
	WithRetryBackoff(5 * time.Millisecond).Apply(c)
	WithMaxOpenRequests(1).Apply(c)

	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	c.Producer.Compression = sarama.CompressionGZIP

	c.Producer.Return.Successes = false
	c.Producer.Return.Errors = true

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
