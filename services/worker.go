package services

import (
	"errors"
	"fmt"

	"github.com/bdavs3/worker/client"
	"github.com/bdavs3/worker/worker"

	"github.com/urfave/cli/v2"
)

type WorkerService struct {
	Client *client.Client
}

func NewWorkerService() (*WorkerService, error) {
	client, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	workerService := &WorkerService{
		Client: client,
	}

	return workerService, nil
}

func (ws *WorkerService) Run(ctx *cli.Context) error {
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

func (ws *WorkerService) Status(ctx *cli.Context) error {
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

func (ws *WorkerService) Out(ctx *cli.Context) error {
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

func (ws *WorkerService) Kill(ctx *cli.Context) error {
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
