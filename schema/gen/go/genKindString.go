package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindString(t schema.Type) typedNodeGenerator {
	return generateKindString{
		t.(schema.TypeString),
		generateKindedRejections_String{t},
	}
}

type generateKindString struct {
	Type schema.TypeString
	generateKindedRejections_String
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type.Name }}{}
		var _ typed.Node = {{ .Type.Name }}{}

		type {{ .Type.Name }} struct{ x string }

	`, w, gk)
}

func (gk generateKindString) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (x {{ .Type.Name }}) AsString() (string, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindString) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}
