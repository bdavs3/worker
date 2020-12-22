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
	Run(result chan<- Result, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) (string, error)
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(result chan<- Result, job Job) { result <- Result{} }
func (dw *DummyWorker) Status(id string) (string, error)  { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)     { return "", nil }
func (dw *DummyWorker) Kill(id string) (string, error)    { return "", nil }

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

// ServerError occurs when the library cannot establish the stdout pipe or when
// an error happens while writing to the log.
type ServerError struct{ msg string }

func (e *ServerError) Error() string { return e.msg }

// CmdSyntaxError occurs when a job cannot be started due to bad syntax.
type CmdSyntaxError struct{ msg string }

func (e *CmdSyntaxError) Error() string { return e.msg }

// Result contains a job ID if a process successfully begins execution. If not,
// it contains the error that occurred.
type Result struct {
	ID  string
	Err error
}

// Run initiates the execution of a Linux process.
func (w *Worker) Run(result chan<- Result, job Job) {
	// TODO (next): Consider storing cancel funcs separately from the log.
	ctx, cancel := context.WithCancel(context.Background())

	// TODO (next): Block on jobs requiring stdin.
	cmd := exec.CommandContext(ctx, job.Command, job.Args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result <- Result{Err: &ServerError{"Job output pipe could not be established."}}
		return
	}

	err = cmd.Start()
	if err != nil {
		result <- Result{Err: &CmdSyntaxError{"Job failed to start due to invalid syntax."}}
		return
	}

	jobID := shortuuid.New()
	result <- Result{ID: jobID}
	w.log.addEntry(jobID, cancel)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := pipeToOutput(w.log, jobID, stdout)
		if err != nil {
			result <- Result{Err: &ServerError{"Failed to write output to the job log."}}
			return
		}
		wg.Done()
	}()
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		w.log.setStatus(jobID, statusError)
		return
	}

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
			// The EOF error is ok since this will occur the command has finished executing.
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

// Kill terminates a given process.
func (w *Worker) Kill(id string) (string, error) {
	cancel, err := w.log.getCancelFunc(id)
	if err != nil {
		return "", err
	}

	cancel()

	return statusKilled, nil
}
