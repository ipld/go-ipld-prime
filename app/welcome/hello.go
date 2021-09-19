/*
	The welcome command package contains both some instruction information commands,
	and is also a demo about how we put the CLI system together and used for some tests.
*/
package welcome

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Cmd_Hello = &cli.Command{
	Name:     "hello",
	Category: "Welcome",
	Usage:    "Be greeted",
	Action: func(args *cli.Context) error {
		fmt.Fprintf(args.App.Writer, "hello!")
		return nil
	},
}
