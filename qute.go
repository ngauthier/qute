package qute

// Job is the unit of work for a Qute queue. All it has to do is Do something.
// There are no return values proxied, so you should store your Job state on your
// Job implementation.
type Job interface {
	Do()
}

// Qute is the gosh darned cutest queue you'll ever meet.
// Qute takes a Job channel of jobs to do and the number of
// workers to run in a pool.
//
// Qute() will block until all jobs are done. So you should
// `go Qute()` if you want to run it in the background. To
// join it back into the foreground at your discretion, use
// the standard done channel wait:
//
//     c := make(chan Job)
//     n := 5
//     done := make(chan bool)
//     go func() {
//       Qute(c, n)
//       close(done)
//     }
//     // do stuff with c
//     c <- someJob
//     // Close c to signal no more jobs
//     close(c)
//     // wait for done to close signifying all workers have stopped
//     <-done
func Qute(c chan Job, n int) {
	queue := make(chan chan Job)
	done := make(chan bool)

	// Start N workers
	for i := 0; i < n; i++ {
		go work(queue, done)
	}

	// Queue loop
	for {
		select {
		// Get a Job (you bum!)
		case j, ok := <-c:
			// Nothing to do
			if !ok {
				// Signal workers to return
				close(done)
				// Retrieve all workers (all work must be done)
				for i := 0; i < n; i++ {
					<-queue
				}
				// Done
				return
			}

			// Get a worker
			w := <-queue
			// Give the worker a job
			w <- j
			// Async wait for job to come back
			go func() {
				// Get the job from the worker and proxy it up the channel
				c <- <-w
			}()
		}
	}
}

func work(queue chan chan Job, done chan bool) {
	// Our work channel where we get jobs
	work := make(chan Job)
	for {
		// Put myself on the queue as an available worker
		queue <- work
		select {
		// Get a job off my own queue (sent to me)
		case j := <-work:
			// Do the work
			j.Do()
			// Send the job back up my queue to the sender so they don't
			// have to remember which job they sent me
			work <- j
		// done signal, so return
		case <-done:
			return
		}
	}
}
