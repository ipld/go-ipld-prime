package gengo

import (
	"io"
)

func emitMinima(f io.Writer) {
	emitFileHeader(f)

	// Iterator rejection thunks.
	f.Write([]byte(`
type mapIteratorReject struct{ err error }
type listIteratorReject struct{ err error }

func (itr mapIteratorReject) Next() (ipld.Node, ipld.Node, error) { return nil, nil, itr.err }
func (itr mapIteratorReject) Done() bool                          { return false }

func (itr listIteratorReject) Next() (int, ipld.Node, error) { return -1, nil, itr.err }
func (itr listIteratorReject) Done() bool                    { return false }
`))
}
