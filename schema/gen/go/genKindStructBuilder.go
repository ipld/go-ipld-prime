package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func (gk generateKindStruct) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindStruct{
		gk.Type,
		genKindedNbRejections_Struct{gk.Type},
	}
}

func (gk generateKindStruct) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type.Name }}__NodeBuilder{}
		}
	`, w, gk)
}

type generateNbKindStruct struct {
	Type schema.TypeStruct
	genKindedNbRejections_Struct
}

func (gk generateNbKindStruct) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type.Name }}__NodeBuilder struct{}
	`, w, gk)
}

func (gk generateNbKindStruct) EmitNodebuilderMethodCreateMap(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type.Name }}__NodeBuilder) CreateMap() (ipld.MapBuilder, error) {
			return &{{ .Type.Name }}__MapBuilder{&{{ .Type.Name }}{}}, nil
		}

		type {{ .Type.Name }}__MapBuilder struct{
			v *{{ .Type.Name }}
			// TODO will probably have more state for validating
		}

		func (mb *{{ .Type.Name }}__MapBuilder) Insert(k, v ipld.Node) error {
			panic("TODO now")
		}

		func (mb *{{ .Type.Name }}__MapBuilder) Delete(k ipld.Node) error {
			panic("TODO later")
		}

		func (mb *{{ .Type.Name }}__MapBuilder) Build() (ipld.Node, error) {
			v := mb.v
			mb = nil
			return v, nil
		}
	`, w, gk)
}
func (gk generateNbKindStruct) EmitNodebuilderMethodAmendMap(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type.Name }}__NodeBuilder) AmendMap() (ipld.MapBuilder, error) {
			panic("TODO later")
		}
	`, w, gk)
}
