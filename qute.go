package qute

type Job interface {
	Do()
}

func Work(queue chan chan Job, done chan bool) {
	work := make(chan Job)
	for {
		queue <- work
		select {
		case j := <-work:
			j.Do()
			work <- j
		case <-done:
			return
		}
	}
}
