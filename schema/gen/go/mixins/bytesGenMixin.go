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

func (BytesTraits) Kind() ipld.Kind {
	return ipld.Kind_Bytes
}
func (g BytesTraits) EmitNodeMethodKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) Kind() ipld.Kind {
			return ipld.Kind_Bytes
		}
	`, w, g)
}
func (g BytesTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodLookupByString(w)
}
func (g BytesTraits) EmitNodeMethodLookupByNode(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodLookupByNode(w)
}
func (g BytesTraits) EmitNodeMethodLookupByIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodLookupByIndex(w)
}
func (g BytesTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodLookupBySegment(w)
}
func (g BytesTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodMapIterator(w)
}
func (g BytesTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodListIterator(w)
}
func (g BytesTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodLength(w)
}
func (g BytesTraits) EmitNodeMethodIsAbsent(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodIsAbsent(w)
}
func (g BytesTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodIsNull(w)
}
func (g BytesTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodAsBool(w)
}
func (g BytesTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodAsInt(w)
}
func (g BytesTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodAsFloat(w)
}
func (g BytesTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodAsString(w)
}
func (g BytesTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Bytes}.emitNodeMethodAsLink(w)
}

type BytesAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (BytesAssemblerTraits) Kind() ipld.Kind {
	return ipld.Kind_Bytes
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodBeginMap(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodBeginList(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignNull(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignBool(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignInt(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignString(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodAssignLink(w)
}
func (g BytesAssemblerTraits) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Bytes}.emitNodeAssemblerMethodPrototype(w)
}
