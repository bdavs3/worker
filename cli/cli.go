package main

import (
	"log"
	"os"

	"github.com/bdavs3/worker/services"

	"github.com/urfave/cli/v2"
)

func main() {
	workerService, err := services.NewWorkerService()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "worker",
		Usage: "run arbitrary Linux jobs",
		Commands: []*cli.Command{
			// By default, the CLI includes a "help" command that displays app
			// info and command usage.
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "give the server a Linux process to execute",
				Action:  workerService.Run,
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "get the status of a services by providing its id",
				Action:  workerService.Status,
			},
			{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "get the output of a services by providing its id",
				Action:  workerService.Out,
			},
			{
				Name:    "kill",
				Aliases: []string{"k"},
				Usage:   "terminate a services by providing its id",
				Action:  workerService.Kill,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
