package job

type Job[T any] interface {
	Combine(jobs []Job[T]) MultiJob[T]
	Hash() string
	SetHash(hash string)
	Do() JobResult[T]
}

type MultiJob[T any] interface {
	Combine(jobs []Job[T]) MultiJob[T]
	Hash() string
	SetHash(hash string)
	Do() map[string]JobResult[T]
}

type MultiJobBase[T any] struct {
	hash string
	job  Job[T]
}

func (m *MultiJobBase[T]) Combine(jobs []Job[T]) MultiJob[T] {
	return m.job.Combine(jobs)
}

func (m *MultiJobBase[T]) Hash() string {
	return m.hash
}

func (m *MultiJobBase[T]) SetHash(hash string) {
	m.hash = hash
}

func (m *MultiJobBase[T]) Do() map[string]JobResult[T] {
	return map[string]JobResult[T]{m.Hash(): m.job.Do()}
}

var _ MultiJob[int] = (*MultiJobBase[int])(nil)

func Wrap[T any](job Job[T]) MultiJob[T] {
	return &MultiJobBase[T]{
		job:  job,
		hash: job.Hash(),
	}
}
