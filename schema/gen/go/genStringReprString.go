package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &stringReprStringGenerator{}

func NewStringReprStringGenerator(pkgName string, typ schema.TypeString, symbolOverride string) TypeGenerator {
	symbol := string(typ.Name())
	if symbolOverride != "" {
		symbol = symbolOverride
	}
	return stringReprStringGenerator{
		stringGenerator{
			mixins.StringTraits{
				pkgName,
				string(typ.Name()),
				symbol,
			},
			pkgName,
			typ,
			symbol,
		},
	}
}

type stringReprStringGenerator struct {
	stringGenerator
}

func (g stringReprStringGenerator) GetRepresentationNodeGen() NodeGenerator {
	return stringReprStringReprGenerator{g.stringGenerator}
}

type stringReprStringReprGenerator struct {
	stringGenerator
}

func (g stringReprStringReprGenerator) EmitNodeType(w io.Writer) {
	// Since this is a "natural" representation... there's just a type alias here.
	//  No new functions are necessary.
	doTemplate(`
		type _{{ .Symbol }}__Repr = _{{ .Symbol }}
	`, w, g)
}
func (g stringReprStringReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Symbol }}__Repr{}
	`, w, g)

}
func (stringReprStringReprGenerator) EmitNodeMethodReprKind(io.Writer)      {}
func (stringReprStringReprGenerator) EmitNodeMethodLookupString(io.Writer)  {}
func (stringReprStringReprGenerator) EmitNodeMethodLookup(io.Writer)        {}
func (stringReprStringReprGenerator) EmitNodeMethodLookupIndex(io.Writer)   {}
func (stringReprStringReprGenerator) EmitNodeMethodLookupSegment(io.Writer) {}
func (stringReprStringReprGenerator) EmitNodeMethodMapIterator(io.Writer)   {}
func (stringReprStringReprGenerator) EmitNodeMethodListIterator(io.Writer)  {}
func (stringReprStringReprGenerator) EmitNodeMethodLength(io.Writer)        {}
func (stringReprStringReprGenerator) EmitNodeMethodIsUndefined(io.Writer)   {}
func (stringReprStringReprGenerator) EmitNodeMethodIsNull(io.Writer)        {}
func (stringReprStringReprGenerator) EmitNodeMethodAsBool(io.Writer)        {}
func (stringReprStringReprGenerator) EmitNodeMethodAsInt(io.Writer)         {}
func (stringReprStringReprGenerator) EmitNodeMethodAsFloat(io.Writer)       {}
func (stringReprStringReprGenerator) EmitNodeMethodAsString(io.Writer)      {}
func (stringReprStringReprGenerator) EmitNodeMethodAsBytes(io.Writer)       {}
func (stringReprStringReprGenerator) EmitNodeMethodAsLink(io.Writer)        {}
func (stringReprStringReprGenerator) EmitNodeMethodStyle(io.Writer)         {}
func (stringReprStringReprGenerator) EmitNodeStyleType(io.Writer)           {}
func (stringReprStringReprGenerator) EmitNodeBuilder(io.Writer)             {}
func (stringReprStringReprGenerator) EmitNodeAssembler(io.Writer)           {}
