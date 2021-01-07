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
	idMatch = "[a-zA-Z0-9]+"
)

func main() {
	worker := worker.NewWorker()
	auth := auth.NewAuth()
	handler := api.NewHandler(worker, auth)

	router := mux.NewRouter()
	router.Use(auth.Secure)

	router.HandleFunc("/jobs/run", handler.PostJob).Methods(http.MethodPost)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/status", handler.GetJobStatus).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/out", handler.GetJobOutput).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{id:"+idMatch+"}/kill", handler.KillJob).Methods(http.MethodPut)

	fmt.Println("Listening...")
	http.ListenAndServeTLS(":"+port, crtFile, keyFile, router)
}
