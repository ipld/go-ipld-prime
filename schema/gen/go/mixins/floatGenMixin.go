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

func (FloatTraits) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Float
}
func (g FloatTraits) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Float
		}
	`, w, g)
}
func (g FloatTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodLookupByString(w)
}
func (g FloatTraits) EmitNodeMethodLookupByNode(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodLookupByNode(w)
}
func (g FloatTraits) EmitNodeMethodLookupByIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodLookupByIndex(w)
}
func (g FloatTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodLookupBySegment(w)
}
func (g FloatTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodMapIterator(w)
}
func (g FloatTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodListIterator(w)
}
func (g FloatTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodLength(w)
}
func (g FloatTraits) EmitNodeMethodIsUndefined(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodIsUndefined(w)
}
func (g FloatTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodIsNull(w)
}
func (g FloatTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodAsBool(w)
}
func (g FloatTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodAsInt(w)
}
func (g FloatTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodAsString(w)
}
func (g FloatTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodAsBytes(w)
}
func (g FloatTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Float}.emitNodeMethodAsLink(w)
}

type FloatAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (FloatAssemblerTraits) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Float
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodBeginMap(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodBeginList(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignNull(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignBool(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignInt(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignString(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodAssignLink(w)
}
func (g FloatAssemblerTraits) EmitNodeAssemblerMethodStyle(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Float}.emitNodeAssemblerMethodStyle(w)
}
