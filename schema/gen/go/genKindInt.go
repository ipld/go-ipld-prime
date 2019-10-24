package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindInt(t schema.Type) typedNodeGenerator {
	return generateKindInt{
		t.(schema.TypeInt),
		generateKindedRejections_Int{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindInt struct {
	Type schema.TypeInt
	generateKindedRejections_Int
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindInt) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

		type {{ .Type | mungeTypeNodeIdent }} struct{ x int }

	`, w, gk)
}

func (gk generateKindInt) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindInt) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Int
		}
	`, w, gk)
}

func (gk generateKindInt) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) AsInt() (int, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindInt) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}

func (gk generateKindInt) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
