package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

// Handler implements API methods for invoking the various capabilities
// of the worker it contains.
type Handler struct {
	Worker worker.JobWorker
}

// NewHandler initializes a Handler with the provided worker.
func NewHandler(worker worker.JobWorker) *Handler {
	return &Handler{
		Worker: worker,
	}
}

// PostJob initiates the worker's execution of the job contained in the
// request and if successful, responds with the unique id assigned to that job.
func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read request."))
		return
	}

	var job worker.Job
	err = json.Unmarshal(reqBody, &job)
	if err != nil || len(job.Command) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request does not contain a valid job."))
		return
	}

	id := h.Worker.Run(job)

	fmt.Fprint(w, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to complete request."))
	}
}

// GetJobStatus queries the worker's log to respond with the status of
// the given job.
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, h.Worker.Status(id))
}

// GetJobOutput queries the worker's log to respond with the output of
// the given job.
func (h *Handler) GetJobOutput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, h.Worker.Out(id))
}

// KillJob initiates the termination of a given job and if successful,
// responds with the new status of the job.
func (h *Handler) KillJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, h.Worker.Kill(id))
}
