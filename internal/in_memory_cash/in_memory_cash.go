package in_memory_cash

import "errors"

var (
	ErrorCashIsInvalid = errors.New("cash is invalid")
	ErrorCashNotFound  = errors.New("cash was not found")
)

type InMemoryCash[K, V comparable, D any] struct {
	cashFactory CashFactory[V, D]
	memory      map[K]Cash[V, D]
}

func NewInMemoryCash[K, V comparable, D any](cashFactory CashFactory[V, D]) *InMemoryCash[K, V, D] {
	return &InMemoryCash[K, V, D]{
		cashFactory: cashFactory,
		memory:      make(map[K]Cash[V, D]),
	}
}

func (c *InMemoryCash[K, V, D]) Get(key K, validateArg any) (res D, err error) {
	cash, ok := c.memory[key]
	if !ok {
		err = ErrorCashNotFound
		return
	}

	if !cash.Validate(validateArg) {
		delete(c.memory, key)
		err = ErrorCashIsInvalid
		return
	}

	res = cash.Value()
	return
}

func (c *InMemoryCash[K, V, D]) Set(key K, data D, validateValue V) (err error) {
	cash, err := c.cashFactory.Create(data, validateValue)
	if err != nil {
		return
	}

	c.memory[key] = cash
	return
}

func (c *InMemoryCash[K, V, D]) Delete(key K) (err error) {
	_, ok := c.memory[key]
	if !ok {
		err = ErrorCashNotFound
		return
	}

	delete(c.memory, key)

	return
}
