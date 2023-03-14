package batchq

type JobResult[T any] interface {
	Hash() string
	Result() (T, error)
	Value() T
	Error() error
	Ok() bool
}

type StackJobResult[T any] struct {
	hash   string
	result T
	err    error
}

func (s StackJobResult[T]) Hash() string {
	return s.hash
}

func (s StackJobResult[T]) Result() (T, error) {
	return s.result, s.err
}

func (s StackJobResult[T]) Value() T {
	return s.result
}

func (s StackJobResult[T]) Error() error {
	return s.err
}

func (s StackJobResult[T]) Ok() bool {
	return s.err == nil
}

var _ JobResult[int] = (*StackJobResult[int])(nil)

func NewStackJobOk[T any](hash string, result T) StackJobResult[T] {
	return StackJobResult[T]{
		hash:   hash,
		result: result,
		err:    nil,
	}
}

func NewStackJobError[T any](hash string, err error) StackJobResult[T] {
	return StackJobResult[T]{
		hash: hash,
		err:  err,
	}
}

type HeapJobResult[T any] struct {
	hash   string
	result T
	err    error
}

func (s *HeapJobResult[T]) Hash() string {
	return s.hash
}

func (s *HeapJobResult[T]) Result() (T, error) {
	return s.result, s.err
}

func (s *HeapJobResult[T]) Value() T {
	return s.result
}

func (s *HeapJobResult[T]) Error() error {
	return s.err
}

func (s *HeapJobResult[T]) Ok() bool {
	return s.err == nil
}

var _ JobResult[int] = (*StackJobResult[int])(nil)

func NewHeapJobOk[T any](hash string, result T) *HeapJobResult[T] {
	return &HeapJobResult[T]{
		hash:   hash,
		result: result,
		err:    nil,
	}
}

func NewHeapJobError[T any](hash string, err error) *HeapJobResult[T] {
	return &HeapJobResult[T]{
		hash: hash,
		err:  err,
	}
}
