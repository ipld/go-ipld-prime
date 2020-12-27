package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type IntTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (IntTraits) Kind() ipld.Kind {
	return ipld.Kind_Int
}
func (g IntTraits) EmitNodeMethodKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) Kind() ipld.Kind {
			return ipld.Kind_Int
		}
	`, w, g)
}
func (g IntTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodLookupByString(w)
}
func (g IntTraits) EmitNodeMethodLookupByNode(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodLookupByNode(w)
}
func (g IntTraits) EmitNodeMethodLookupByIndex(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodLookupByIndex(w)
}
func (g IntTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodLookupBySegment(w)
}
func (g IntTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodMapIterator(w)
}
func (g IntTraits) EmitNodeMethodListIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodListIterator(w)
}
func (g IntTraits) EmitNodeMethodLength(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodLength(w)
}
func (g IntTraits) EmitNodeMethodIsAbsent(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodIsAbsent(w)
}
func (g IntTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodIsNull(w)
}
func (g IntTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodAsBool(w)
}
func (g IntTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodAsFloat(w)
}
func (g IntTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodAsString(w)
}
func (g IntTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodAsBytes(w)
}
func (g IntTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_Int}.emitNodeMethodAsLink(w)
}

type IntAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (IntAssemblerTraits) Kind() ipld.Kind {
	return ipld.Kind_Int
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodBeginMap(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodBeginList(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignNull(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignBool(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignString(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodAssignLink(w)
}
func (g IntAssemblerTraits) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_Int}.emitNodeAssemblerMethodPrototype(w)
}
