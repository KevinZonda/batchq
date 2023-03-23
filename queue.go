package batchq

import (
	"github.com/KevinZonda/batchq/job"
	"github.com/KevinZonda/batchq/map"
	"time"
)

type BatchQ[T any] struct {
	jobChan   chan job.Job[T]
	n         int
	resultMap _map.Map[T]
	stopChan  chan bool
	dur       time.Duration
	// canAppend Can append this job to origin jobs?
	// if false, will start a new job list
	canAppend func(origin []job.Job[T], newOne job.Job[T]) bool
}

func NewBatchQ[T any](numToBatch int, resultMap _map.Map[T], unitTime time.Duration, canAppend func(origin []job.Job[T], newOne job.Job[T]) bool) *BatchQ[T] {
	return &BatchQ[T]{
		jobChan:   make(chan job.Job[T]),
		n:         numToBatch,
		resultMap: resultMap,
		stopChan:  make(chan bool),
		dur:       unitTime,
		canAppend: canAppend,
	}
}

func NewBatchQEasy[T any](numToBatch int, unitTime time.Duration, canAppend func(origin []job.Job[T], newOne job.Job[T]) bool) *BatchQ[T] {
	return &BatchQ[T]{
		jobChan:   make(chan job.Job[T]),
		n:         numToBatch,
		resultMap: _map.NewMapBase[T](),
		stopChan:  make(chan bool),
		dur:       unitTime,
		canAppend: canAppend,
	}
}

func (q *BatchQ[T]) Start() {
	go q.StartBlock()
}

func (q *BatchQ[T]) Check(hash string) (found bool, result job.JobResult[T]) {
	if result, found := q.resultMap.Get(hash); found {
		return true, result
	}
	return false, result
}

func (q *BatchQ[T]) Add(job job.Job[T]) string {
	q.jobChan <- job
	return job.Hash()
}

func (q *BatchQ[T]) RemoveResult(hash string) {
	q.resultMap.Remove(hash)
}

func (q *BatchQ[T]) Stop() {
	q.stopChan <- true
}

func (q *BatchQ[T]) SetBatchSize(n int) {
	q.n = n
}

func (q *BatchQ[T]) process(jobs []job.Job[T]) {
	var multi job.MultiJob[T]
	if len(jobs) == 1 {
		multi = job.Wrap[T](jobs[0])
	} else {
		multi = jobs[0].Combine(jobs[1:])
	}
	rst := multi.Do()
	for hash, result := range rst {
		q.resultMap.Set(hash, result)
	}
}

func (q *BatchQ[T]) StartBlock() {
	q.resultMap.Start()
	var jobs []job.Job[T]
	firstTime := time.Now()
	f := func() {
		if len(jobs) > 0 {
			j := jobs
			jobs = nil
			go q.process(j)
		}
	}
	for {
		select {
		case <-q.stopChan:
			q.resultMap.Stop()
			return
		case jobC := <-q.jobChan:
			if len(jobs) == 0 {
				firstTime = time.Now()
			}
			if q.canAppend != nil && len(jobs) > 0 {
				if !q.canAppend(jobs, jobC) {
					f()
					jobs = []job.Job[T]{jobC}
					continue
				}
			}
			jobs = append(jobs, jobC)
			if len(jobs) == q.n {
				f()
			}
		case <-time.After(q.dur):
			f()
		default:
			if time.Since(firstTime) > q.dur {
				f()
			}
		}
	}
}
