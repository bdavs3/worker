package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bdavs3/worker/client"
	"github.com/bdavs3/worker/worker"

	"github.com/urfave/cli/v2"
)

func main() {
	workerService, err := newWorkerService()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "worker",
		Usage: "run arbitrary Linux processes",
		Commands: []*cli.Command{
			// By default, the CLI includes a "help" command that displays app
			// info and command usage.
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "give the server a Linux process to execute",
				Action:  workerService.run,
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "get the status of a process by providing its id",
				Action:  workerService.status,
			},
			{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "get the output of a process by providing its id",
				Action:  workerService.out,
			},
			{
				Name:    "kill",
				Aliases: []string{"k"},
				Usage:   "terminate a process by providing its id",
				Action:  workerService.kill,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type workerService struct {
	Client *client.Client
}

func newWorkerService() (*workerService, error) {
	client, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	workerService := &workerService{
		Client: client,
	}

	return workerService, nil
}

func (ws *workerService) run(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("no job supplied to 'run' command")
	}

	job := worker.Job{
		Command: ctx.Args().Get(0),
		Args:    ctx.Args().Slice()[1:],
	}

	responseBody, err := ws.Client.PostJob(job)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

func (ws *workerService) status(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job id supplied to 'status' command")
	}

	id := ctx.Args().Get(0)

	responseBody, err := ws.Client.GetJobStatus(id)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}

func (ws *workerService) out(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job ID supplied to 'out' command")
	}

	id := ctx.Args().Get(0)

	responseBody, err := ws.Client.GetJobOutput(id)
	if err != nil {
		return err
	}

	fmt.Print(responseBody)

	return nil
}

func (ws *workerService) kill(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("no job ID supplied to 'kill' command")
	}

	id := ctx.Args().Get(0)

	responseBody, err := ws.Client.KillJob(id)
	if err != nil {
		return err
	}

	fmt.Println(responseBody)

	return nil
}
