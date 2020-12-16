package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bdavs3/worker/server/worker"
)

func TestAPIRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(worker.Run))

	defer ts.Close()

	var tests = []struct {
		job  worker.Job
		want int
	}{
		{worker.Job{Command: "echo", Args: []string{"hello"}}, 200},
		{worker.Job{Command: "", Args: []string{}}, 400},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%+v, %d", test.job, test.want)
		t.Run(testname, func(t *testing.T) {
			requestBody, err := json.Marshal(test.job)
			if err != nil {
				log.Fatal(err)
			}

			res, err := http.Post(
				ts.URL,
				"application/json",
				bytes.NewBuffer(requestBody),
			)
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != test.want {
				t.Errorf("got %d, want %d", res.StatusCode, test.want)
			}
		})
	}
}
