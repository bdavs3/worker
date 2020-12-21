package worker

import (
	"os/exec"

	"github.com/lithammer/shortuuid"
)

// TODO (next): Execute jobs passed to this library concurrently using
// goroutines. Keep track of job execution in a log stored in memory,
// ensuring that access to this log is synchronized but does not cause
// deadlock. Allow active processes to be terminated.

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(id chan string, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) string
}

// DummyWorker is a JobWorker that can be passed to a handler in order
// to test the API independently.
type DummyWorker struct{}

// Run simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Run(id chan string, job Job) {}

// Status simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }

// Out simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Out(id string) (string, error) { return "", nil }

// Kill simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Kill(id string) string { return "" }

// Worker is a JobWorker that
type Worker struct {
	log *Log
}

// NewWorker initalizes a worker with the provided log.
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

// Run will initiate the execution of a Linux process.
func (w *Worker) Run(id chan string, job Job) {
	jobID := shortuuid.New()
	w.log.addEntry(jobID)
	id <- jobID

	out, err := exec.Command(job.Command, job.Args...).Output()
	if err != nil {
		w.log.setStatus(jobID, err.Error())
		return
	}

	w.log.setOutput(jobID, string(out))
	w.log.setStatus(jobID, "finished")
}

// Status will query the log for the status of a given process.
func (w *Worker) Status(id string) (string, error) {
	status, err := w.log.getStatus(id)
	if err != nil {
		return "", err
	}

	return status, nil
}

// Out will query the log for the output of a given process.
func (w *Worker) Out(id string) (string, error) {
	output, err := w.log.getOutput(id)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Kill will terminate a given process.
func (w *Worker) Kill(id string) string {
	return "Killing job " + id
}
