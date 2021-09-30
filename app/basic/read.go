package basic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/printer"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/traversal"
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
		`   Another special output format, activated with "--output=raw", can be used in order to get the original raw serial stream, directly as it was loaded.  In this case, no codec is used at all, and the data is not validated or mutated in any way.  (Raw mode does not stack with most other features, including pathing.)` + "\n" +
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
		&cli.StringFlag{
			Name:  "input",
			Usage: `Defines what format the input should be expected to be in.  Only relevant in the input is from a file or stdin; if the data source is a CID, that already implies a codec.  Valid arguments must start with "codec:" followed by a multicodec name, or "codec:0x" followed by a multicodec indicator number in hexidecimal.`,
		},
		// TODO more
	},
	Action: func(args *cli.Context) error {
		// Parse positional args.
		var sourceArg string
		var pathArg string
		switch args.Args().Len() {
		case 2:
			pathArg = args.Args().Get(1)
			fallthrough
		case 1:
			sourceArg = args.Args().Get(0)
		default:
			fmt.Errorf("read command needs one or two positional arguments")
		}

		// Let's get some data!
		var reader *bufio.Reader
		var link datamodel.Link
		switch {
		case sourceArg == "-": // stdin
			reader = bufio.NewReader(os.Stdin) // FIXME does this cli package not have a way to attach a stream so I don't have to use a global for this?
		case looksPathish(sourceArg): // looks like a filename
			panic("todo")
		default: // hope this is a CID
			panic("todo")
		}

		// Early exit: if "raw" mode is requested, pass the data through direction.  Skip *everything* else.  (No need to determine codec, nothing.)
		//  (Future: maybe we can path.  However, it would only work as long as the lands on a block edge.  Unclear how useful this would be; PRs welcome.)
		if args.String("output") == "raw" {
			_, err := io.Copy(args.App.Writer, reader)
			return err
		}

		// Determine the input codec.
		//  This can involve peeking at the bytes, if there's no explicit statements.
		// The dominance is:
		//  1. Listen to the flag, if there is one.
		//  2. Listen to the CID, if that was the data source.
		//  3. Peek and guess as a last resort.
		//  4. If we can't guess usefully, give up and error.
		var decoder codec.Decoder
		switch {
		case args.IsSet("input"):
			panic("todo")
		case link != nil:
			panic("todo")
		default:
			peeked, err := reader.Peek(1)
			if err != nil {
				return err
			}
			switch peeked[0] {
			case '{': // it's probably json.  we'll assume dag-json.
				decoder = dagjson.Decode
			default:
				return fmt.Errorf("no input codec specified by args, and gave up trying to guess one from the content")
			}
		}

		// Was there a schema?  Load that, compile it, and get a NodePrototype from that.
		// Otherwise?  Basicnode will do.
		np := basicnode.Prototype.Any

		// Was there an ADL hint?  Haven't implemented that yet,
		//  but we'd probably at least parse it here.

		// Figure out the output format too.
		//  We don't need this yet, but it's good practice to at make sure all the args are sane and we know what to do with them before starting real work.
		var encoder codec.Encoder
		switch args.String("output") {
		case "", "debug":
			encoder = func(n datamodel.Node, wr io.Writer) error {
				printer.Fprint(wr, n)
				return nil
			}
		case "raw":
			panic("unreachable (already handled this case earlier)")
		case "html":
			panic("todo")
		default:
			switch {
			case strings.HasPrefix(args.String("output"), "codec:0x"):
				panic("todo")
			case strings.HasPrefix(args.String("output"), "codec:"):
				panic("todo")
			default:
				return fmt.Errorf("output argument format not recognized")
			}
		}

		// Finally, we have the codec, the input stream, and the NodePrototype.
		// And all the other args-parsing we'll need by the end is done too.
		// Let's go!
		n, err := ipld.DecodeStreamingUsingPrototype(reader, decoder, np)
		if err != nil {
			return err
		}

		// Pathing time!
		//  Drop back down to representation level first, too, if we had types, and also the flag requesting representation-level pathing.
		if tn, ok := n.(schema.TypedNode); ok && args.String("path-mode") == "representation" {
			n = tn.Representation()
		}
		n, err = traversal.Get(n, datamodel.ParsePath(pathArg))
		if err != nil {
			return err
		}

		// TODO: can we... actually readily switch back up to type-view after pathing at repr level?  I feel like that should be possible (at least in some cases; not all).

		// Finally: print back out whatever we've read (and possibly transformed, and pathed to).
		err = ipld.EncodeStreaming(args.App.Writer, n, encoder)

		// And push one last trailing linebreak out, because that's considered a normative ending thing in most CLI composition.
		args.App.Writer.Write([]byte{'\n'})

		return err
	},
}

func looksPathish(x string) bool {
	return strings.HasPrefix(x, "./") ||
		strings.HasPrefix(x, "../") ||
		strings.HasPrefix(x, "/")
}
