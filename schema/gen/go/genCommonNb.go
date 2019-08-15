package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

type genKindedNbRejections struct{}

func (genKindedNbRejections) emitNodebuilderMethodCreateMap(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodAmendMap(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) AmendMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateList(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodAmendList(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) AmendList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateNull(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateNull() (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateBool(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateBool(bool) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateInt(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateInt(int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateFloat(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateFloat(float64) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateString(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateString(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateBytes(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateBytes([]byte) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}
func (genKindedNbRejections) emitNodebuilderMethodCreateLink(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}__NodeBuilder) CreateLink(ipld.Link) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_String struct {
	Type schema.Type // used so we can generate error messages with the type name.
}

func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateMap(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodAmendMap(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateList(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodAmendList(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateNull(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateBool(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateInt(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateFloat(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateBytes(w, gk.Type)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateLink(w, gk.Type)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_Map struct {
	Type schema.Type // used so we can generate error messages with the type name.
}

func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateList(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodAmendList(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateNull(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateBool(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateInt(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateFloat(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateString(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateBytes(w, gk.Type)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{}.emitNodebuilderMethodCreateLink(w, gk.Type)
}
