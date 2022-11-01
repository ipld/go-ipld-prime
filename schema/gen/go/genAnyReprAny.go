package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

var _ TypeGenerator = &anyReprAnyGenerator{}

func NewAnyReprAnyGenerator(pkgName string, typ *schema.TypeAny, adjCfg *AdjunctCfg) TypeGenerator {
	return anyReprAnyGenerator{
		anyGenerator{
			adjCfg,
			pkgName,
			typ,
		},
	}
}

type anyReprAnyGenerator struct {
	anyGenerator
}

func (g anyReprAnyGenerator) EmitNodeMethodKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Kind() datamodel.Kind {
			return datamodel.Kind_Invalid
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodLookupByString(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) LookupByString(string) (datamodel.Node, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodLookupByNode(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) LookupByNode(datamodel.Node) (datamodel.Node, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodLookupByIndex(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) LookupByIndex(idx int64) (datamodel.Node, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodLookupBySegment(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) MapIterator() datamodel.MapIterator {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodListIterator(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) ListIterator() datamodel.ListIterator {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Length() int64 {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodIsAbsent(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) IsAbsent() bool {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodIsNull(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) IsNull() bool {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsInt() (int64, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsFloat(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsFloat() (float64, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsString() (string, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsBytes(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsBytes() ([]byte, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsLink(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsLink() (datamodel.Link, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) EmitNodeMethodAsBool(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) AsBool() (bool, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}

func (g anyReprAnyGenerator) GetRepresentationNodeGen() NodeGenerator {
	return anyReprAnyReprGenerator{
		g.AdjCfg,
		g.Type,
	}
}

type anyReprAnyReprGenerator struct {
	AdjCfg *AdjunctCfg
	Type   *schema.TypeAny
}

func (g anyReprAnyReprGenerator) EmitNodeType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr = _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}
func (g anyReprAnyReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ datamodel.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}
func (anyReprAnyReprGenerator) EmitNodeMethodKind(io.Writer)            {}
func (anyReprAnyReprGenerator) EmitNodeMethodLookupByString(io.Writer)  {}
func (anyReprAnyReprGenerator) EmitNodeMethodLookupByNode(io.Writer)    {}
func (anyReprAnyReprGenerator) EmitNodeMethodLookupByIndex(io.Writer)   {}
func (anyReprAnyReprGenerator) EmitNodeMethodLookupBySegment(io.Writer) {}
func (anyReprAnyReprGenerator) EmitNodeMethodMapIterator(io.Writer)     {}
func (anyReprAnyReprGenerator) EmitNodeMethodListIterator(io.Writer)    {}
func (anyReprAnyReprGenerator) EmitNodeMethodLength(io.Writer)          {}
func (anyReprAnyReprGenerator) EmitNodeMethodIsAbsent(io.Writer)        {}
func (anyReprAnyReprGenerator) EmitNodeMethodIsNull(io.Writer)          {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsBool(io.Writer)          {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsInt(io.Writer)           {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsFloat(io.Writer)         {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsString(io.Writer)        {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsBytes(io.Writer)         {}
func (anyReprAnyReprGenerator) EmitNodeMethodAsLink(io.Writer)          {}
func (anyReprAnyReprGenerator) EmitNodeMethodPrototype(io.Writer)       {}
func (g anyReprAnyReprGenerator) EmitNodePrototypeType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprPrototype = _{{ .Type | TypeSymbol }}__Prototype
	`, w, g.AdjCfg, g)
}
func (g anyReprAnyReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return anyReprAnyReprBuilderGenerator(g)
}

type anyReprAnyReprBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	Type   *schema.TypeAny
}

func (anyReprAnyReprBuilderGenerator) EmitNodeBuilderType(io.Writer)    {}
func (anyReprAnyReprBuilderGenerator) EmitNodeBuilderMethods(io.Writer) {}
func (g anyReprAnyReprBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprAssembler = _{{ .Type | TypeSymbol }}__Assembler
	`, w, g.AdjCfg, g)
}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodBeginMap(io.Writer)     {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodBeginList(io.Writer)    {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignNull(io.Writer)   {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignBool(io.Writer)   {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignInt(io.Writer)    {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignFloat(io.Writer)  {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignString(io.Writer) {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignBytes(io.Writer)  {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignLink(io.Writer)   {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodAssignNode(io.Writer)   {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerMethodPrototype(io.Writer)    {}
func (anyReprAnyReprBuilderGenerator) EmitNodeAssemblerOtherBits(io.Writer)          {}
