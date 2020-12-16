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

	router.HandleFunc("/jobs/run", worker.Run)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/status", worker.Status)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/out", worker.Out)
	router.HandleFunc("/jobs/{id:"+idRegex+"}/kill", worker.Kill)

	fmt.Println("Listening...")
	http.ListenAndServe(":"+port, auth.Secure(router))
}
