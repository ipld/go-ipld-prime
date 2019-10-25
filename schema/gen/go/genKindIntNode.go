package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

// --- type-semantics node interface satisfaction --->

func (gk generateKindInt) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

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

// --- type-semantics nodebuilder --->

func (gk generateKindInt) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateKindInt) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindInt{
		gk.Type,
		genKindedNbRejections_Int{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
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

// --- entrypoints to representation --->

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
