package schema

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/ipld/go-ipld-prime/app/shared"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
)

var Cmd_Schema = &cli.Command{
	Name:     "schema",
	Category: "Advanced",
	Usage:    "Manipulate schemas -- parsing, compiling, transforming, and storing.",
	Subcommands: []*cli.Command{{
		Name:  "parse",
		Usage: "Parse a schema DSL document, and produce the DMT form, emitted in JSON by default.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-compile",
				Usage: `Skip the compilation phase, and just emit the DMT (regardless of whether it's logically valid).`,
			},
			&cli.BoolFlag{
				Name:  "save",
				Usage: `Put the parsed schema into storage, and return a CID pointing to it.  (Roughly equivalent to piping the schema parse command into a put command.)`,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       `Defines what format the DMT should be produced in.  Valid arguments are codecs, specified as the word "codec:" followed by a multicodec name, or "codec:0x" followed by a multicodec indicator number in hexidecimal.`,
				DefaultText: "codec:json",
			},
		},
		Action: func(args *cli.Context) error {
			// Parse positional args.
			var sourceArg string
			switch args.Args().Len() {
			case 1:
				sourceArg = args.Args().Get(0)
			default:
				fmt.Errorf("schema parse command needs exactly one positional argument")
			}

			// Let's get some data!
			inputReader, _ := shared.ParseDataSourceArg(sourceArg)

			// Parse!
			dmt, err := schemadsl.Parse(sourceArg, inputReader)
			if err != nil {
				return err // TODO probably need an error tagging strategy here.
			}

			// Compile!  Maybe.  Just to make sure we can.
			var ts schema.TypeSystem
			ts.Init()
			if err := schemadmt.Compile(&ts, dmt); err != nil {
				return err // TODO probably need an error tagging strategy here.
			}

			return nil
		},
	}, {
		Name:  "compile",
		Usage: "Compile a schema DMT document, exiting nonzero and reporting errors if anything is logically invalid.",
	}},
	// Someday: it may be neat to have a handful of well-known transforms, like: strip all rename directives, or make all representations default, etc.
}
