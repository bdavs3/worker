package worker

import (
	"os/exec"

	"github.com/lithammer/shortuuid"
)

// JobWorker implements methods to run/terminate Linux processes and
// query their output/status.
type JobWorker interface {
	Run(id chan string, job Job)
	Status(id string) (string, error)
	Out(id string) (string, error)
	Kill(id string) string
}

// DummyWorker implements the JobWorker interface so that the API can be tested
// independently.
type DummyWorker struct{}

func (dw *DummyWorker) Run(id chan string, job Job)      {}
func (dw *DummyWorker) Status(id string) (string, error) { return "", nil }
func (dw *DummyWorker) Out(id string) (string, error)    { return "", nil }
func (dw *DummyWorker) Kill(id string) string            { return "" }

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

// Kill will terminate a given process.
func (w *Worker) Kill(id string) string {
	return "Killing job " + id
}
