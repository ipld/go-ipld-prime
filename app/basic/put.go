package basic

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Cmd_Put = &cli.Command{
	Name:     "put",
	Category: "Basic",
	Usage:    "Put a single block of data into storage.",
	Action: func(args *cli.Context) error {
		return fmt.Errorf("not yet implemented")
	},
}
