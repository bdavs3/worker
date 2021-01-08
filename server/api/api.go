package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

// A Response contains information relevant to a particular process in the
// worker library.
type Response struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
	Output string `json:"output,omitempty"`
}

// Handler is an HTTP handler that manages processes on behalf of clients.
type Handler struct {
	Worker worker.JobWorker
	Auth   auth.SecurityLayer
}

// NewHandler initalizes a Handler with the given JobWorker and UserAuthLayer.
func NewHandler(worker worker.JobWorker, auth auth.SecurityLayer) *Handler {
	return &Handler{
		Worker: worker,
		Auth:   auth,
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

	username, _, _ := r.BasicAuth()
	h.Auth.SetOwner(username, id)

	response := &Response{ID: id}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling json", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

// GetJobStatus responds with the status of the process represented by the given id.
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status, err := h.Worker.Status(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := &Response{ID: id, Status: status}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling json", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

// GetJobOutput responds with the output of the process represented by the given id.
func (h *Handler) GetJobOutput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	output, err := h.Worker.Out(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := &Response{ID: id, Output: output}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling json", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

// KillJob terminates the job represented by the given id.
func (h *Handler) KillJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.Worker.Kill(id)
	if err != nil {
		switch err.(type) {
		case *worker.ErrJobNotActive:
			http.Error(w, err.Error(), http.StatusConflict)
		case *worker.ErrJobNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := &Response{ID: id, Status: "job successfully killed"}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling json", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
