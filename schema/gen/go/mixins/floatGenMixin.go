package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type FloatTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (FloatTraits) Kind() ipld.Kind {
	return ipld.Kind_Float
}
func (g FloatTraits) EmitNodeMethodKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) Kind() ipld.Kind {
			return ipld.Kind_Float
		}
	`, w, g)
}
func (g FloatTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodLookupByString(w)
}
func (g FloatTraits) EmitNodeMethodLookupByNode(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodLookupByNode(w)
}
func (g FloatTraits) EmitNodeMethodLookupByIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodLookupByIndex(w)
}
func (g FloatTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodLookupBySegment(w)
}
func (g FloatTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodMapIterator(w)
}
func (g FloatTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodListIterator(w)
}
func (g FloatTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodLength(w)
}
func (g FloatTraits) EmitNodeMethodIsAbsent(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodIsAbsent(w)
}
func (g FloatTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodIsNull(w)
}
func (g FloatTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodAsBool(w)
}
func (g FloatTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodAsInt(w)
}
func (g FloatTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodAsString(w)
}
func (g FloatTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodAsBytes(w)
}
func (g FloatTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Float}.emitNodeMethodAsLink(w)
}

type FloatAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (FloatAssemblerTraits) Kind() ipld.Kind {
	return ipld.Kind_Float
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodBeginMap(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodBeginList(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignNull(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignBool(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignInt(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignString(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodAssignLink(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Float}.emitNodeAssemblerMethodPrototype(w)
}
