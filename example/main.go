package main

import (
	"fmt"
	"github.com/KevinZonda/batchq"
	"github.com/KevinZonda/batchq/job"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var q *batchq.BatchQ[string]

func main() {
	q = batchq.NewBatchQEasy[string](8*time.Second, batchq.NewCountConstraint[string](3))

	q.Start()
	http.HandleFunc("/", post)

	err := http.ListenAndServe("127.0.0.1:3306", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		return
	}
	fmt.Println("post:", string(body))
	j := &J{
		hash:    randString(16),
		content: string(body),
	}
	jhash := q.Add(j)
	fmt.Println("Job Hash", jhash)
	go func(hash string) {
		time.Sleep(10 * time.Second)
		fmt.Println("CHECK", hash)
		fmt.Println(q.Check(hash))
		fmt.Println("QLEN", q.QueueLength())
	}(jhash)
	io.WriteString(w, jhash)
}

type J struct {
	hash    string
	content string
}

func randString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

func (j *J) Combine(jobs []job.Job[string]) job.MultiJob[string] {
	fmt.Println("combine job:", j.content)
	js := &Js{
		hash:    randString(16),
		jh:      []string{j.hash},
		content: []string{j.content},
	}
	fmt.Println("combine job:", j.hash, "->", js.hash, j.content)
	for _, v := range jobs {
		x := v.(*J)
		fmt.Println("combine job:", x.hash, "->", js.hash, x.content)
		js.content = append(js.content, x.content)
		js.jh = append(js.jh, x.hash)
	}
	return js
}

func (j *J) Hash() string {
	return j.hash
}

func (j *J) SetHash(hash string) {
	j.hash = hash
}

func (j *J) Do() job.JobResult[string] {
	fmt.Println("do job:", j.content)
	return job.NewHeapJobOk[string](j.hash, j.content)
}

func (j *J) CanCombine(job job.Job[string]) bool {
	return true
}

var _ job.Job[string] = (*J)(nil)

type Js struct {
	hash    string
	content []string
	jh      []string
}

func (j *Js) Combine(jobs []job.Job[string]) job.MultiJob[string] {
	for _, v := range jobs {
		x := v.(*J)
		fmt.Println("combine m job:", x.hash, "->", j.hash, x.content)
		j.content = append(j.content, x.content)
		j.jh = append(j.jh, x.hash)
	}
	return j
}

func (j *Js) Hash() string {
	return j.hash
}

func (j *Js) SetHash(hash string) {
	j.hash = hash
}

func (j *Js) Do() map[string]job.JobResult[string] {
	fmt.Println("do m job:", j.content)
	m := make(map[string]job.JobResult[string])
	for i, v := range j.content {
		m[j.jh[i]] = job.NewHeapJobOk[string](j.hash, v)
	}
	return m
}

func (j *Js) CanCombine(jx job.Job[string]) bool {
	_, ok := jx.(*J)
	return ok
}

var _ job.MultiJob[string] = (*Js)(nil)
