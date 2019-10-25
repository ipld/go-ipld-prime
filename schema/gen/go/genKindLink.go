package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindLink(t schema.Type) typedNodeGenerator {
	return generateKindLink{
		t.(schema.TypeLink),
		generateKindedRejections_Link{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindLink struct {
	Type schema.TypeLink
	generateKindedRejections_Link
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindLink) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

		type {{ .Type | mungeTypeNodeIdent }} struct{ x ipld.Link }

	`, w, gk)
}

func (gk generateKindLink) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindLink) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Link
		}
	`, w, gk)
}

func (gk generateKindLink) EmitNodeMethodAsLink(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) AsLink() (ipld.Link, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindLink) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}

func (gk generateKindLink) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
