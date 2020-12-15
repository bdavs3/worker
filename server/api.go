package main

import (
	"net/http"

	"github.com/bdavs3/worker/server/worker"
)

const port = "8080"

func main() {
	http.HandleFunc("/jobs/run", worker.Run)
	http.HandleFunc("/jobs/status", worker.Status)
	http.HandleFunc("/jobs/out", worker.Out)
	http.HandleFunc("/jobs/kill", worker.Kill)

	http.ListenAndServe(":"+port, nil)
}
