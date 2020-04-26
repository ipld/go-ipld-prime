package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type ListTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (g ListTraits) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_List
		}
	`, w, g)
}
func (g ListTraits) EmitNodeMethodLookupString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodLookupString(w)
}
func (g ListTraits) EmitNodeMethodLookupSegment(w io.Writer) {
	doTemplate(`
		func (n {{ .TypeSymbol }}) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
			i, err := seg.Index()
			if err != nil {
				return nil, ipld.ErrInvalidSegmentForList{TypeName: "{{ .PkgName }}.{{ .TypeName }}", TroubleSegment: seg, Reason: err}
			}
			return n.LookupIndex(i)
		}
	`, w, g)
}
func (g ListTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodMapIterator(w)
}
func (g ListTraits) EmitNodeMethodIsUndefined(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodIsUndefined(w)
}
func (g ListTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodIsNull(w)
}
func (g ListTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsBool(w)
}
func (g ListTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsInt(w)
}
func (g ListTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsFloat(w)
}
func (g ListTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsString(w)
}
func (g ListTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsBytes(w)
}
func (g ListTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.ReprKind_List}.emitNodeMethodAsLink(w)
}

type ListAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (ListAssemblerTraits) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_List
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodBeginMap(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignNull(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignBool(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignInt(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignString(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodAssignLink(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodStyle(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.ReprKind_List}.emitNodeAssemblerMethodStyle(w)
}
