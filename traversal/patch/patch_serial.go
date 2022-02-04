package patch

import (
	"io"
	"strings"

	"github.com/ipld/go-ipld-prime/node/bindnode"

	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
)

var ts = func() schema.TypeSystem {
	sch, err := schemadsl.Parse("", strings.NewReader(
		// This could be more accurately modelled as an inline union,
		// but that seems like work, given how high the overlap is.
		`
		type Operation struct {
			op String
			path String
			value optional Any
			from optional String
		}
		type OperationList [Operation]
	`))
	if err != nil {
		panic(err)
	}
	var ts schema.TypeSystem
	ts.Init()
	if err := schemadmt.Compile(&ts, sch); err != nil {
		panic(err)
	}
	return ts
}()

// FIXME this should surely accept a codec.Decoder parameter
func LoadPatch(r io.Reader) ([]Operation, error) {
	npt := bindnode.Prototype((*[]operationRaw)(nil), ts.TypeByName("OperationList"))
	nb := npt.Representation().NewBuilder()
	if err := json.Decode(nb, r); err != nil {
		return nil, err
	}
	opsRaw := bindnode.Unwrap(nb.Build()).(*[]operationRaw)
	var ops []Operation
	for _, opRaw := range *opsRaw {
		// TODO check the Op string
		op := Operation{
			Op:   Op(opRaw.Op),
			Path: datamodel.ParsePath(opRaw.Path),
		}
		if opRaw.Value != nil {
			op.Value = *opRaw.Value
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
	Value *datamodel.Node
	From  *string
}
