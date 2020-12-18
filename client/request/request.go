package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bdavs3/worker/server/worker"

	"github.com/urfave/cli/v2"
)

// TODO (next): Change host protocol to 'https' and port to 443 once API
// serves with TLS.

// TODO (out of scope): Rather than using hard-coded user credentials, provide
// the user with a way to create an account and log in. Once authenticated
// with the API, the client could receive a session token that precludes the
// need to authenticate on each subsequent request.

const host = "http://localhost:8080"

// Run passes a job to the worker library for execution.
func Run(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("No job supplied to 'run' command")
	}

	job := worker.Job{
		Command: c.Args().Get(0),
		Args:    c.Args().Slice()[1:],
	}

	requestBody, err := json.Marshal(job)
	if err != nil {
		return err
	}

	responseBody, err := makeRequestWithAuth(
		http.MethodPost,
		host+"/jobs/run",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// Status queries the status of a job being handled by the worker library.
func Status(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'status' command")
	}
	if _, err := strconv.Atoi(c.Args().Get(0)); err != nil {
		return errors.New("Job ID must be an integer")
	}

	jobID := c.Args().Get(0)

	responseBody, err := makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("%s/jobs/%s/status", host, jobID),
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// Out queries the output of a job being handled by the worker library.
func Out(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'out' command")
	}
	if _, err := strconv.Atoi(c.Args().Get(0)); err != nil {
		return errors.New("Job ID must be an integer")
	}

	jobID := c.Args().Get(0)

	responseBody, err := makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("%s/jobs/%s/out", host, jobID),
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

// Kill terminates a job being handled by the worker library.
func Kill(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'kill' command")
	}
	if _, err := strconv.Atoi(c.Args().Get(0)); err != nil {
		return errors.New("Job ID must be an integer")
	}

	jobID := c.Args().Get(0)

	responseBody, err := makeRequestWithAuth(
		http.MethodPut,
		fmt.Sprintf("%s/jobs/%s/kill", host, jobID),
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
func makeRequestWithAuth(method, endpoint string, requestBody io.Reader) (string, error) {
	req, err := http.NewRequest(method, endpoint, requestBody)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(os.Getenv("username"), os.Getenv("pw"))

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
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
