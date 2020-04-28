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

/*

	- lookups are essentially the same except you add a '.Representation' on the end.
		- *maybe* frumped by dealing with nullable.  hopefully we can write that a textually brief way and let the compiler inline/optimize it.
		- "cast" yourself back to the type-level node and call that method explicitly, then do the suffix.

	- iterators are in a similar boat to lookups.
		- still do also definitely need to generate a new named iterator type, though, which... probably bombs a lot.
			- actually, can we just do one of those `type foo bar` things again?  save a bit of gsloc that way?  probably.

	- infuriatingly, i don't see a way to skimp on the darn stub methods.
		- to be fair, they're supposed to return a moderately different error message anyway -- ".Repr" in the type name text.
		- at least we have reduced sloc here in the generator, since the mixin methods apply again.

	- assembler has a different child assembler type to embed (poosibly with radically different logical behavior), so that's deffo a new type.
		- but a LOT of the logic is subsequently the same right up until that hand-off to the child assembler.
			- most of the major 'state' and 'm' transition logic is the same (Finish especially)
			- slice growth and maybe-allocs are the same, because the child node type is the same, even though the repr assembler is divergent
			- all the 'tidy' logic falls in with the above
			- we might even be able to extract all of these, if we can make them regard just '*state' and '*m' parameters.
				- i'm not sure if this would have negative effects on binary size or optimizations though.

	- AssignNode legitimately differs (but only for the bottom third) because... wait, no, *textually*, it's identical.
		- it calls out to AssembleValue, which will differ in that it calls the Repr assembler for child, but that's it.

	- BeginList is also textually identical except for the type it has to be attached to >:I

*/
