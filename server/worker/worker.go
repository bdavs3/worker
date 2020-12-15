package worker

import (
	"fmt"
	"net/http"
)

// Job represents a Linux process that the client passes to the server.
type Job struct {
	Command string
	Args    []string
}

// Run initiates the execution of a process passed by the client.
func Run(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Run")
}

// Status queries the log for the status of a given process.
func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status")
}

// Out queries the log for the output of a given process.
func Out(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Out")
}

// Kill terminates a given process.
func Kill(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kill")
}
