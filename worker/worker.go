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
	statusKilled   = "killed"
)

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(id chan<- string, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) (string, error)
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(id chan<- string, job Job)    { id <- "" }
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(id string) (string, error)   { return "", nil }

// Worker is a JobWorker containing a log for the status/output of jobs.
type Worker struct {
	log *Log
}

// NewWorker returns a Worker containing a new log.
func NewWorker() *Worker {
	return &Worker{
		log: NewLog(),
	}
}

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Run initiates the execution of a Linux process.
func (w *Worker) Run(id chan<- string, job Job) {
	jobID := shortuuid.New()
	id <- jobID

	// TODO (next): Consider storing cancel funcs separately from the log.
	ctx, cancel := context.WithCancel(context.Background())

	// TODO (next): Block on jobs requiring stdin.
	cmd := exec.CommandContext(ctx, job.Command, job.Args...)
	stdout, _ := cmd.StdoutPipe()
	w.log.addEntry(jobID, cancel)

	cmd.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pipeToOutput(w.log, jobID, stdout)
		wg.Done()
	}()
	wg.Wait()

	cmd.Wait()

	w.log.setStatus(jobID, statusComplete)
}

func pipeToOutput(log *Log, id string, r io.Reader) error {
	output := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(output)
		if n > 0 {
			err = log.appendOutput(id, output)
			if err != nil {
				return err
			}
		}
		if err != nil {
			// The EOF error is ok since it simply means the reader has closed.
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

// KillResult contains a message indicating whether a process was killed and
// an error that may have occured during termination.
type KillResult struct {
	Message string
	Err     error
}

// Kill terminates a given process.
func (w *Worker) Kill(id string) (string, error) {
	cancel, err := w.log.getCancelFunc(id)
	if err != nil {
		return "", err
	}

	cancel()

	return statusKilled, nil
}
