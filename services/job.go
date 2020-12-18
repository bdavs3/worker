package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// TODO (next): Change host protocol to 'https' and port to 443 once API
// serves with TLS.

// TODO (out of scope): Rather than using hard-coded user credentials, provide
// the user with a way to create an account and log in. Once authenticated
// with the API, the client could receive a session token that precludes the
// need to authenticate on each subsequent request.

const (
	host    = "http://localhost:8080"
	timeout = 5
)

// Job represents a Linux process to be passed to the API.
type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// JobService contains a BaseURL to prepend to API endpoints and an HTTPClient
// for making requests to those endpoints.
type JobService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewJobService creates a JobService struct with a default BaseURL and HTTPClient.
func NewJobService() *JobService {
	return &JobService{
		BaseURL:    host,
		HTTPClient: &http.Client{Timeout: timeout * time.Second},
	}
}

// PostJob passes a job to the worker library for execution.
func (js *JobService) PostJob(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("No job supplied to 'run' command")
	}

	job := Job{
		Command: c.Args().Get(0),
		Args:    c.Args().Slice()[1:],
	}

	requestBody, err := json.Marshal(job)
	if err != nil {
		return err
	}

	responseBody, err := js.makeRequestWithAuth(
		http.MethodPost,
		"/jobs/run",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// GetJobStatus queries the status of a job being handled by the worker library.
func (js *JobService) GetJobStatus(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'status' command")
	}

	jobID := c.Args().Get(0)

	responseBody, err := js.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/status", jobID),
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// GetJobOutput queries the output of a job being handled by the worker library.
func (js *JobService) GetJobOutput(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'out' command")
	}

	jobID := c.Args().Get(0)

	responseBody, err := js.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/out", jobID),
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// KillJob terminates a job being handled by the worker library.
func (js *JobService) KillJob(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'kill' command")
	}

	jobID := c.Args().Get(0)

	responseBody, err := js.makeRequestWithAuth(
		http.MethodPut,
		fmt.Sprintf("/jobs/%s/kill", jobID),
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// makeRequestWithAuth makes an HTTP request to a given endpoint after setting
// the request's Authorization header. It reads the response body and returns
// it as a string.
func (js *JobService) makeRequestWithAuth(method, endpoint string, requestBody io.Reader) (string, error) {
	req, err := http.NewRequest(method, js.BaseURL+endpoint, requestBody)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(os.Getenv("username"), os.Getenv("pw"))

	resp, err := js.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
