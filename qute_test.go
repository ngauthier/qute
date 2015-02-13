package qute

import "testing"

type NoopJob struct{}

func (j *NoopJob) Do() {}

func TestWork(t *testing.T) {
	n := 1

	j := &NoopJob{}

	queue := make(chan chan Job, n)
	done := make(chan bool)

	go Work(queue, done)

	c := <-queue
	c <- j

	j2 := <-c

	if j != j2 {
		t.Fatal("didn't get job back")
	}

	close(done)
}
