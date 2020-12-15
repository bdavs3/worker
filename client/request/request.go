package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bdavs3/worker/server/worker"
	"github.com/urfave/cli/v2"
)

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

func Status(c *cli.Context) error {
	fmt.Println("Status")
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
