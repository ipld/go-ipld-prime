package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

// --- type-semantics node interface satisfaction --->

func (gk generateKindLink) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

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

// --- type-semantics nodebuilder --->

func (gk generateKindLink) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateKindLink) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindLink{
		gk.Type,
		genKindedNbRejections_Link{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

type generateNbKindLink struct {
	Type schema.TypeLink
	genKindedNbRejections_Link
}

func (gk generateNbKindLink) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodebuilderIdent }} struct{}

	`, w, gk)
}

func (gk generateNbKindLink) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateNbKindLink) EmitNodebuilderMethodCreateLink(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateLink(v ipld.Link) (ipld.Node, error) {
			return {{ .Type | mungeTypeNodeIdent }}{v}, nil
		}
	`, w, gk)
}

// --- entrypoints to representation --->

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
