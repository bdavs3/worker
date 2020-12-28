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

// NewHandler returns a Handler with the provided JobWorker.
func NewHandler(worker worker.JobWorker) *Handler {
	return &Handler{
		Worker: worker,
	}
}

// PostJob initiates the worker's execution of the job contained in the
// request and if successful, responds with the id assigned to that job.
func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to read request"))
		return
	}

	var job worker.Job
	err = json.Unmarshal(reqBody, &job)
	if err != nil || len(job.Command) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request does not contain a valid job"))
		return
	}

	reqComplete := make(chan bool, 1)
	id := h.Worker.Run(job, reqComplete)

	select {
	case <-reqComplete:
		fmt.Fprint(w, id)
	case <-r.Context().Done():
		// Request cancelled.
	}
}

// GetJobStatus responds with the status of the given job.
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status, err := h.Worker.Status(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprint(w, status)
}

// GetJobOutput responds with the output of the given job.
func (h *Handler) GetJobOutput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	output, err := h.Worker.Out(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprint(w, output)
}

// KillJob terminates a given job and responds with its new status.
func (h *Handler) KillJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result, err := h.Worker.Kill(id)
	if err == nil {
		fmt.Fprint(w, result)
		return
	}
	switch err.(type) {
	case *worker.ErrJobNotActive:
		w.WriteHeader(http.StatusConflict)
	case *worker.ErrJobNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}
