package batchq

import (
	"errors"
	"github.com/KevinZonda/batchq/job"
	"time"
)

func (q *BatchQ[T]) WaitJobResult(hash string, interval time.Duration) job.JobResult[T] {
	found := false
	var result job.JobResult[T]

startCheck:
	found, result = q.Check(hash)
	if found {
		return result
	}
	time.Sleep(interval)
	goto startCheck
}

func (q *BatchQ[T]) WaitResult(hash string, interval time.Duration) (result T, err error) {
	jr := q.WaitJobResult(hash, interval)
	if jr.Ok() {
		return jr.Value(), nil
	}
	return result, jr.Error()
}

func (q *BatchQ[T]) WaitJobResultChain(hash string, interval time.Duration) <-chan job.JobResult[T] {
	ch := make(chan job.JobResult[T])
	go func() {
		ch <- q.WaitJobResult(hash, interval)
		close(ch)
	}()
	return ch
}

var ErrTimeout = errors.New("timeout")

func (q *BatchQ[T]) WaitJobResultWithTimeout(hash string, interval, timout time.Duration) (result job.JobResult[T], err error) {
	select {
	case result = <-q.WaitJobResultChain(hash, interval):
		return
	case <-time.After(timout):
		err = ErrTimeout
	}
	return
}

func (q *BatchQ[T]) WaitResultWithTimeout(hash string, interval, timout time.Duration) (result T, err error) {
	jr, err := q.WaitJobResultWithTimeout(hash, interval, timout)
	if err != nil {
		return
	}
	if jr.Ok() {
		return jr.Value(), nil
	}
	return result, jr.Error()
}
