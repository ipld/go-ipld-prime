package gengo

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type generateKindedRejections struct {
	TypeIdent string // the identifier in code (sometimes is munged internals like "_Thing__Repr" corresponding to no publicly admitted schema.Type.Name).
	TypeProse string // as will be printed in messages (e.g. can be goosed up a bit, like "Thing.Repr" instead of "_Thing__Repr").
	Kind      ipld.ReprKind
}

func (d generateKindedRejections) emitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) LookupString(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) Lookup(ipld.Node) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodLookupIndex(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) LookupIndex(idx int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodLookupSegment(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) MapIterator() ipld.MapIterator {
			return mapIteratorReject{ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind | ReprKindConst }}}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodListIterator(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) ListIterator() ipld.ListIterator {
			return listIteratorReject{ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind | ReprKindConst }}}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) Length() int {
			return -1
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodIsUndefined(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) IsUndefined() bool {
			return false
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodIsNull(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) IsNull() bool {
			return false
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsBool(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsBool() (bool, error) {
			return false, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsInt() (int, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsFloat(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsFloat() (float64, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsString() (string, error) {
			return "", ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsBytes(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

func (d generateKindedRejections) emitNodeMethodAsLink(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AsLink() (ipld.Link, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_String struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk generateKindedRejections_String) EmitNodeMethodLookupString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodLookupString(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodLookup(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodLookup(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodLookupIndex(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodLookupSegment(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodLookupSegment(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodMapIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodMapIterator(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodListIterator(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodLength(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodLength(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodIsUndefined(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodIsNull(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodAsBool(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodAsInt(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodAsFloat(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodAsBytes(w)
}
func (gk generateKindedRejections_String) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodeMethodAsLink(w)
}

// Embeddable to do all the "nope" methods at once.
//
// Used for anything that "acts like" map (so, also struct).
type generateKindedRejections_Map struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk generateKindedRejections_Map) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodLookupIndex(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodLookupSegment(w io.Writer) {
	doTemplate(`
		func (n {{ .TypeIdent }}) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
			return n.LookupString(seg.String())
		}
	`, w, gk)
}
func (gk generateKindedRejections_Map) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodListIterator(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodIsUndefined(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodIsNull(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsBool(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsInt(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsFloat(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsString(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsBytes(w)
}
func (gk generateKindedRejections_Map) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodeMethodAsLink(w)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_Int struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk generateKindedRejections_Int) EmitNodeMethodLookupString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodLookupString(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodLookup(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodLookup(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodLookupIndex(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodLookupSegment(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodLookupSegment(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodMapIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodMapIterator(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodListIterator(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodLength(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodLength(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodIsUndefined(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodIsNull(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodAsBool(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodAsString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodAsString(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodAsFloat(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodAsBytes(w)
}
func (gk generateKindedRejections_Int) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodeMethodAsLink(w)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_Bytes struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk generateKindedRejections_Bytes) EmitNodeMethodLookupString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodLookupString(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodLookup(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodLookup(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodLookupIndex(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodLookupSegment(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodLookupSegment(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodMapIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodMapIterator(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodListIterator(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodLength(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodLength(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodIsUndefined(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodIsNull(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodAsBool(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodAsInt(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodAsFloat(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodAsString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodAsString(w)
}
func (gk generateKindedRejections_Bytes) EmitNodeMethodAsLink(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodeMethodAsLink(w)
}

// Embeddable to do all the "nope" methods at once.
type generateKindedRejections_Link struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk generateKindedRejections_Link) EmitNodeMethodLookupString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodLookupString(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodLookup(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodLookup(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodLookupIndex(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodLookupIndex(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodLookupSegment(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodLookupSegment(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodMapIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodMapIterator(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodListIterator(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodListIterator(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodLength(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodLength(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodIsUndefined(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodIsUndefined(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodIsNull(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodIsNull(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodAsBool(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodAsBool(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodAsInt(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodAsInt(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodAsFloat(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodAsFloat(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodAsBytes(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodAsBytes(w)
}
func (gk generateKindedRejections_Link) EmitNodeMethodAsString(w io.Writer) {
	generateKindedRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodeMethodAsString(w)
}
