# Qute

The cutest queue you'll ever meet.

To use Qute, write jobs that implement `Job`, then run `Qute(c, n)` where `c` is a Job channel and `n` is your pool size. For example:

```golang
c := make(chan Job)
n := 10

go Qute(c,n)

// Send job
c <- myJob

// You'll get the same job back when it's done
myJobCameBack := <-c

// Close to gracefully shutdown
close(c)
```

#### func  Qute

```go
func Qute(c chan Job, n int)
```
Qute is the gosh darned cutest queue you'll ever meet. Qute takes a Job channel
of jobs to do and the number of workers to run in a pool.

Qute() will block until all jobs are done. So you should `go Qute()` if you want
to run it in the background. To join it back into the foreground at your
discretion, use the standard done channel wait:

    c := make(chan Job)
    n := 5
    done := make(chan bool)
    go func() {
      Qute(c, n)
      close(done)
    }
    // do stuff with c
    c <- someJob
    // Close c to signal no more jobs
    close(c)
    // wait for done to close signifying all workers have stopped
    <-done

#### type Job

```go
type Job interface {
	Do()
}
```

Job is the unit of work for a Qute queue. All it has to do is Do something.
There are no return values proxied, so you should store your Job state on your
Job implementation.
