package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &listReprListGenerator{}

func NewListReprListGenerator(pkgName string, typ schema.TypeList, adjCfg *AdjunctCfg) TypeGenerator {
	return listReprListGenerator{
		listGenerator{
			adjCfg,
			mixins.ListTraits{
				pkgName,
				string(typ.Name()),
				adjCfg.TypeSymbol(typ),
			},
			pkgName,
			typ,
		},
	}
}

type listReprListGenerator struct {
	listGenerator
}

func (g listReprListGenerator) GetRepresentationNodeGen() NodeGenerator {
	return listReprListReprGenerator{
		g.AdjCfg,
		mixins.ListTraits{
			g.PkgName,
			string(g.Type.Name()) + ".Repr",
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Repr",
		},
		g.PkgName,
		g.Type,
	}
}

type listReprListReprGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.ListTraits
	PkgName string
	Type    schema.TypeList
}

func (g listReprListReprGenerator) EmitNodeType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr = _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}
func (g listReprListReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}
func (listReprListReprGenerator) EmitNodeMethodReprKind(io.Writer)      {}
func (listReprListReprGenerator) EmitNodeMethodLookupString(io.Writer)  {}
func (listReprListReprGenerator) EmitNodeMethodLookup(io.Writer)        {}
func (listReprListReprGenerator) EmitNodeMethodLookupIndex(io.Writer)   {}
func (listReprListReprGenerator) EmitNodeMethodLookupSegment(io.Writer) {}
func (listReprListReprGenerator) EmitNodeMethodMapIterator(io.Writer)   {}
func (listReprListReprGenerator) EmitNodeMethodListIterator(io.Writer)  {}
func (listReprListReprGenerator) EmitNodeMethodLength(io.Writer)        {}
func (listReprListReprGenerator) EmitNodeMethodIsUndefined(io.Writer)   {}
func (listReprListReprGenerator) EmitNodeMethodIsNull(io.Writer)        {}
func (listReprListReprGenerator) EmitNodeMethodAsBool(io.Writer)        {}
func (listReprListReprGenerator) EmitNodeMethodAsInt(io.Writer)         {}
func (listReprListReprGenerator) EmitNodeMethodAsFloat(io.Writer)       {}
func (listReprListReprGenerator) EmitNodeMethodAsString(io.Writer)      {}
func (listReprListReprGenerator) EmitNodeMethodAsBytes(io.Writer)       {}
func (listReprListReprGenerator) EmitNodeMethodAsLink(io.Writer)        {}
func (listReprListReprGenerator) EmitNodeMethodStyle(io.Writer)         {}
func (g listReprListReprGenerator) EmitNodeStyleType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprStyle = _{{ .Type | TypeSymbol }}__Style
	`, w, g.AdjCfg, g)
}
func (g listReprListReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return nil // TODO
}
