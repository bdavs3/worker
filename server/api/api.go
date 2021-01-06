package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

// Handler is an HTTP handler that manages processes on behalf of clients.
type Handler struct {
	Worker worker.JobWorker
}

// NewHandler creates a new Handler instance with the given JobWorker.
func NewHandler(worker worker.JobWorker) *Handler {
	return &Handler{
		Worker: worker,
	}
}

// PostJob initiates the worker's execution of the process contained in the
// request and if successful, responds with the id assigned to that process.
func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request", http.StatusBadRequest)
		return
	}

	var job worker.Job
	err = json.Unmarshal(reqBody, &job)
	if err != nil || len(job.Command) == 0 {
		http.Error(w, "request does not contain a valid job", http.StatusBadRequest)
		return
	}

	id := h.Worker.Run(job)
	fmt.Fprint(w, id)
}

// GetJobStatus responds with the status of the process represented by the given id.
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status, err := h.Worker.Status(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(w, status)
}

// GetJobOutput responds with the output of the process represented by the given id.
func (h *Handler) GetJobOutput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	output, err := h.Worker.Out(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(w, output)
}

// KillJob terminates the job represented by the given id.
func (h *Handler) KillJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result, err := h.Worker.Kill(id)
	if err == nil {
		fmt.Fprint(w, result)
		return
	}
	switch err.(type) {
	case *worker.ErrJobNotActive:
		http.Error(w, err.Error(), http.StatusConflict)
	case *worker.ErrJobNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
