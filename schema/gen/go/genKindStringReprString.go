package gengo

import (
	"io"
)

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Name }}{}

		type {{ .Name }} struct{ x string }

	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}

// FUTURE: consider breaking the nodebuilder methods down like the node methods are; a lot of the "nope" variants could be reused.
func (gk generateKindString) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Name }}__NodeBuilder{}
		}

		type {{ .Name }}__NodeBuilder struct{}

		func (nb {{ .Name }}__NodeBuilder) CreateMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) AmendMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{MethodName: "AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) AmendList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{MethodName: "AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateNull() (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateBool(v bool) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateInt(v int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateFloat(v float64) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateString(v string) (ipld.Node, error) {
			return {{ .Name }}{v}, nil
		}
		func (nb {{ .Name }}__NodeBuilder) CreateBytes(v []byte) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
		}
		func (nb {{ .Name }}__NodeBuilder) CreateLink(v ipld.Link) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}
