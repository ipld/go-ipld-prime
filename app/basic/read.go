package basic

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Cmd_Read = &cli.Command{
	Name:     "read",
	Category: "Basic",
	Usage:    "Read and print out one block of data (or a specific part of it, if a path is used); optionally, the data can be transcoded on the way out.",
	UsageText: `Read is for inspecting data.` + "\n" +
		"\n" +
		`   ### Synopsis` + "\n" +
		"\n" +
		`   ipld [...global args...] read <CID|filename|"-"> [<datamodel-path>]` + "\n" +
		`           [--output=<"debug"|"raw"|"html"|"codec:"<multicodec-name-or-hex>>]` + "\n" +
		`           [--input="codec:"<multicodec-name-or-hex>]` + "\n" +
		`           [--schema=<filename>|--schema-cid=<CID> --type=<starting-typename> [--schema-lens=<"representation"|"typed">] [--path-mode=<"representation"|"typed">]]` + "\n" +
		`           [--ADL=<adlhook>]` + "\n" +
		"\n" +
		`   ### Data Sources` + "\n" +
		"\n" +
		`   The first positional argument (which is required) tells the read command what the data source is.` + "\n" +
		`   Read can load data from storage if given a CID.` + "\n" +
		`   Read can work with data passed in on stdin (stated by using a dash ("-") as the parameter).` + "\n" +
		`   Read can consume data from a file (stated by using a "./" or "/" prefix, to disambiguate it from a CID!).` + "\n" +
		"\n" +
		`   When the input is from stdin or a file, the codec can be specified with the "--input" flag.  If it is not specified, a very simple heuristic will be used.  (This heuristic may change over time, and you should not rely on its behavior for noninteractive scripts.)` + "\n" +
		"\n" +
		`   ### Output Formats` + "\n" +
		"\n" +
		`   The default output format is a diagnostic printout format, meant for human readability.  Other formats and codecs can be specified (you'll probably want to do this if constructing some data pipeline; the diagnostic format is not meant to be parsed).` + "\n" +
		"\n" +
		`   Any IPLD codec can be used, by saying "--output=codec:<multicodec-name>" or "--output=codec:0x<multicodec-hex>".` + "\n" +
		"\n" +
		`   Another special output format, activated with "--output=raw", can be used in order to get the original raw serial stream, directly as it was loaded.  In this case, no codec is used at all, and the data is not validated or mutated in any way.` + "\n" +
		"\n" +
		`   An HTML output can be produced with "--output=html", which has similar purpose to the default textual debug format, but may include clickable links, etc.` + "\n" +
		"\n" +
		`   ### Transformations` + "\n" +
		"\n" +
		`   Several kinds of very basic transformation and filtering can be performed with additional options to the read command.` + "\n" +
		"\n" +
		`   If a path is provided as the second positional argument, after the data is loaded, it is traversed according to the path, and only the reached data will be emitted.` + "\n" +
		"\n" +
		`   If a schema is provided (either as a document in another file, or as a CID to be loaded from storage), it will be used to validate the data.  The name of the type in the schema that we expect to see at the root of the document must also be provided.  The output will default to the typed view, as with the pathing mode if is path parameter was provided, but both can be switched back to representation mode if desired by use of additional flags.` + "\n" +
		"\n" +
		`   Specifying a single ADL transformation to use will be supported in the future.  The API for this is not yet finalized.` + "\n" +
		"\n" +
		`   ### Multiple Blocks` + "\n" +
		"\n" +
		`   The read command is for handling one block of data at a time.  The read command does not support compositing a view of data taken from across multiple blocks.` + "\n" +
		"\n" +
		`   However, do note two features of the read command may still trigger block loading in the course of their work: Pathing may traverse links, and ADLs may also produce views of data which has involved link loading.` + "\n",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Usage:       `Defines what format the output should use.  Valid arguments are "debug", "raw", "html", or the word "codec:" followed by a multicodec name, or "codec:0x" followed by a multicodec indicator number in hexidecimal.`,
			DefaultText: "debug",
		},
		// TODO more
	},
	Action: func(args *cli.Context) error {
		return fmt.Errorf("not yet implemented")
	},
}
