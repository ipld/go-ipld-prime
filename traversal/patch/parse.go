package patch

import (
	_ "embed"

	"bytes"
	"io"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"

	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
)

//go:embed patch.ipldsch
var embedSchema []byte

var ts = func() *schema.TypeSystem {
	ts, err := ipld.LoadSchemaBytes(embedSchema)
	if err != nil {
		panic(err)
	}
	return ts
}()

func ParseBytes(b []byte, dec codec.Decoder) ([]Operation, error) {
	return Parse(bytes.NewReader(b), dec)
}

func Parse(r io.Reader, dec codec.Decoder) ([]Operation, error) {
	npt := bindnode.Prototype((*[]operationRaw)(nil), ts.TypeByName("OperationSequence"))
	nb := npt.Representation().NewBuilder()
	if err := json.Decode(nb, r); err != nil {
		return nil, err
	}
	opsRaw := bindnode.Unwrap(nb.Build()).(*[]operationRaw)
	var ops []Operation
	for _, opRaw := range *opsRaw {
		// TODO check the Op string
		op := Operation{
			Op:    Op(opRaw.Op),
			Path:  datamodel.ParsePath(opRaw.Path),
			Value: opRaw.Value,
		}
		if opRaw.From != nil {
			op.From = datamodel.ParsePath(*opRaw.From)
		}
		ops = append(ops, op)
	}
	return ops, nil
}

// operationRaw is roughly the same structure as Operation, but more amenable to serialization
// (it doesn't use high level library types that don't have a data model equivalent).
type operationRaw struct {
	Op    string
	Path  string
	Value datamodel.Node
	From  *string
}
