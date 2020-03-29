package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type MapTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (g MapTraits) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Map
		}
	`, w, g)
}
func (g MapTraits) EmitNodeMethodLookupIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodLookupIndex(w)
}
func (g MapTraits) EmitNodeMethodLookupSegment(w io.Writer) {
	doTemplate(`
		func (n {{ .TypeSymbol }}) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
			return n.LookupString(seg.String())
		}
	`, w, g)
}
func (g MapTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodListIterator(w)
}
func (g MapTraits) EmitNodeMethodIsUndefined(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodIsUndefined(w)
}
func (g MapTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodIsNull(w)
}
func (g MapTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsBool(w)
}
func (g MapTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsInt(w)
}
func (g MapTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsFloat(w)
}
func (g MapTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsString(w)
}
func (g MapTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsBytes(w)
}
func (g MapTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Map}.emitNodeMethodAsLink(w)
}
