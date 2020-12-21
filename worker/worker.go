package worker

import (
	"context"
	"os/exec"

	"github.com/lithammer/shortuuid"
)

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(id chan<- string, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(result chan<- KillResult, id string)
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(id chan<- string, job Job)    { id <- "" }
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(result chan<- KillResult, id string) {
	result <- KillResult{"", nil}
}

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
	w.log.addEntry(jobID, cancel)

	// TODO (next): Block on jobs requiring stdin and capture their output.
	out, err := exec.CommandContext(ctx, job.Command, job.Args...).Output()
	w.log.setOutput(jobID, string(out))
	if err != nil {
		w.log.setStatus(jobID, "Error - "+err.Error())
		return
	}

	w.log.setStatus(jobID, "finished")
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
	output, err := w.log.getOutput(id)
	if err != nil {
		return "", err
	}

	return output, nil
}

// KillResult contains a message indicated whether a process was killed and
// any error that occured during termination.
type KillResult struct {
	Message string
	Err     error
}

// Kill terminates a given process.
func (w *Worker) Kill(result chan<- KillResult, id string) {
	cancel, err := w.log.getCancelFunc(id)
	if err != nil {
		result <- KillResult{"", err}
	}

	cancel()

	result <- KillResult{"killed", nil}
}
