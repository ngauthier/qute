package qute

import "testing"

type NoopJob struct{}

func (j *NoopJob) Do() {}

func TestStart(t *testing.T) {
	j := &NoopJob{}
	c := make(chan Job)
	done := Start(c, 8)

	for i := 0; i < 5; i++ {
		c <- j
		j2 := <-c
		if j != j2 {
			t.Fatal("Didn't get job back")
		}
	}

	close(c)
	<-done
}

func TestRun(t *testing.T) {
	j := &NoopJob{}
	c := make(chan Job)

	go func() {
		c <- j
	}()

	var j2 Job
	go func() {
		j2 = <-c
		close(c)
	}()

	Run(c, 5)

	if j != j2 {
		t.Fatal("Didn't get job back")
	}
}
