package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO (next): Execute jobs passed to this library concurrently using
// goroutines. Keep track of job execution in a log stored in memory,
// ensuring that access to this log is synchronized but does not cause
// deadlock. Allow active processes to be terminated.

// Job represents a Linux process to be handled by the worker library.
type Job struct {
	Command string
	Args    []string
}

// Run will initiate the execution of a Linux process.
func Run(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read request."))
		return
	}

	var job Job
	err = json.Unmarshal(reqBody, &job)
	if err != nil || len(job.Command) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request does not contain a valid job."))
		return
	}

	// TODO (next): Rather than echoing the job back to the client, respond with
	// the unique ID assigned to the job.
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to complete request."))
	}
}

// Status will query the log for the status of a given process.
func Status(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Status of "+id)
}

// Out will query the log for the output of a given process.
func Out(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Output of "+id)
}

// Kill will terminate a given process.
func Kill(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, "Killing job "+id)
}
