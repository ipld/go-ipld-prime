package gengo

import (
	"io"
)

func emitMinima(f io.Writer) {
	emitFileHeader(f)

	// Avoid moderate annoyance keeping track of imports.
	f.Write([]byte(`
var _ typed.Node = typed.Node(nil)
`))

	// Iterator rejection thunks.
	f.Write([]byte(`
type mapIteratorReject struct{ err error }
type listIteratorReject struct{ err error }

func (itr mapIteratorReject) Next() (ipld.Node, ipld.Node, error) { return nil, nil, itr.err }
func (itr mapIteratorReject) Done() bool                          { return false }

func (itr listIteratorReject) Next() (int, ipld.Node, error) { return -1, nil, itr.err }
func (itr listIteratorReject) Done() bool                    { return false }
`))

	// Box type for map keys.
	// f.Write([]byte(`
	// type boxedString struct { x string }
	// `))
	//
	// ... nevermind; we already need strings in the prelude.  Use em.
}
