package ttl_cash

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cash"
	"time"
)

type CashData[D any] struct {
	validUntil time.Time
	data       D
}

func (c *CashData[D]) Value() D {
	return c.data
}

func (c *CashData[D]) Validate(now time.Time) bool {
	return c.validUntil.Before(now)
}

type TTLCashFactory[K comparable, D any] struct {
	validTime time.Duration
}

func NewTTLCashFactory[K comparable, D any](validTime time.Duration) *TTLCashFactory[K, D] {
	return &TTLCashFactory[K, D]{validTime: validTime}
}

func (f *TTLCashFactory[K, D]) Create(
	data D, now time.Time,
) (in_memory_cash.Cash[time.Time, D], error) {
	return &CashData[D]{
		validUntil: now.Add(f.validTime),
		data:       data,
	}, nil
}

func NewTTLCash[K comparable, D any](validTime time.Duration) *in_memory_cash.InMemoryCash[K, time.Time, D] {
	cashFactory := NewTTLCashFactory[K, D](validTime)
	return in_memory_cash.NewInMemoryCash[K, time.Time, D](cashFactory)
}
