package batchq

import "time"

type BatchQ[T any] struct {
	jobChan   chan Job[T]
	n         int
	resultMap Map[JobResult[T]]
	stopChan  chan bool
	dur       time.Duration
}

func NewBatchQ[T any]() *BatchQ[T] {
	return &BatchQ[T]{}
}

func (q *BatchQ[T]) Start() {
	go q.StartBlock()
}

func (q *BatchQ[T]) Check(hash string) (found bool, result JobResult[T]) {
	if result, found := q.resultMap.Get(hash); found {
		return true, result
	}
	return false, result
}

func (q *BatchQ[T]) Add(job Job[T]) string {
	q.jobChan <- job
	return job.Hash()
}

func (q *BatchQ[T]) Stop() {
	q.stopChan <- true
}

func (q *BatchQ[T]) SetBatchSize(n int) {
	q.n = n
}

func (q *BatchQ[T]) process(jobs []Job[T]) {
	var multi MultiJob[T]
	if len(jobs) == 1 {
		multi = Wrap[T](jobs[0])
	} else {
		multi = jobs[0].Combine(jobs[1:])
	}
	multi.Do()
}

func (q *BatchQ[T]) StartBlock() {
	var jobs []Job[T]
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
			return
		case job := <-q.jobChan:
			if len(jobs) == 0 {
				firstTime = time.Now()
			}
			jobs = append(jobs, job)
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
