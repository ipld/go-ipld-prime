package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type StringTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (g StringTraits) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, g)
}
func (g StringTraits) EmitNodeMethodLookupString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodLookupString(w)
}
func (g StringTraits) EmitNodeMethodLookup(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodLookup(w)
}
func (g StringTraits) EmitNodeMethodLookupIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodLookupIndex(w)
}
func (g StringTraits) EmitNodeMethodLookupSegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodLookupSegment(w)
}
func (g StringTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodMapIterator(w)
}
func (g StringTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodListIterator(w)
}
func (g StringTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodLength(w)
}
func (g StringTraits) EmitNodeMethodIsUndefined(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodIsUndefined(w)
}
func (g StringTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodIsNull(w)
}
func (g StringTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodAsBool(w)
}
func (g StringTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodAsInt(w)
}
func (g StringTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodAsFloat(w)
}
func (g StringTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodAsBytes(w)
}
func (g StringTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_String}.emitNodeMethodAsLink(w)
}
