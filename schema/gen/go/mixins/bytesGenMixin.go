package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type BytesTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (BytesTraits) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bytes
}
func (g BytesTraits) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Bytes
		}
	`, w, g)
}
func (g BytesTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodLookupByString(w)
}
func (g BytesTraits) EmitNodeMethodLookupByNode(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodLookupByNode(w)
}
func (g BytesTraits) EmitNodeMethodLookupByIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodLookupByIndex(w)
}
func (g BytesTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodLookupBySegment(w)
}
func (g BytesTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodMapIterator(w)
}
func (g BytesTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodListIterator(w)
}
func (g BytesTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodLength(w)
}
func (g BytesTraits) EmitNodeMethodIsUndefined(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodIsUndefined(w)
}
func (g BytesTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodIsNull(w)
}
func (g BytesTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodAsBool(w)
}
func (g BytesTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodAsInt(w)
}
func (g BytesTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodAsFloat(w)
}
func (g BytesTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodAsString(w)
}
func (g BytesTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_Bytes}.emitNodeMethodAsLink(w)
}

type BytesAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (BytesAssemblerTraits) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bytes
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodBeginMap(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodBeginList(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignNull(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignBool(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignInt(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignString(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodAssignLink(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodStyle(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_Bytes}.emitNodeAssemblerMethodStyle(w)
}
