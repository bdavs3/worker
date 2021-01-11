package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bdavs3/worker/server/api"
	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

const (
	crtFile = "../worker.crt"
	keyFile = "../worker.key"
	idMatch = "[a-zA-Z0-9]+"
)

// TODO (out of scope): In the interest of high availability, use a load balancer to
// distribute network traffic.

// TODO (out of scope): In the interest of high performance, start optimizing in the
// following ways...
// - Reduce lock contention using atomic operations or other data structures.
// - Pre-allocate memory where possible.
// - Avoid data copies, but don't overdo it.

func main() {
	port := os.Getenv("port")
	if len(port) == 0 {
		port = "443"
	}

	worker := worker.NewWorker()
	auth := auth.NewAuth()
	handler := api.NewHandler(worker, auth)

	router := mux.NewRouter()
	router.Use(auth.Authenticate)

	// A subrouter is used to avoid extraneous authorization checks.
	sub := router.Methods(http.MethodGet, http.MethodPut).Subrouter()
	sub.Use(auth.Authorize)

	router.HandleFunc("/jobs/run", handler.PostJob).Methods(http.MethodPost)
	sub.HandleFunc("/jobs/{id:"+idMatch+"}/status", handler.GetJobStatus).Methods(http.MethodGet)
	sub.HandleFunc("/jobs/{id:"+idMatch+"}/out", handler.GetJobOutput).Methods(http.MethodGet)
	sub.HandleFunc("/jobs/{id:"+idMatch+"}/kill", handler.KillJob).Methods(http.MethodPut)

	fmt.Println("Listening...")
	http.ListenAndServeTLS(":"+port, crtFile, keyFile, router)
}
