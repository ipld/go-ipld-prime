package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func (gk generateKindInt) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindInt{
		gk.Type,
		genKindedNbRejections_Int{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

func (gk generateKindInt) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

type generateNbKindInt struct {
	Type schema.TypeInt
	genKindedNbRejections_Int
}

func (gk generateNbKindInt) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodebuilderIdent }} struct{}
	`, w, gk)
}

func (gk generateNbKindInt) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateNbKindInt) EmitNodebuilderMethodCreateInt(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateInt(v int) (ipld.Node, error) {
			return {{ .Type | mungeTypeNodeIdent }}{v}, nil
		}
	`, w, gk)
}
