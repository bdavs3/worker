package request

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	fmt.Println("Run")
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
