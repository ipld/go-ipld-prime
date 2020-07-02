package gengo

import (
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &unionReprKeyedGenerator{}

func NewUnionReprKeyedGenerator(pkgName string, typ schema.TypeUnion, adjCfg *AdjunctCfg) TypeGenerator {
	return unionReprKeyedGenerator{
		unionGenerator{
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

type unionReprKeyedGenerator struct {
	unionGenerator
}

func (g unionReprKeyedGenerator) GetRepresentationNodeGen() NodeGenerator {
	return nil /* TODO */
}
