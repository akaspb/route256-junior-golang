package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// Option is a configuration callback
type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

// WithProducerPartitioner - установить алгоритм выбор партиции
//   - sarama.NewManualPartitioner - ручной
//   - sarama.NewRandomPartitioner - случайная партиция
//   - sarama.NewRoundRobinPartitioner - по кругу
//   - sarama.NewHashPartitioner - по ключу
func WithProducerPartitioner(pfn sarama.PartitionerConstructor) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Partitioner = pfn
		return nil
	})
}

// WithRequiredAcks - установить параметр acks
//   - sarama.NoResponse - ничего не ждем
//   - sarama.WaitForLocal - ждем успешной записи ТОЛЬКО на лидер партиции
//   - sarama.WaitForLocal - ждем успешной записи на лидер партиции и всех in-sync реплик (настроено в кафка кластере)
func WithRequiredAcks(acks sarama.RequiredAcks) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.RequiredAcks = acks
		return nil
	})
}

// WithIdempotent - установить семантики exactly once (установить Idempotent = true)
/*
	У продюсера есть счетчик (count).
	Каждое успешно отправленное сообщение увеличивает счетчик (count++).
	Если продюсер не смог отправить сообщение, то счетчик не меняется и отправляется в таком виде в другом сообщение.
	Кафка это видит и начинает сравнивать (в том числе Key) сообщения с одинаковыми счетчиками.
	Далее не дает отправить дубль, если Idempotent = true.
*/
func WithIdempotent() Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Idempotent = true
		return nil
	})
}

// WithMaxRetries - установить число попыток отправить сообщение
func WithMaxRetries(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Max = n
		return nil
	})
}

// WithRetryBackoff - установить интервалы между попытками отправить сообщение
func WithRetryBackoff(d time.Duration) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Backoff = d
		return nil
	})
}

// WithMaxOpenRequests - установить пропускную способность
//
// При значении 1 гарантируем строгий порядок отправки сообщений/батчей
func WithMaxOpenRequests(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Net.MaxOpenRequests = n
		return nil
	})
}

// WithProducerFlushMessages - установить количество сообщений, которые должны быть собраны в очереди перед их отправкой
//
// Когда количество сообщений в очереди достигает этого значения, все собранные сообщения отправляются в брокеры.
// Установка этого значения позволяет уменьшить количество вызовов к Kafka, что может существенно повысить производительность.
// Вместо отправки каждого сообщения сразу, продюсер собирает несколько сообщений и отправляет их за одну операцию.
func WithProducerFlushMessages(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Flush.Messages = n
		return nil
	})
}

// WithProducerFlushFrequency - установить интервал времени, после которого все сообщения в очереди будут отправлены
func WithProducerFlushFrequency(d time.Duration) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Flush.Frequency = d
		return nil
	})
}
