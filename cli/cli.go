package main

import (
	"log"
	"os"

	"github.com/bdavs3/worker/client"

	"github.com/urfave/cli/v2"
)

func main() {
	client, err := client.NewClient()
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
				Action:  client.PostJob,
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "get the status of a services by providing its ID",
				Action:  client.GetJobStatus,
			},
			{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "get the output of a services by providing its ID",
				Action:  client.GetJobOutput,
			},
			{
				Name:    "kill",
				Aliases: []string{"k"},
				Usage:   "terminate a services by providing its ID",
				Action:  client.KillJob,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
