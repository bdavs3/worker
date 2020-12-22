package main

import (
	"fmt"
	"net/http"

	"github.com/bdavs3/worker/server/api"
	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"

	"github.com/gorilla/mux"
)

const (
	port    = "443"
	crtFile = "../worker.crt"
	keyFile = "../worker.key"
)

func main() {
	worker := worker.NewWorker()
	handler := api.NewHandler(worker)

	router := mux.NewRouter()

	router.HandleFunc("/jobs/run", auth.Secure(handler.PostJob)).Methods(http.MethodPost)
	router.HandleFunc("/jobs/{id}/status", auth.Secure(handler.GetJobStatus)).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id}/out", auth.Secure(handler.GetJobOutput)).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id}/kill", auth.Secure(handler.KillJob)).Methods(http.MethodPut)

	fmt.Println("Listening...")
	http.ListenAndServeTLS(":"+port, crtFile, keyFile, router)
}
