package batchq

import (
	"time"

	"github.com/KevinZonda/batchq/job"
	_map "github.com/KevinZonda/batchq/map"
)

type BatchQ[T any] struct {
	jobChan   chan job.Job[T]
	resultMap _map.Map[T]
	stopChan  chan bool
	dur       time.Duration
	// canAppend Can append this job to origin jobs?
	// if false, will start a new job list
	canAppend func(origin []job.Job[T], newOne job.Job[T]) bool
}

func NewBatchQ[T any](resultMap _map.Map[T], unitTime time.Duration, canAppend func(origin []job.Job[T], newOne job.Job[T]) bool) *BatchQ[T] {
	if canAppend == nil {
		canAppend = NewCountConstraint[T](10)
	}
	return &BatchQ[T]{
		jobChan:   make(chan job.Job[T], 100),
		resultMap: resultMap,
		stopChan:  make(chan bool, 1),
		dur:       unitTime,
		canAppend: canAppend,
	}
}

func NewBatchQEasy[T any](unitTime time.Duration, canAppend func(origin []job.Job[T], newOne job.Job[T]) bool) *BatchQ[T] {
	if canAppend == nil {
		canAppend = NewCountConstraint[T](10)
	}
	return &BatchQ[T]{
		jobChan:   make(chan job.Job[T], 100),
		resultMap: _map.NewMapBase[T](true),
		stopChan:  make(chan bool, 1),
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

func (q *BatchQ[T]) QueueLength() int {
	return len(q.jobChan)
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
	//go q.resultMap.Start()
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
				jobs = append(jobs, jobC)
				continue
			}
			// jobs > 0
			if !q.canAppend(jobs, jobC) {
				f()
				jobs = []job.Job[T]{jobC}
				continue
			}
			jobs = append(jobs, jobC)
		default:
			if time.Since(firstTime) > q.dur {
				f()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func NewCountConstraint[T any](maxCount int) func(origin []job.Job[T], newOne job.Job[T]) bool {
	return func(origin []job.Job[T], newOne job.Job[T]) bool {
		return len(origin) < maxCount
	}
}
