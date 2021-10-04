`schema` subcommand
===================

Docs
----

[testmark]:# (docs/script)
```
ipld schema --help
```

[testmark]:# (docs/output)
```text
NAME:
   ipld schema - Manipulate schemas -- parsing, compiling, transforming, and storing.

USAGE:
   ipld schema command [command options] [arguments...]

COMMANDS:
   parse    Parse a schema DSL document, and produce the DMT form, emitted in JSON by default.
   compile  Compile a schema DMT document, exiting nonzero and reporting errors if anything is logically invalid.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
   
```

Parsing
-------

### Hello, parse

To demonstrate the `ipld schema parse` command, we'll first need a small schema document (in the DSL format).

We'll put this in a file called "`theschema.ipldsch`":

[testmark]:# (hello-parse/fs/theschema.ipldsch)
```ipldsch
type Hello string

type World struct {
	field Hello
}
```

The parse command simply takes the usual kinds of inputs: a link, or a filename, or "-" for standard input.
Here we use the name of the file with the contents above:

[testmark]:# (hello-parse/script)
```bash
ipld schema parse ./theschema.ipldsch
```

This prints out the parsed DMT form of the schema:

[testmark]:# (hello-parse/output)
```text
{
	"types": {
		"Hello": {
			"string": {}
		},
		"World": {
			"struct": {
				"fields": {
					"field": {
						"type": "Hello"
					}
				},
				"representation": {
					"map": {
						"fields": {}
					}
				}
			}
		}
	}
}
```
