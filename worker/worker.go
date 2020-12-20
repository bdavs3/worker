package worker

import (
	"github.com/lithammer/shortuuid"
)

// TODO (next): Execute jobs passed to this library concurrently using
// goroutines. Keep track of job execution in a log stored in memory,
// ensuring that access to this log is synchronized but does not cause
// deadlock. Allow active processes to be terminated.

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(job Job) string
	Status(id string) string
	Out(id string) string
	Kill(id string) string
}

// DummyWorker is a JobWorker that can be passed to a handler in order
// to test the API independently.
type DummyWorker struct{}

// Run simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Run(job Job) string { return "" }

// Status simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Status(id string) string { return "" }

// Out simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Out(id string) string { return "" }

// Kill simply ensures that DummyWorker implements the JobWorker interface.
func (dw *DummyWorker) Kill(id string) string { return "" }

// Worker is a JobWorker that
type Worker struct {
	log *Log
}

// NewWorker initalizes a worker with the provided log.
func NewWorker(log *Log) *Worker {
	return &Worker{
		log: log,
	}
}

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Run will initiate the execution of a Linux process.
func (w *Worker) Run(job Job) string {
	id := shortuuid.New()
	w.log.addEntry(id)

	return id
}

// Status will query the log for the status of a given process.
func (w *Worker) Status(id string) string {
	return w.log.getStatus(id)
}

// Out will query the log for the output of a given process.
func (w *Worker) Out(id string) string {
	return w.log.getOutput(id)
}

// Kill will terminate a given process.
func (w *Worker) Kill(id string) string {
	return "Killing job " + id
}
