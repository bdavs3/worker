package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/bdavs3/worker/server/worker"
	"github.com/urfave/cli/v2"
)

// TODO: Change to 'https' and port 443 once API serves with TLS.
const (
	host = "http://localhost"
	port = "8080"
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

	resp, err := http.Post(
		host+":"+port+"/jobs/run",
		"application/json",
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

// Status queries the status of job being handled by the worker library.
func Status(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("No job ID supplied to 'status' command")
	}
	if _, err := strconv.Atoi(c.Args().Get(0)); err != nil {
		return errors.New("Job ID must be an integer")
	}

	// TODO: Pass job ID as path paramater once API is configured to accept it.
	resp, err := http.Get(host + ":" + port + "/jobs/status")
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

func Out(c *cli.Context) error {
	fmt.Println("Out")
	return nil
}

func Kill(c *cli.Context) error {
	fmt.Println("Kill")
	return nil
}
