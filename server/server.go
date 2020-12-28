package main

import (
	"fmt"
	"net/http"

	"github.com/bdavs3/worker/server/api"
	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

// TODO (next): Change to 443 once serving with TLS.
const (
	port    = "8080"
	idMatch = "[a-zA-Z0-9]*"
)

func main() {
	worker := worker.NewWorker()
	handler := api.NewHandler(worker)

	router := mux.NewRouter()

	router.HandleFunc("/jobs/run", auth.Secure(handler.PostJob)).Methods(http.MethodPost)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/status", auth.Secure(handler.GetJobStatus)).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/out", auth.Secure(handler.GetJobOutput)).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/kill", auth.Secure(handler.KillJob)).Methods(http.MethodPut)

	fmt.Println("Listening...")
	// TODO (next): ListenAndServeTLS by using a pre-generated private key
	// and self-signed certificate located inside the repository.
	http.ListenAndServe(":"+port, router)
}
