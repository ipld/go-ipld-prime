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

func (mapReprMapReprGenerator) IsRepr() bool { return true } // hint used in some generalized templates.

func (g mapReprMapReprGenerator) EmitNodeType(w io.Writer) {
	// Even though this is a "natural" representation... we need a new type here,
	//  because maps are recursive, and so all our functions that access
	//   children need to remember to return the representation node of those child values.
	// It's still structurally the same, though (and we'll be able to cast in the methodset pattern).
	// Error-thunking methods also have a different string in their error, so those are unique even if they don't seem particularly interesting.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}

func (g mapReprMapReprGenerator) EmitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func (nr *_{{ .Type | TypeSymbol }}__Repr) LookupString(k string) (ipld.Node, error) {
			v, err := ({{ .Type | TypeSymbol }})(nr).LookupString(k)
			if err != nil || v == ipld.Null {
				return v, err
			}
			return v.({{ .Type.ValueType | TypeSymbol}}).Representation(), nil
		}
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeMethodLookup(w io.Writer) {
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
func (g mapReprMapReprGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	// FUTURE: trying to get this to share the preallocated memory if we get iterators wedged into their node slab will be ... fun.
	doTemplate(`
		func (nr *_{{ .Type | TypeSymbol }}__Repr) MapIterator() ipld.MapIterator {
			return &_{{ .Type | TypeSymbol }}__ReprMapItr{({{ .Type | TypeSymbol }})(nr), 0}
		}

		type _{{ .Type | TypeSymbol }}__ReprMapItr _{{ .Type | TypeSymbol }}__MapItr

		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Next() (k ipld.Node, v ipld.Node, err error) {
			k, v, err = (*_{{ .Type | TypeSymbol }}__MapItr)(itr).Next()
			if err != nil || v == ipld.Null {
				return
			}
			return k, v.({{ .Type.ValueType | TypeSymbol}}).Representation(), nil
		}
		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Done() bool {
			return (*_{{ .Type | TypeSymbol }}__MapItr)(itr).Done()
		}

	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func (rn *_{{ .Type | TypeSymbol }}__Repr) Length() int {
			return len(rn.t)
		}
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeMethodStyle(w io.Writer) {
	emitNodeMethodStyle_typical(w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) EmitNodeStyleType(w io.Writer) {
	// FIXME this alias is a lie, but keep it around for one more sec so we can do an incremental commit until we finish fixing it
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprStyle = _{{ .Type | TypeSymbol }}__Style
	`, w, g.AdjCfg, g)
}
func (g mapReprMapReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return nil // TODO
}
