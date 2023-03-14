package batchq

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

type Map[T any] interface {
	Get(key string) (JobResult[T], bool)
	Set(key string, value JobResult[T])
	Remove(key string)
	Count() int
	Start()
	Stop()
}

type ValuePair[T any] struct {
	Key   string
	Value T
}

type MapBase[T any] struct {
	super cmap.ConcurrentMap[string, JobResult[T]]
	stop  chan bool
	dur   time.Duration
}

func (m *MapBase[T]) Get(key string) (JobResult[T], bool) {
	return m.super.Get(key)
}

func (m *MapBase[T]) Set(key string, value JobResult[T]) {
	m.super.Set(key, value)
}

func (m *MapBase[T]) Remove(key string) {
	m.super.Remove(key)
}

func (m *MapBase[T]) Count() int {
	return m.super.Count()
}

func (m *MapBase[T]) Start() {
	for {
		select {
		case <-m.stop:
			return
		default:
			for _, key := range m.super.Keys() {
				if value, found := m.super.Get(key); found {
					if value == nil {
						m.super.Remove(key)
						continue
					}
					if value.FinishedAt().Add(m.dur).Before(time.Now()) {
						m.super.Remove(key)
					}
				}
			}
		}
	}
}

func (m *MapBase[T]) Stop() {
	m.stop <- true
}

func NewMapBase[T any]() *MapBase[T] {
	return &MapBase[T]{
		super: cmap.New[JobResult[T]](),
		stop:  make(chan bool),
	}
}

var _ Map[int] = (*MapBase[int])(nil)
