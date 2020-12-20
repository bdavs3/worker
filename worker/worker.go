package worker

import "fmt"

// TODO (next): Execute jobs passed to this library concurrently using
// goroutines. Keep track of job execution in a log stored in memory,
// ensuring that access to this log is synchronized but does not cause
// deadlock. Allow active processes to be terminated.

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Run will initiate the execution of a Linux process.
func Run(job Job) string {
	return fmt.Sprint(5) // TODO (next): Return a UUID/GUID.
}

// Status will query the log for the status of a given process.
func Status(jobID string) string {
	return "Status of " + jobID
}

// Out will query the log for the output of a given process.
func Out(jobID string) string {
	return "Output of " + jobID
}

// Kill will terminate a given process.
func Kill(jobID string) string {
	return "Killing job " + jobID
}
