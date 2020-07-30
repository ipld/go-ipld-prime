package gengo

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ipld/go-ipld-prime/schema"
)

func Generate(pth string, pkgName string, ts schema.TypeSystem, adjCfg *AdjunctCfg) {
	// Emit fixed bits.
	withFile(filepath.Join(pth, "minima.go"), func(f io.Writer) {
		EmitInternalEnums(pkgName, f)
	})

	// Emit a file for each type.
	for _, typ := range ts.GetTypes() {
		withFile(filepath.Join(pth, "t"+typ.Name().String()+".go"), func(f io.Writer) {
			EmitFileHeader(pkgName, f)
			switch t2 := typ.(type) {
			case *schema.TypeBool:
				EmitEntireType(NewBoolReprBoolGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeInt:
				EmitEntireType(NewIntReprIntGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeFloat:
				EmitEntireType(NewFloatReprFloatGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeString:
				EmitEntireType(NewStringReprStringGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeBytes:
				EmitEntireType(NewBytesReprBytesGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeLink:
				EmitEntireType(NewLinkReprLinkGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeStruct:
				switch t2.RepresentationStrategy().(type) {
				case schema.StructRepresentation_Map:
					EmitEntireType(NewStructReprMapGenerator(pkgName, t2, adjCfg), f)
				case schema.StructRepresentation_Tuple:
					EmitEntireType(NewStructReprTupleGenerator(pkgName, t2, adjCfg), f)
				case schema.StructRepresentation_Stringjoin:
					EmitEntireType(NewStructReprStringjoinGenerator(pkgName, t2, adjCfg), f)
				default:
					panic("unrecognized struct representation strategy")
				}
			case *schema.TypeMap:
				EmitEntireType(NewMapReprMapGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeList:
				EmitEntireType(NewListReprListGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeUnion:
				switch t2.RepresentationStrategy().(type) {
				case schema.UnionRepresentation_Keyed:
					EmitEntireType(NewUnionReprKeyedGenerator(pkgName, t2, adjCfg), f)
				default:
					panic("unrecognized union representation strategy")
				}
			default:
				panic("add more type switches here :)")
			}
		})
	}

	// Emit the unified type table.
	withFile(filepath.Join(pth, "typeTable.go"), func(f io.Writer) {
		EmitTypeTable(pkgName, ts, adjCfg, f)
	})
}

func withFile(filename string, fn func(io.Writer)) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fn(f)
}
