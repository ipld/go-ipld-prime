package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func (gk generateKindString) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindString{
		gk.Type,
		genKindedNbRejections_String{gk.Type},
	}
}

func (gk generateKindString) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type.Name }}__NodeBuilder{}
		}
	`, w, gk)
}

type generateNbKindString struct {
	Type schema.TypeString
	genKindedNbRejections_String
}

func (gk generateNbKindString) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type.Name }}__NodeBuilder struct{}
	`, w, gk)
}

func (gk generateNbKindString) EmitNodebuilderMethodCreateString(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type.Name }}__NodeBuilder) CreateString(v string) (ipld.Node, error) {
			return {{ .Type.Name }}{v}, nil
		}
	`, w, gk)
}
