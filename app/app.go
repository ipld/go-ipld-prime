package app

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ipld/go-ipld-prime/app/basic"
	"github.com/ipld/go-ipld-prime/app/welcome"
)

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	app := &cli.App{
		Name:      "ipld",
		Usage:     "a data wangling and mangling tool, for munging and wunging, yurling and curling",
		Writer:    stdout,
		ErrWriter: stderr,
		Commands: []*cli.Command{
			welcome.Cmd_Hello,
			basic.Cmd_Put,
			basic.Cmd_Read,
		},
	}

	err := app.Run(args)
	if err == nil {
		return 0, nil
	}
	// Future: some kind of routing table of error code to exit code.
	//  (You can ignore the exit code, and still just look at the error.
	//   This will be useful to have because we'll also have a web daemon mode which calls mostly the same functions,
	//    but will have to route the same error codes into the different int space of HTTP status codes.)
	fmt.Fprintf(os.Stderr, "error: %s", err)
	return 1, err
}
