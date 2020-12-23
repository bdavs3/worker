package worker

import (
	"context"
	"io"
	"os/exec"
	"sync"

	"github.com/lithammer/shortuuid"
)

const (
	statusActive   = "active"
	statusComplete = "complete"
	statusError    = "error"
	statusKilled   = "killed"
)

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(ctx context.Context, result chan<- RunResult, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) (string, error)
}

// Worker provides the machinery for executing and controlling Linux processes.
// A zero value of this type is invalid - use NewWorker to create a new instance.
type Worker struct {
	log   *log
	killC map[string]chan bool
	mu    sync.RWMutex
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

func (dw *DummyWorker) Run(ctx context.Context, result chan<- RunResult, job Job) {
	result <- RunResult{}
}
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(id string) (string, error)   { return "", nil }

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// ErrOutputPipe occurs when the worker's StdoutPipe fails to be established.
type ErrOutputPipe struct{ msg string }

func (e *ErrOutputPipe) Error() string { return e.msg }

// ErrInvalidCmd occurs when a job cannot be started due to bad syntax.
type ErrInvalidCmd struct{ msg string }

func (e *ErrInvalidCmd) Error() string { return e.msg }

// ErrJobNotFound occurs when a job cannot be found in the worker log.
type ErrJobNotFound struct{ msg string }

func (e *ErrJobNotFound) Error() string { return e.msg }

// ErrJobNotActive occurs when termination is attempted on a job that
// is no longer active.
type ErrJobNotActive struct{ msg string }

func (e *ErrJobNotActive) Error() string { return e.msg }

// RunResult contains the job ID for a process that successfully began execution
// or an error for one that did not.
type RunResult struct {
	ID  string
	Err error
}

// Run initiates the execution of a Linux process.
func (w *Worker) Run(ctx context.Context, result chan<- RunResult, job Job) {
	// TODO (next): Block on jobs requiring stdin.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, job.Command, job.Args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result <- RunResult{Err: &ErrOutputPipe{"job failed to start: bad output pipe"}}
		return
	}

	err = cmd.Start()
	if err != nil {
		result <- RunResult{Err: &ErrInvalidCmd{"job failed to start: invalid syntax"}}
		return
	}

	jobID := shortuuid.New()
	result <- RunResult{ID: jobID}
	w.log.addEntry(jobID)

	quitListening := make(chan bool)
	go w.listenForKill(jobID, cancel, quitListening)

	done := make(chan bool)
	go w.writeOutput(jobID, stdout, done)
	<-done

	quitListening <- true

	err = cmd.Wait()
	if err != nil {
		if err.Error() != "signal: killed" { // Prefer to keep custom message.
			w.log.setStatus(jobID, statusError)
		}
		return
	}

	w.log.setStatus(jobID, statusComplete)
}

// listenForKill calls the provided CancelFunc if the worker's channel associated
// with the specified ID receives a value.
func (w *Worker) listenForKill(id string, cancel context.CancelFunc, quit chan bool) {
	w.mu.Lock()
	w.killC[id] = make(chan bool)
	killC := w.killC[id]
	w.mu.Unlock()

	select {
	case <-killC:
		cancel()
		w.log.setStatus(id, statusKilled)
	case <-quit:
	}

	w.mu.Lock()
	delete(w.killC, id)
	w.mu.Unlock()

	return
}

// writeOutput writes to the worker's log using the provided io.Reader.
func (w *Worker) writeOutput(id string, stdout io.Reader, done chan bool) {
	err := pipeToLog(id, w.log, stdout)
	if err != nil {
		w.log.setStatus(id, statusError)
		return
	}
	done <- true
}

func pipeToLog(id string, log *log, stdout io.Reader) error {
	bytes := make([]byte, 1024)
	for {
		n, err := stdout.Read(bytes)
		if err != nil {
			return err
		}
		if n > 0 {
			err = log.appendOutput(id, bytes[:n])
			if err != nil {
				return err
			}
		}
		if err != nil {
			// The EOF error is to be expected when the command has finished executing.
			if err == io.EOF {
				return nil
			}
			return err
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
	status, err := w.log.getStatus(id)
	if err != nil {
		return "", err
	}
	if status != statusActive {
		return "", &ErrJobNotActive{"can't kill inactive job"}
	}

	w.mu.Lock()
	w.killC[id] <- true
	w.mu.Unlock()

	return statusKilled, nil
}
