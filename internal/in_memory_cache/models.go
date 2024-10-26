package in_memory_cache

type Cache[V comparable, D any] interface {
	Value() D
	Validate(validateArg V) bool
}

type CacheFactory[V comparable, D any] interface {
	Create(data D, validateValue V) Cache[V, D]
}
