package qute

import "testing"

type NoopJob struct{}

func (j *NoopJob) Do() {}

func TestQute(t *testing.T) {
	n := 8

	j := &NoopJob{}

	c := make(chan Job)

	done := make(chan bool)
	go func() {
		Qute(c, n)
		close(done)
	}()

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
