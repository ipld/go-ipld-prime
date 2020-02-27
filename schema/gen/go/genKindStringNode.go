package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

// --- type-semantics node interface satisfaction --->

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ schema.TypedNode = {{ .Type | mungeTypeNodeIdent }}{}

	`, w, gk)
}

func (gk generateKindString) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) AsString() (string, error) {
			return x.x, nil
		}
	`, w, gk)
}

// --- type-semantics nodebuilder --->

func (gk generateKindString) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateKindString) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindString{
		gk.Type,
		genKindedNbRejections_String{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

type generateNbKindString struct {
	Type schema.TypeString
	genKindedNbRejections_String
}

func (gk generateNbKindString) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodebuilderIdent }} struct{}
	`, w, gk)
}

func (gk generateNbKindString) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateNbKindString) EmitNodebuilderMethodCreateString(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateString(v string) (ipld.Node, error) {
			return {{ .Type | mungeTypeNodeIdent }}{v}, nil
		}
	`, w, gk)
}

// --- entrypoints to representation --->

func (gk generateKindString) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}
func (gk generateKindString) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
