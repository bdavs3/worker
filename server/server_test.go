package main

import (
	"bytes"
	"encoding/json"
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
		comment string
		job     worker.Job
		want    int
	}{
		{
			comment: "well-formed request to /jobs/run",
			job:     worker.Job{Command: "echo", Args: []string{"hello"}},
			want:    http.StatusOK,
		},
		{
			comment: "poorly-formed request to /jobs/run",
			job:     worker.Job{Command: "", Args: []string{}},
			want:    http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.comment, func(t *testing.T) {
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
