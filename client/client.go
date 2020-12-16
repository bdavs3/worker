package main

import (
	"log"
	"os"

	"github.com/bdavs3/worker/client/request"

	"github.com/urfave/cli/v2"
)

func main() {
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
				Action:  request.Run,
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "get the status of a job by providing its ID",
				Action:  request.Status,
			},
			{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "get the output of a job by providing its ID",
				Action:  request.Out,
			},
			{
				Name:    "kill",
				Aliases: []string{"k"},
				Usage:   "terminate a job by providing its ID",
				Action:  request.Kill,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
