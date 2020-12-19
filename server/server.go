package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

const (
	// TODO (next): Change this regex to match chosen UUID/GUID format.
	idRegex = "[1-9][0-9]*"
	port    = "8080"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/jobs/run", postJob).Methods(http.MethodPost)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/status", getJobStatus).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/out", getJobOutput).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/kill", killJob).Methods(http.MethodPut)

	fmt.Println("Listening...")
	// TODO (next): ListenAndServeTLS by using a pre-generated private key
	// and self-signed certificate located inside the repository.
	http.ListenAndServe(":"+port, auth.Secure(router))
}

func postJob(w http.ResponseWriter, r *http.Request) {
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

	// TODO (next): Rather than echoing the job back to the client, respond with
	// the unique ID assigned to the job.
	err = json.NewEncoder(w).Encode(worker.Run(job))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to complete request."))
	}
}

func getJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, worker.Status(id))
}

func getJobOutput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, worker.Out(id))
}

func killJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprint(w, worker.Kill(id))
}
