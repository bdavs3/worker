package worker

import (
	"context"
	"io"
	"os/exec"

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
	log   *Log
	killC chan string
}

// NewWorker returns a Worker containing a new log.
func NewWorker() *Worker {
	return &Worker{
		log:   NewLog(),
		killC: make(chan string),
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	go func() {
		for {
			if status, err := w.Status(jobID); err != nil || status != statusActive {
				return
			}
			id := <-w.killC
			if id == jobID {
				w.log.setStatus(jobID, statusKilled)
				cancel()
			}
		}
	}()

	done := make(chan struct{}) // Empty struct to not use memory.
	go func() {
		err := pipeToOutput(w.log, jobID, stdout)
		if err != nil {
			result <- Result{Err: &ServerError{"Failed to write output to the job log."}}
			return
		}
		done <- struct{}{}
	}()
	<-done

	err = cmd.Wait()
	if err != nil {
		if err.Error() != "signal: killed" {
			w.log.setStatus(jobID, statusError)
		}
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
			// The EOF error is ok since it will occur when the command has finished executing.
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
	status, err := w.log.getStatus(id)
	if err != nil {
		return "", err
	}
	if status != statusActive {
		return "", &NotActiveErr{"Can't kill inactive job."}
	}

	w.killC <- id

	return statusKilled, nil
}
