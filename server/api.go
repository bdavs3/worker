package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/server/worker"
)

const (
	idRegex = "[1-9][0-9]*"
	port    = "8080"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/jobs/run", auth.AuthenticateUser(worker.Run))
	router.HandleFunc("/jobs/{id:"+idRegex+"}/status", auth.AuthenticateUser(worker.Status))
	router.HandleFunc("/jobs/{id:"+idRegex+"}/out", auth.AuthenticateUser(worker.Out))
	router.HandleFunc("/jobs/{id:"+idRegex+"}/kill", auth.AuthenticateUser(worker.Kill))

	http.Handle("/", router)

	fmt.Println("Listening...")
	http.ListenAndServe(":"+port, nil)
}
