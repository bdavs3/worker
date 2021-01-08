package client

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

	"github.com/bdavs3/worker/server/api"
	"github.com/bdavs3/worker/worker"

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
	timeout = 5 * time.Second
)

// Client is an HTTP client for making requests to control and assess Linux
// processes on the server. Use NewClient to create a new instance.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Client instance.
func NewClient() *Client {
	return &Client{
		BaseURL:    host,
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

// PostJob passes a Linux process to the worker library for execution.
func (c *Client) PostJob(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("no job supplied to 'run' command")
	}

	job := worker.Job{
		Command: ctx.Args().Get(0),
		Args:    ctx.Args().Slice()[1:],
	}

	requestBody, err := json.Marshal(job)
	if err != nil {
		return err
	}

	response, err := c.makeRequestWithAuth(
		http.MethodPost,
		"/jobs/run",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	fmt.Println(response.ID)

	return nil
}

// GetJobStatus queries the status of a process being handled by the worker library.
func (c *Client) GetJobStatus(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job ID supplied to 'status' command")
	}

	jobID := ctx.Args().Get(0)

	response, err := c.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/status", jobID),
		nil,
	)
	if err != nil {
		return err
	}
	if response.ID != jobID {
		return errors.New("response contains incorrect job id")
	}

	fmt.Println(response.Status)

	return nil
}

// GetJobOutput queries the output of a process being handled by the worker library.
func (c *Client) GetJobOutput(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job ID supplied to 'out' command")
	}

	jobID := ctx.Args().Get(0)

	response, err := c.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/out", jobID),
		nil,
	)
	if err != nil {
		return err
	}
	if response.ID != jobID {
		return errors.New("response contains incorrect job id")
	}

	fmt.Print(response.Output)

	return nil
}

// KillJob terminates a process being handled by the worker library.
func (c *Client) KillJob(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job ID supplied to 'kill' command")
	}

	jobID := ctx.Args().Get(0)

	response, err := c.makeRequestWithAuth(
		http.MethodPut,
		fmt.Sprintf("/jobs/%s/kill", jobID),
		nil,
	)
	if err != nil {
		return err
	}
	if response.ID != jobID {
		return errors.New("response contains incorrect job id")
	}

	fmt.Println(response.Status)

	return nil
}

// makeRequestWithAuth makes an HTTP request to the given endpoint
// after setting the Authorization header. It then returns the response.
func (c *Client) makeRequestWithAuth(method, endpoint string, requestBody io.Reader) (*api.Response, error) {
	req, err := http.NewRequest(
		method,
		c.BaseURL+endpoint,
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("username"), os.Getenv("pw"))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s\n%s", http.StatusText(resp.StatusCode), body)
	}

	var response *api.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
