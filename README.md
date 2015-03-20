# Qute

The cutest queue you'll ever meet.

To use Qute, write jobs that implement `Job`, then run `Start(c, n)` where `c` is a Job channel and `n` is your pool size. For example:

```go
c := make(chan Job)

done := Start(c, 10)

// Send job
c <- myJob

// You'll get the same job back when it's done
myJobCameBack := <-c

// Close to initiate graceful shutdown
close(c)
// Wait for queue to stop
<-done
```
## Usage

#### func  Run

```go
func Run(c chan Job, n int)
```
For example:

    c := make(chan Job)
    c <- someJob
    close(c)
    Run(c, 5)
    // someJob is done

#### func  Start

```go
func Start(c chan Job, n int) chan bool
```
Start() takes a Job channel and pool size and will run the queue in the
background. To stop the queue, close the Job channel. Start returns a bool
signal channel that will close when all jobs have stopped gracefully.

For example:

    c := make(chan Job)
    done := Start(c, 5)
    c <- someJob
    close(c) // signal to stop processing jobs
    <-done // this blocks until all jobs are done

#### type Job

```go
type Job interface {
	Do()
}
```

Job is the unit of work for a Qute queue. All it has to do is Do something.
There are no return values proxied, so you should store your Job state on your
Job implementation.
