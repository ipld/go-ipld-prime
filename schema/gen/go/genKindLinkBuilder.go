package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func (gk generateKindLink) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindLink{
		gk.Type,
		genKindedNbRejections_Link{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

func (gk generateKindLink) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
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
