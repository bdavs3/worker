package worker

import (
	"fmt"
	"net/http"
)

type Job struct {
	Command string
	Args    []string
}

func Run(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Run")
}

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status")
}

func Out(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Out")
}

func Kill(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kill")
}
