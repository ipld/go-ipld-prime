package shared

import (
	"bufio"
	"os"
	"strings"

	"github.com/ipld/go-ipld-prime/datamodel"
)

// StringIsPathish returns true if the string explicitly looks like a filesystem path
// (starts with `./`, `../`, or `/`).
func StringIsPathish(x string) bool {
	return strings.HasPrefix(x, "./") ||
		strings.HasPrefix(x, "../") ||
		strings.HasPrefix(x, "/")
}

// ParseDataSourceArg returns a reader for data based on the argument,
// and a Link if the argument was of that kind.
func ParseDataSourceArg(inputArg string) (reader *bufio.Reader, link datamodel.Link, err error) {
	switch {
	case inputArg == "-": // stdin
		reader = bufio.NewReader(os.Stdin) // FIXME does this cli package not have a way to attach a stream so I don't have to use a global for this?
	case StringIsPathish(inputArg): // looks like a filename
		f, err := os.Open(inputArg)
		if err != nil {
			return nil, nil, err
		}
		reader = bufio.NewReader(f)
	default: // hope this is a CID
		panic("todo")
	}
	return
}
