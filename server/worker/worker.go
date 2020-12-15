package worker

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Status of "+id)
}

// Out queries the log for the output of a given process.
func Out(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Output of "+id)
}

// Kill terminates a given process.
func Kill(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Killing job "+id)
}
