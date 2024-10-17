package main

import (
	"time"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/kafka"
)

func initProducer(config kafka.Config) (sarama.AsyncProducer, error) {
	return kafka.NewAsyncProducer(config,
		// kafka.WithIdempotent(),
		kafka.WithRequiredAcks(sarama.WaitForAll),
		// kafka.WithMaxOpenRequests(1),
		kafka.WithMaxRetries(5),
		kafka.WithRetryBackoff(10*time.Millisecond),
		// kafka.WithProducerPartitioner(sarama.NewManualPartitioner),
		// kafka.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		// kafka.WithProducerPartitioner(sarama.NewRandomPartitioner),
		kafka.WithProducerPartitioner(sarama.NewHashPartitioner),
		kafka.WithProducerFlushMessages(3),
		kafka.WithProducerFlushFrequency(1*time.Second),
	)
}
