package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bdavs3/worker/worker"
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
				t.Errorf("Error marshalling job as JSON: %v", err)
			}

			req, err := http.NewRequest(
				http.MethodPost,
				ts.URL,
				bytes.NewBuffer(requestBody),
			)
			if err != nil {
				t.Errorf("Error forming request: %v", err)
			}

			client := &http.Client{Timeout: 5 * time.Second}

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Did not receive response before timeout: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != test.want {
				t.Errorf("got %d, want %d", resp.StatusCode, test.want)
			}
		})
	}
}
