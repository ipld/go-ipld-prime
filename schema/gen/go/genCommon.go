package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

type generateKindedRejections struct{}

func (generateKindedRejections) emitNodeMethodLookupString(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) LookupString(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodLookup(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) Lookup(ipld.Node) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodLookupIndex(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) LookupIndex(idx int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodMapIterator(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) MapIterator() ipld.MapIterator {
			return mapIteratorReject{ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodListIterator(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) ListIterator() ipld.ListIterator {
			return listIteratorReject{ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}}
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

func (generateKindedRejections) emitNodeMethodIsUndefined(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) IsUndefined() bool {
			return false
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
			return false, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsInt(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsInt() (int, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsFloat(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsFloat() (float64, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsString(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsString() (string, error) {
			return "", ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsBytes(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

func (generateKindedRejections) emitNodeMethodAsLink(w io.Writer, t schema.Type) {
	doTemplate(`
		func ({{ .Name }}) AsLink() (ipld.Link, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .Name }}", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: {{ .Kind.ActsLike | ReprKindConst }}}
		}
	`, w, t)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_String struct {
	Type schema.Type // used so we can generate error messages with the type name.
}

func (gk generateKindedRejections_String) EmitNodeMethodLookupString(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodLookupString(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodLookup(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodLookup(w, gk.Type)
}
func (gk generateKindedRejections_String) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodLookupIndex(w, gk.Type)
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
func (gk generateKindedRejections_String) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodIsUndefined(w, gk.Type)
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

// Embeddable to do all the "nope" methods at once.
//
// Used for anything that "acts like" map (so, also struct).
type generateKindedRejections_Map struct {
	Type schema.Type // used so we can generate error messages with the type name.
}

func (gk generateKindedRejections_Map) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodLookupIndex(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodListIterator(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodIsUndefined(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodIsNull(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsBool(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsInt(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsFloat(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsString(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsString(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsBytes(w, gk.Type)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{}.emitNodeMethodAsLink(w, gk.Type)
}
