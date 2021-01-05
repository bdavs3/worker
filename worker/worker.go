package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/lithammer/shortuuid"
)

const (
	statusActive   = "active"
	statusComplete = "complete"
	statusError    = "error"
	statusKilled   = "killed"
)

// A JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(job Job, reqComplete chan bool) string
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) (string, error)
}

// Worker provides the machinery for executing and controlling Linux processes.
// A zero value of this type is invalid - use NewWorker to create a new instance.
type Worker struct {
	log   *log
	killC map[string]chan bool
	mu    sync.Mutex
}

// NewWorker returns a Worker containing an empty log and channel map.
func NewWorker() *Worker {
	return &Worker{
		log:   newLog(),
		killC: make(map[string]chan bool),
	}
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(job Job, reqComplete chan bool) string {
	reqComplete <- true
	return ""
}
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(id string) (string, error)   { return "", nil }

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// ErrJobNotFound occurs when a job cannot be found in the worker log.
type ErrJobNotFound struct{ msg string }

func (e *ErrJobNotFound) Error() string { return e.msg }

// ErrJobNotActive occurs when termination is attempted on a job that
// is no longer active.
type ErrJobNotActive struct{ msg string }

func (e *ErrJobNotActive) Error() string { return e.msg }

// Run assigns a UUID to a Linux process and initiates its execution.
func (w *Worker) Run(job Job, reqComplete chan bool) string {
	id := shortuuid.New()
	w.log.addEntry(id)

	// Requests may be cancelled up until this point.
	reqComplete <- true

	go w.execJob(id, job)

	return id
}

func (w *Worker) execJob(id string, job Job) {
	cmdctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(cmdctx, job.Command, job.Args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err))
		return
	}

	// Direct stderr through stdout to interleave them as expected by command order.
	cmd.Stderr = cmd.Stdout

	err = cmd.Start()
	if err != nil {
		w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err))
		return
	}

	go w.listenForKill(cmdctx, cancel, id)
	w.writeOutput(id, stdout)

	err = cmd.Wait()
	if err != nil {
		// Prefer to keep 'kill' status if the process was terminated.
		status, _ := w.log.getStatus(id)
		if status != statusKilled {
			w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err))
		}
		return
	}

	w.log.setStatus(id, statusComplete)
}

// listenForKill calls the provided CancelFunc if the channel associated with
// the specified ID receives a value.
func (w *Worker) listenForKill(ctx context.Context, cancel context.CancelFunc, id string) {
	w.mu.Lock()
	w.killC[id] = make(chan bool)
	killC := w.killC[id]
	w.mu.Unlock()

	select {
	case <-killC:
		cancel()
		w.log.setStatus(id, statusKilled)
		killC <- true // Reply on the channel to signify that the job has been killed.
	case <-ctx.Done():
		// Job execution completed. Do nothing and proceed.
	}

	w.mu.Lock()
	delete(w.killC, id)
	w.mu.Unlock()
}

// writeOutput writes to the worker's log using the provided io.Reader.
func (w *Worker) writeOutput(id string, r io.ReadCloser) {
	bytes := make([]byte, 1024)
	for {
		n, err := r.Read(bytes)
		if err != nil {
			if err == io.EOF {
				return
			}
			w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err.Error()))
			return
		}
		if n > 0 {
			err := w.log.appendOutput(id, bytes[:n])
			if err != nil {
				w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err.Error()))
				return
			}
		}
	}
}

// Status queries the log for the status of a given process.
func (w *Worker) Status(id string) (string, error) {
	status, err := w.log.getStatus(id)
	if err != nil {
		return "", err
	}

	return status, nil
}

// Out queries the log for the output of a given process.
func (w *Worker) Out(id string) (string, error) {
	out, err := w.log.getOutput(id)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Kill terminates a given process using the channel associated with the provided ID.
func (w *Worker) Kill(id string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	killC, ok := w.killC[id]
	if !ok {
		return "", &ErrJobNotActive{"job not active"}
	}
	killC <- true

	// Await verification that the job has been killed, or return an error after
	// a timeout.
	select {
	case <-killC:
		return statusKilled, nil
	case <-time.After(50 * time.Millisecond):
		return "", errors.New("job not killed before timeout")
	}
}
