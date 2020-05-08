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

func (listReprListReprGenerator) IsRepr() bool { return true } // hint used in some generalized templates.

func (g listReprListReprGenerator) EmitNodeType(w io.Writer) {
	// Even though this is a "natural" representation... we need a new type here,
	//  because lists are recursive, and so all our functions that access
	//   children need to remember to return the representation node of those child values.
	// It's still structurally the same, though (and we'll be able to cast in the methodset pattern).
	// Error-thunking methods also have a different string in their error, so those are unique even if they don't seem particularly interesting.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeMethodLookup(w io.Writer) {
	// Null is also already a branch in the method we're calling; hopefully the compiler inlines and sees this and DTRT.
	// REVIEW: these unchecked casts are definitely safe at compile time, but I'm not sure if the compiler considers that provable,
	//  so we should investigate if there's any runtime checks injected here that waste time.  If so: write this with more gsloc to avoid :(
	doTemplate(`
		func (nr *_{{ .Type | TypeSymbol }}__Repr) Lookup(k ipld.Node) (ipld.Node, error) {
			v, err := ({{ .Type | TypeSymbol }})(nr).Lookup(k)
			if err != nil || v == ipld.Null {
				return v, err
			}
			return v.({{ .Type.ValueType | TypeSymbol}}).Representation(), nil
		}
	`, w, g.AdjCfg, g)

}

func (g listReprListReprGenerator) EmitNodeMethodLookupIndex(w io.Writer) {
	doTemplate(`
		func (nr *_{{ .Type | TypeSymbol }}__Repr) LookupIndex(idx int) (ipld.Node, error) {
			v, err := ({{ .Type | TypeSymbol }})(nr).LookupIndex(idx)
			if err != nil || v == ipld.Null {
				return v, err
			}
			return v.({{ .Type.ValueType | TypeSymbol}}).Representation(), nil
		}
	`, w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeMethodListIterator(w io.Writer) {
	// FUTURE: trying to get this to share the preallocated memory if we get iterators wedged into their node slab will be ... fun.
	doTemplate(`
		func (nr *_{{ .Type | TypeSymbol }}__Repr) ListIterator() ipld.ListIterator {
			return &_{{ .Type | TypeSymbol }}__ReprListItr{({{ .Type | TypeSymbol }})(nr), 0}
		}

		type _{{ .Type | TypeSymbol }}__ReprListItr _{{ .Type | TypeSymbol }}__ListItr

		func (itr *_{{ .Type | TypeSymbol }}__ReprListItr) Next() (idx int, v ipld.Node, err error) {
			idx, v, err = (*_{{ .Type | TypeSymbol }}__ListItr)(itr).Next()
			if err != nil || v == ipld.Null {
				return
			}
			return idx, v.({{ .Type.ValueType | TypeSymbol}}).Representation(), nil
		}
		func (itr *_{{ .Type | TypeSymbol }}__ReprListItr) Done() bool {
			return (*_{{ .Type | TypeSymbol }}__ListItr)(itr).Done()
		}

	`, w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func (rn *_{{ .Type | TypeSymbol }}__Repr) Length() int {
			return len(rn.x)
		}
	`, w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeMethodStyle(w io.Writer) {
	emitNodeMethodStyle_typical(w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) EmitNodeStyleType(w io.Writer) {
	// FIXME this alias is a lie, but keep it around for one more sec so we can do an incremental commit until we finish fixing it
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprStyle = _{{ .Type | TypeSymbol }}__Style
	`, w, g.AdjCfg, g)
	//emitNodeStyleType_typical(w, g.AdjCfg, g)
}

func (g listReprListReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return nil // TODO
}

/*

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
