package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindBytes(t schema.Type) typedNodeGenerator {
	return generateKindBytes{
		t.(schema.TypeBytes),
		generateKindedRejections_Bytes{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindBytes struct {
	Type schema.TypeBytes
	generateKindedRejections_Bytes
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindBytes) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

		type {{ .Type | mungeTypeNodeIdent }} struct{ x []byte }

	`, w, gk)
}

func (gk generateKindBytes) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindBytes) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Bytes
		}
	`, w, gk)
}

func (gk generateKindBytes) EmitNodeMethodAsBytes(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) AsBytes() ([]byte, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindBytes) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}

func (gk generateKindBytes) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
