package gengo

import (
	"io"
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/schema"
)

func TestNuevo(t *testing.T) {
	os.Mkdir("_test", 0755)
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}

	emitType := func(tg typeGenerator, w io.Writer) {
		tg.EmitNodeType(w)
		tg.EmitNodeMethodReprKind(w)
		tg.EmitNodeMethodTraverseField(w)
		tg.EmitNodeMethodTraverseIndex(w)
		tg.EmitNodeMethodMapIterator(w)
		tg.EmitNodeMethodListIterator(w)
		tg.EmitNodeMethodLength(w)
		tg.EmitNodeMethodIsNull(w)
		tg.EmitNodeMethodAsBool(w)
		tg.EmitNodeMethodAsInt(w)
		tg.EmitNodeMethodAsFloat(w)
		tg.EmitNodeMethodAsString(w)
		tg.EmitNodeMethodAsBytes(w)
		tg.EmitNodeMethodAsLink(w)
		tg.EmitNodeMethodNodeBuilder(w)
	}

	f := openOrPanic("_test/thunks.go")
	emitFileHeader(f)
	doTemplate(`
		type mapIteratorReject struct{ err error }
		type listIteratorReject struct{ err error }

		func (itr mapIteratorReject) Next() (ipld.Node, ipld.Node, error) { return nil, nil, itr.err }
		func (itr mapIteratorReject) Done() bool                          { return false }

		func (itr listIteratorReject) Next() (int, ipld.Node, error) { return -1, nil, itr.err }
		func (itr listIteratorReject) Done() bool                    { return false }
	`, f, nil)

	f = openOrPanic("_test/tStrang.go")
	emitFileHeader(f)
	emitType(generateKindString{"Strang", schema.TypeString{}}, f)
}
