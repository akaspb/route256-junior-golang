package ttl_cash

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cash"
	"time"
)

type KeyType string
type Data string

type CashData struct {
	validUntil time.Time
	data       Data
}

func (c *CashData) Value() Data {
	return c.data
}

func (c *CashData) Validate(now time.Time) bool {
	return c.validUntil.Before(now)
}

type TTLCashFactory struct {
	validTime time.Duration
}

func NewTTLCashFactory(validTime time.Duration) *TTLCashFactory {
	return &TTLCashFactory{validTime: validTime}
}

func (f *TTLCashFactory) Create(
	data Data, now time.Time,
) (in_memory_cash.Cash[time.Time, Data], error) {
	return &CashData{
		validUntil: now.Add(f.validTime),
		data:       data,
	}, nil
}

func NewTTLCash(validTime time.Duration) *in_memory_cash.InMemoryCash[KeyType, time.Time, Data] {
	cashFactory := NewTTLCashFactory(validTime)
	return in_memory_cash.NewInMemoryCash[KeyType, time.Time, Data](cashFactory)
}
