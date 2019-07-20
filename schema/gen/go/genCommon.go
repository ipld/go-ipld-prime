package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

type generateKindedRejections struct{}

func (generateKindedRejections) emitNodeMethodTraverseField(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) TraverseField(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "{{ .Name }}.TraverseField", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodTraverseIndex(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) TraverseIndex(idx int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "{{ .Name }}.TraverseIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodMapIterator(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) MapIterator() ipld.MapIterator {
			return mapIteratorReject{ipld.ErrWrongKind{MethodName: "{{ .Name }}.MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodListIterator(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) ListIterator() ipld.ListIterator {
			return listIteratorReject{ipld.ErrWrongKind{MethodName: "{{ .Name }}.ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodLength(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) Length() int {
			return -1
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodIsNull(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) IsNull() bool {
			return false
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsBool(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsBool() (bool, error) {
			return false, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsInt(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsInt() (int, error) {
			return 0, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsFloat(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsFloat() (float64, error) {
			return 0, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsString(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsBytes", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsBytes(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsLink(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsLink() (ipld.Link, error) {
			return nil, ipld.ErrWrongKind{MethodName: "{{ .Name }}.AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_String struct {
	Type schema.Type // used so we can generate error messages with the type name.
}

func (gk generateKindedRejections_String) EmitNodeMethodTraverseField(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodTraverseField(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodTraverseIndex(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodTraverseIndex(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodMapIterator(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodMapIterator(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodListIterator(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodLength(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodLength(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodIsNull(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsBool(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsInt(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsFloat(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsBytes(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsLink(w, gk.Type)
}
