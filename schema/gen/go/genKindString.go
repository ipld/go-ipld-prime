package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

type generateKindString struct {
	Name schema.TypeName
	Type schema.Type
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

// FUTURE: quite a few of these "nope" methods can be reused widely.

func (gk generateKindString) EmitNodeMethodTraverseField(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) TraverseField(key string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "TraverseField", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodTraverseIndex(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) TraverseIndex(idx int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{MethodName: "TraverseIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) MapIterator() ipld.MapIterator {
			return mapIteratorReject{ipld.ErrWrongKind{MethodName: "MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodListIterator(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) ListIterator() ipld.ListIterator {
			return listIteratorReject{ipld.ErrWrongKind{MethodName: "ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}}
		}
	`, w, gk) // REVIEW: maybe that rejection thunk should be in main package?  don't really want to flash it at folks though.  very impl detail.
}

func (gk generateKindString) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) Length() int {
			return -1
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodIsNull(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) IsNull() bool {
			return false
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsBool(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) AsBool() (bool, error) {
			return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) AsInt() (int, error) {
			return 0, ipld.ErrWrongKind{MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsFloat(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) AsFloat() (float64, error) {
			return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (x {{ .Name }}) AsString() (string, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsBytes(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsLink(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) AsLink() (ipld.Link, error) {
			return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Name }}__NodeBuilder{}
		}
	`, w, gk)
}
