package in_memory_cash

type Cash[V comparable, D any] interface {
	Value() D
	Validate(validateArg V) bool
}

type CashFactory[V comparable, D any] interface {
	Create(data D, validateValue V) (Cash[V, D], error)
}
