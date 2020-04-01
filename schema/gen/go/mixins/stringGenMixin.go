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

type StringAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	TypeSymbol    string // the core symbol of the typed value (we assume the 'w' pointer will be a pointer to this type)
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (g StringAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodBeginMap(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodBeginList(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignNull(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignBool(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignInt(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	// REVIEW: unresolved if this is worth the effort or not; makes this trait need to know TypeSymbol;
	//  others are shaping up such that things that know TypeSymbol have been over the line for reusability anyway;
	//   should this one really be special?
	doTemplate(`
		func (na *{{ .AppliedPrefix }}Assembler) AssignString(v string) error {
			*na.w = {{ .TypeSymbol }}{v}
			return nil
		}
	`, w, g)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodAssignLink(w)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *{{ .AppliedPrefix }}Assembler) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				return na.AssignString(v2)
			}
		}
	`, w, g)
}
func (g StringAssemblerTraits) EmitNodeAssemblerMethodStyle(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_String}.emitNodeAssemblerMethodStyle(w)
}
