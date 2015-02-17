// Qute is the gosh darned cutest queue you'll ever meet.
package qute

// Job is the unit of work for a Qute queue. All it has to do is Do something.
// There are no return values proxied, so you should store your Job state on your
// Job implementation.
type Job interface {
	Do()
}

// Start() takes a Job channel and pool size and will run the queue in the background.
// To stop the queue, close the Job channel. Start returns a bool signal channel that
// will close when all jobs have stopped gracefully.
//
// For example:
//     c := make(chan Job)
//     done := Start(n, 5)
//     c <- someJob
//     close(c) // signal to stop processing jobs
//     <-done // this blocks until all jobs are done
func Start(c chan Job, n int) chan bool {
	done := make(chan bool)
	go func() {
		Run(c, n)
		close(done)
	}()
	return done
}

// Run() takes a Job channel and pool size and will block until all jobs are done.

// For example:
//     c := make(chan Job)
//     c <- someJob
//     close(c)
//     Run(c, 5)
//     // someJob is done
func Run(c chan Job, n int) {
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
