package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &mapReprMapGenerator{}

func NewMapReprMapGenerator(pkgName string, typ schema.TypeMap, adjCfg *AdjunctCfg) TypeGenerator {
	return mapReprMapGenerator{
		mapGenerator{
			adjCfg,
			mixins.MapTraits{
				pkgName,
				string(typ.Name()),
				adjCfg.TypeSymbol(typ),
			},
			pkgName,
			typ,
		},
	}
}

type mapReprMapGenerator struct {
	mapGenerator
}

func (g mapReprMapGenerator) GetRepresentationNodeGen() NodeGenerator {
	return mapReprMapReprGenerator{
		g.AdjCfg,
		mixins.MapTraits{
			g.PkgName,
			string(g.Type.Name()) + ".Repr",
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Repr",
		},
		g.PkgName,
		g.Type,
	}
}

type mapReprMapReprGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapTraits
	PkgName string
	Type    schema.TypeMap
}

// FIXME: the representation for maps is NOT natural if the map has a complex key!

func (g mapReprMapReprGenerator) EmitNodeType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr = _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}
func (mapReprMapReprGenerator) EmitNodeMethodReprKind(io.Writer)      {}
func (mapReprMapReprGenerator) EmitNodeMethodLookupString(io.Writer)  {}
func (mapReprMapReprGenerator) EmitNodeMethodLookup(io.Writer)        {}
func (mapReprMapReprGenerator) EmitNodeMethodLookupIndex(io.Writer)   {}
func (mapReprMapReprGenerator) EmitNodeMethodLookupSegment(io.Writer) {}
func (mapReprMapReprGenerator) EmitNodeMethodMapIterator(io.Writer)   {}
func (mapReprMapReprGenerator) EmitNodeMethodListIterator(io.Writer)  {}
func (mapReprMapReprGenerator) EmitNodeMethodLength(io.Writer)        {}
func (mapReprMapReprGenerator) EmitNodeMethodIsUndefined(io.Writer)   {}
func (mapReprMapReprGenerator) EmitNodeMethodIsNull(io.Writer)        {}
func (mapReprMapReprGenerator) EmitNodeMethodAsBool(io.Writer)        {}
func (mapReprMapReprGenerator) EmitNodeMethodAsInt(io.Writer)         {}
func (mapReprMapReprGenerator) EmitNodeMethodAsFloat(io.Writer)       {}
func (mapReprMapReprGenerator) EmitNodeMethodAsString(io.Writer)      {}
func (mapReprMapReprGenerator) EmitNodeMethodAsBytes(io.Writer)       {}
func (mapReprMapReprGenerator) EmitNodeMethodAsLink(io.Writer)        {}
func (mapReprMapReprGenerator) EmitNodeMethodStyle(io.Writer)         {}
func (g mapReprMapReprGenerator) EmitNodeStyleType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprStyle = _{{ .Type | TypeSymbol }}__Style
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return nil // TODO
}
