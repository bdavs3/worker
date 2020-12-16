package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

const (
	host     = "http://localhost"
	port     = "8080"
	username = "default_user"
	pw       = "123456"
)

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
		log.Fatal(err)
	}

	resp, err := makeRequestWithAuth(
		http.MethodPost,
		host+":"+port+"/jobs/run",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

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

	resp, err := makeRequestWithAuth(
		http.MethodGet,
		host+":"+port+"/jobs/"+jobID+"/status",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

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

	resp, err := makeRequestWithAuth(
		http.MethodGet,
		host+":"+port+"/jobs/"+jobID+"/out",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

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

	resp, err := makeRequestWithAuth(
		http.MethodPut,
		host+":"+port+"/jobs/"+jobID+"/kill",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	return nil
}

// makeRequestWithAuth makes an HTTP request to a given endpoint after setting
// the request's Authorization header.
func makeRequestWithAuth(method, endpoint string, body io.Reader) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(username, pw)

	return client.Do(req)
}
