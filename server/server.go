package main

import (
	"fmt"
	"net/http"

	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/server/worker"

	"github.com/gorilla/mux"
)

const (
	idRegex = "[1-9][0-9]*"
	port    = "8080"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/jobs/run", worker.Run).Methods(http.MethodPost)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/status", worker.Status).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/out", worker.Out).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/kill", worker.Kill).Methods(http.MethodPut)

	fmt.Println("Listening...")
	// TODO (next): ListenAndServeTLS by using a pre-generated private key
	// and self-signed certificate located inside the repository.
	http.ListenAndServe(":"+port, auth.Secure(router))
}
