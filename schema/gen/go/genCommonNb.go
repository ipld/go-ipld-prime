package gengo

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type genKindedNbRejections struct {
	TypeIdent string // the identifier in code (sometimes is munged internals like "_Thing__Repr" corresponding to no publicly admitted schema.Type.Name).
	TypeProse string // as will be printed in messages (e.g. can be goosed up a bit, like "Thing.Repr" instead of "_Thing__Repr").
	Kind      ipld.ReprKind
}

func (d genKindedNbRejections) emitNodebuilderMethodCreateMap(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodAmendMap(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AmendMap() (ipld.MapBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateList(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodAmendList(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) AmendList() (ipld.ListBuilder, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateNull(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateNull() (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateBool(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateBool(bool) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateInt(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateInt(int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateFloat(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateFloat(float64) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateString(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateString(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateBytes(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateBytes([]byte) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}
func (d genKindedNbRejections) emitNodebuilderMethodCreateLink(w io.Writer) {
	doTemplate(`
		func ({{ .TypeIdent }}) CreateLink(ipld.Link) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .TypeProse }}", MethodName: "CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: {{ .Kind | ReprKindConst }}}
		}
	`, w, d)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_String struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateMap(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodAmendMap(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateList(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodAmendList(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateInt(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateBytes(w)
}
func (gk genKindedNbRejections_String) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_String}.emitNodebuilderMethodCreateLink(w)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_Map struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateList(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodAmendList(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateInt(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateString(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateBytes(w)
}
func (gk genKindedNbRejections_Map) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Map}.emitNodebuilderMethodCreateLink(w)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_Int struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateMap(w)
}

func (gk genKindedNbRejections_Int) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodAmendMap(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateList(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodAmendList(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateString(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateBytes(w)
}
func (gk genKindedNbRejections_Int) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Int}.emitNodebuilderMethodCreateLink(w)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_Bytes struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateMap(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodAmendMap(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateList(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodAmendList(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateInt(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateString(w)
}
func (gk genKindedNbRejections_Bytes) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Bytes}.emitNodebuilderMethodCreateLink(w)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_Link struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateMap(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodAmendMap(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateList(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodAmendList(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodAmendList(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateInt(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateBytes(w)
}
func (gk genKindedNbRejections_Link) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_Link}.emitNodebuilderMethodCreateString(w)
}

// Embeddable to do all the "nope" methods at once.
type genKindedNbRejections_List struct {
	TypeIdent string // see doc in generateKindedRejections
	TypeProse string // see doc in generateKindedRejections
}

func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateMap(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodAmendMap(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodAmendMap(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateNull(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateNull(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateBool(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateBool(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateInt(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateInt(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateFloat(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateFloat(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateString(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateString(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateBytes(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateBytes(w)
}
func (gk genKindedNbRejections_List) EmitNodebuilderMethodCreateLink(w io.Writer) {
	genKindedNbRejections{gk.TypeIdent, gk.TypeProse, ipld.ReprKind_List}.emitNodebuilderMethodCreateLink(w)
}
