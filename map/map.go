package _map

import (
	"github.com/KevinZonda/batchq/job"
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

type Map[T any] interface {
	Get(key string) (job.JobResult[T], bool)
	Set(key string, value job.JobResult[T])
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
	super cmap.ConcurrentMap[string, job.JobResult[T]]
	stop  chan bool
	dur   time.Duration
	gc    bool
}

func (m *MapBase[T]) Get(key string) (job.JobResult[T], bool) {
	return m.super.Get(key)
}

func (m *MapBase[T]) Set(key string, value job.JobResult[T]) {
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
			if !m.gc {
				time.Sleep(time.Millisecond * 100)
				continue
			}
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
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (m *MapBase[T]) Stop() {
	m.stop <- true
}

func NewMapBase[T any](gc bool) *MapBase[T] {
	return &MapBase[T]{
		super: cmap.New[job.JobResult[T]](),
		stop:  make(chan bool),
		gc:    gc,
		dur:   time.Minute * 10,
	}
}

var _ Map[int] = (*MapBase[int])(nil)
