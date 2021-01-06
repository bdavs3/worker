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
	Run(job Job) string
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) error
}

// Worker provides the machinery for executing and controlling Linux processes.
// A zero value of this type is invalid - use NewWorker to create a new instance.
type Worker struct {
	log *log
	mu  sync.Mutex // Used to synchronize the termination of processes.
}

// NewWorker creates a new instance of the process worker.
func NewWorker() *Worker {
	return &Worker{
		log: newLog(),
	}
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(job Job) string               { return "" }
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(id string) error             { return nil }

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// ErrJobNotFound occurs when a process cannot be found in the worker log.
type ErrJobNotFound struct{ msg string }

func (e *ErrJobNotFound) Error() string { return e.msg }

// ErrJobNotActive occurs when termination is attempted on a process that
// is no longer active.
type ErrJobNotActive struct{ msg string }

func (e *ErrJobNotActive) Error() string { return e.msg }

// Run initiates the execution of a Linux process.
func (w *Worker) Run(job Job) string {
	id := shortuuid.New()

	w.log.addEntry(id)
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

// listenForKill is used to terminate a running process specified by Kill.
func (w *Worker) listenForKill(ctx context.Context, cancel context.CancelFunc, id string) {
	killC, _ := w.log.getKillC(id)

	select {
	case <-killC:
		cancel()
		w.log.setStatus(id, statusKilled)
		killC <- true // Reply on the channel to signify that the job has been killed.
	case <-ctx.Done():
		// Process execution completed. Do nothing and proceed.
	}
}

// writeOutput reads process output from r and writes it to the log entry
// associated with the given id.
func (w *Worker) writeOutput(id string, r io.ReadCloser) {
	bytes := make([]byte, 1024)
	for {
		n, err := r.Read(bytes)
		if err != nil {
			if err == io.EOF {
				return
			}
			w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err))
			return
		}
		if n > 0 {
			err := w.log.appendOutput(id, bytes[:n])
			if err != nil {
				w.log.setStatus(id, fmt.Sprintf("%s - %s", statusError, err))
				return
			}
		}
	}
}

// Status returns the status of the process represented by the given id.
func (w *Worker) Status(id string) (string, error) {
	status, err := w.log.getStatus(id)
	if err != nil {
		return "", err
	}

	return status, nil
}

// Out returns the output of the process represented by the given id.
func (w *Worker) Out(id string) (string, error) {
	out, err := w.log.getOutput(id)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Kill terminates the process represented by the given id.
func (w *Worker) Kill(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	status, err := w.log.getStatus(id)
	if err != nil {
		return &ErrJobNotFound{"job not found"}
	}
	if status != statusActive {
		return &ErrJobNotActive{"job not active"}
	}

	killC, _ := w.log.getKillC(id) // getStatus took care of 'not found' err.
	killC <- true

	select {
	case <-killC:
		return nil
	case <-time.After(50 * time.Millisecond):
		return errors.New("job not killed before timeout")
	}
}
