package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bdavs3/worker/server/auth"
	"github.com/bdavs3/worker/worker"
)

func TestAPIRequest(t *testing.T) {
	// Using a DummyWorker allows the API to be tested without entangling it
	// in the functionality of the worker library.
	dummyWorker := &worker.DummyWorker{}
	auth := &auth.DummyAuth{}
	handler := NewHandler(dummyWorker, auth)

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
		rec := httptest.NewRecorder()

		t.Run(test.comment, func(t *testing.T) {
			requestBody, err := json.Marshal(test.job)
			if err != nil {
				t.Errorf("Error marshalling job as JSON: %v", err)
			}

			req, err := http.NewRequest(
				http.MethodPost,
				"/jobs/run",
				bytes.NewBuffer(requestBody),
			)
			if err != nil {
				t.Errorf("Error forming request: %v", err)
			}

			http.HandlerFunc(handler.PostJob).ServeHTTP(rec, req)

			if rec.Code != test.want {
				t.Errorf("got %d, want %d", rec.Code, test.want)
			}
		})
	}
}
