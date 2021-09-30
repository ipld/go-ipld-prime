`read` subcommand
=================

Docs
----

[testmark]:# (docs/sequence)
```
ipld read --help
```

[testmark]:# (docs/output)
```text
NAME:
   basic.test read - Read and print out one block of data (or a specific part of it, if a path is used); optionally, the data can be transcoded on the way out.

USAGE:
   Read is for inspecting data.

   ### Synopsis

   ipld [...global args...] read <CID|filename|"-"> [<datamodel-path>]
           [--output=<"debug"|"raw"|"html"|"codec:"<multicodec-name-or-hex>>]
           [--input="codec:"<multicodec-name-or-hex>]
           [--schema=<filename>|--schema-cid=<CID> --type=<starting-typename> [--schema-lens=<"representation"|"typed">] [--path-mode=<"representation"|"typed">]]
           [--ADL=<adlhook>]

   ### Data Sources

   The first positional argument (which is required) tells the read command what the data source is.
   Read can load data from storage if given a CID.
   Read can work with data passed in on stdin (stated by using a dash ("-") as the parameter).
   Read can consume data from a file (stated by using a "./" or "/" prefix, to disambiguate it from a CID!).

   When the input is from stdin or a file, the codec can be specified with the "--input" flag.  If it is not specified, a very simple heuristic will be used.  (This heuristic may change over time, and you should not rely on its behavior for noninteractive scripts.)

   ### Output Formats

   The default output format is a diagnostic printout format, meant for human readability.  Other formats and codecs can be specified (you'll probably want to do this if constructing some data pipeline; the diagnostic format is not meant to be parsed).

   Any IPLD codec can be used, by saying "--output=codec:<multicodec-name>" or "--output=codec:0x<multicodec-hex>".

   Another special output format, activated with "--output=raw", can be used in order to get the original raw serial stream, directly as it was loaded.  In this case, no codec is used at all, and the data is not validated or mutated in any way.  (Raw mode does not stack with most other features, including pathing.)

   An HTML output can be produced with "--output=html", which has similar purpose to the default textual debug format, but may include clickable links, etc.

   ### Transformations

   Several kinds of very basic transformation and filtering can be performed with additional options to the read command.

   If a path is provided as the second positional argument, after the data is loaded, it is traversed according to the path, and only the reached data will be emitted.

   If a schema is provided (either as a document in another file, or as a CID to be loaded from storage), it will be used to validate the data.  The name of the type in the schema that we expect to see at the root of the document must also be provided.  The output will default to the typed view, as with the pathing mode if is path parameter was provided, but both can be switched back to representation mode if desired by use of additional flags.

   Specifying a single ADL transformation to use will be supported in the future.  The API for this is not yet finalized.

   ### Multiple Blocks

   The read command is for handling one block of data at a time.  The read command does not support compositing a view of data taken from across multiple blocks.

   However, do note two features of the read command may still trigger block loading in the course of their work: Pathing may traverse links, and ADLs may also produce views of data which has involved link loading.


CATEGORY:
   Basic

OPTIONS:
   --output value  Defines what format the output should use.  Valid arguments are "debug", "raw", "html", or the word "codec:" followed by a multicodec name, or "codec:0x" followed by a multicodec indicator number in hexidecimal. (default: debug)
   --input value   Defines what format the input should be expected to be in.  Only relevant in the input is from a file or stdin; if the data source is a CID, that already implies a codec.  Valid arguments must start with "codec:" followed by a multicodec name, or "codec:0x" followed by a multicodec indicator number in hexidecimal.
   --help, -h      show help (default: false)
   
```

Simple Operations
-----------------

### Hello, read

Let's start with a very simple example of using the `ipld read` command:
piping some data into the command, and getting the debug representation of the same data printed back out:

[testmark]:# (hello-read/script)
```bash
echo '{"hello": "world"}' | ipld read -
```

The dash ("`-`") argument here meant "read from stdin".

What's printed out is the same data, in a debug representation:

[testmark]:# (hello-read/output)
```text
map{
	string{"hello"}: string{"world"}
}
```

### Pathing

We can also "path" through data -- stepping through a structure and looking at only a point that interests us.

This is done with a second argument:

[testmark]:# (hello-path/script)
```bash
echo '{"hello": {"pathing": "world"}}' | ipld read - hello/pathing
```

As you probably expected: asking the `ipld read` command to step through the "hello" and "pathing" keys in that object
reached the value (the string `"world"`):

[testmark]:# (hello-path/output)
```text
string{"world"}
```

### Raw passthrough mode

The read command can be operated in a "raw" mode, in which it returns whatever data it loads, without modification.

[testmark]:# (hello-raw/script)
```bash
echo '{"hello": "world"}' | ipld read --output=raw -
```

(Caution: the args parser requests `--arguments` to become before positional arguments.)

[testmark]:# (hello-raw/output)
```text
{"hello": "world"}
```

This example is someone mundane, because we're passing data into stdin already, so the read command is doing essentially nothing.
However, this can also be used with other forms of input, like a link for loading, which is more useful.
