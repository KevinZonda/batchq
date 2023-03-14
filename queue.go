package batchq

import "time"

type BatchQ[T any] struct {
	jobChan chan Job[T]
	n       int
}

func NewBatchQ[T any]() *BatchQ[T] {
	return &BatchQ[T]{}
}

func (q *BatchQ[T]) Start() {

}

func (q *BatchQ[T]) Check(hash string) bool {
	return false
}

func (q *BatchQ[T]) Add(job Job[T]) string {
	q.jobChan <- job
	return job.Hash()
}

func (q *BatchQ[T]) Stop() {

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

func (q *BatchQ[T]) consume() {
	var jobs []Job[T]
	for {
		select {
		case job := <-q.jobChan:
			jobs = append(jobs, job)
			if len(jobs) == q.n {
				go q.process(jobs)
			}
		case <-time.After(1 * time.Second):
			if len(jobs) > 0 {
				go q.process(jobs)
			}

		}
	}
}
