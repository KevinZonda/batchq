package batchq

type Map[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
	Count() int
	Start()
	Stop()
}

type ValuePair[T any] struct {
	Key   string
	Value T
}
