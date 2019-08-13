package gengo

import (
	"io"
)

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type.Name }}{}
		var _ typed.Node = typed.Node(nil) // TODO

		type {{ .Type.Name }} struct{ x string }

	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}

// FUTURE: consider breaking the nodebuilder methods down like the node methods are; a lot of the "nope" variants could be reused.
func (gk generateKindString) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type.Name }}__NodeBuilder{}
		}

		type {{ .Type.Name }}__NodeBuilder struct{}

		func (nb {{ .Type.Name }}__NodeBuilder) CreateMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) AmendMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) AmendList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateNull() (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateBool(v bool) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateInt(v int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateFloat(v float64) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateString(v string) (ipld.Node, error) {
			return {{ .Type.Name }}{v}, nil
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateBytes(v []byte) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Type.Name }}__NodeBuilder) CreateLink(v ipld.Link) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Type.Name }}", MethodName: "NodeBuilder.CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}
