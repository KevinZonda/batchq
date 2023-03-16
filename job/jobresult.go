package job

import "time"

type JobResult[T any] interface {
	Hash() string
	Result() (T, error)
	Value() T
	Error() error
	Ok() bool
	FinishedAt() time.Time
}

type StackJobResult[T any] struct {
	hash   string
	result T
	err    error
	finish time.Time
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
func (s StackJobResult[T]) FinishedAt() time.Time {
	return s.finish
}

var _ JobResult[int] = (*StackJobResult[int])(nil)

func NewStackJobOk[T any](hash string, result T) StackJobResult[T] {
	return StackJobResult[T]{
		hash:   hash,
		result: result,
		err:    nil,
		finish: time.Now(),
	}
}

func NewStackJobError[T any](hash string, err error) StackJobResult[T] {
	return StackJobResult[T]{
		hash:   hash,
		err:    err,
		finish: time.Now(),
	}
}

type HeapJobResult[T any] struct {
	hash   string
	result T
	err    error
	finish time.Time
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

func (s *HeapJobResult[T]) FinishedAt() time.Time {
	return s.finish
}

var _ JobResult[int] = (*HeapJobResult[int])(nil)

func NewHeapJobOk[T any](hash string, result T) *HeapJobResult[T] {
	return &HeapJobResult[T]{
		hash:   hash,
		result: result,
		err:    nil,
		finish: time.Now(),
	}
}

func NewHeapJobError[T any](hash string, err error) *HeapJobResult[T] {
	return &HeapJobResult[T]{
		hash:   hash,
		err:    err,
		finish: time.Now(),
	}
}
