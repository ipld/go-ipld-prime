package gengo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/ipld/go-ipld-prime/schema"
)

// Generate takes a typesystem and the adjunct config for codegen,
// and emits generated code in the given path with the given package name.
//
// All of the files produced will match the pattern "ipldsch.*.gen.go".
func Generate(pth string, pkgName string, ts schema.TypeSystem, adjCfg *AdjunctCfg) {
	// Emit fixed bits.
	withFile(filepath.Join(pth, "ipldsch_minima.go"), func(f io.Writer) {
		EmitInternalEnums(pkgName, f)
	})

	externs, err := getExternTypes(pth)
	if err != nil {
		// Consider warning that duplication may be present due to inability to parse destination.
		externs = make(map[string]struct{})
	}

	// Local helper function for applying generation logic to each type.
	//  We will end up doing this more than once because in this layout, more than one file contains part of the story for each type.
	applyToEachType := func(fn func(tg TypeGenerator, w io.Writer), f io.Writer) {
		// Sort the type names so we have a determinisic order; this affects output consistency.
		//  Any stable order would do, but we don't presently have one, so a sort is necessary.
		types := ts.GetTypes()
		keys := make(sortableTypeNames, 0, len(types))
		for tn := range types {
			if _, exists := externs[tn.String()]; !exists {
				keys = append(keys, tn)
			}
		}
		sort.Sort(keys)
		for _, tn := range keys {
			switch t2 := types[tn].(type) {
			case *schema.TypeBool:
				fn(NewBoolReprBoolGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeInt:
				fn(NewIntReprIntGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeFloat:
				fn(NewFloatReprFloatGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeString:
				fn(NewStringReprStringGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeBytes:
				fn(NewBytesReprBytesGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeLink:
				fn(NewLinkReprLinkGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeStruct:
				switch t2.RepresentationStrategy().(type) {
				case schema.StructRepresentation_Map:
					fn(NewStructReprMapGenerator(pkgName, t2, adjCfg), f)
				case schema.StructRepresentation_Tuple:
					fn(NewStructReprTupleGenerator(pkgName, t2, adjCfg), f)
				case schema.StructRepresentation_Stringjoin:
					fn(NewStructReprStringjoinGenerator(pkgName, t2, adjCfg), f)
				default:
					panic("unrecognized struct representation strategy")
				}
			case *schema.TypeMap:
				fn(NewMapReprMapGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeList:
				fn(NewListReprListGenerator(pkgName, t2, adjCfg), f)
			case *schema.TypeUnion:
				switch t2.RepresentationStrategy().(type) {
				case schema.UnionRepresentation_Keyed:
					fn(NewUnionReprKeyedGenerator(pkgName, t2, adjCfg), f)
				case schema.UnionRepresentation_Kinded:
					fn(NewUnionReprKindedGenerator(pkgName, t2, adjCfg), f)
				default:
					panic("unrecognized union representation strategy")
				}
			default:
				panic("add more type switches here :)")
			}
		}
	}

	// Emit a file with the type table, and the golang type defns for each type.
	withFile(filepath.Join(pth, "ipldsch_types.go"), func(f io.Writer) {
		// Emit headers, import statements, etc.
		fmt.Fprintf(f, "package %s\n\n", pkgName)
		fmt.Fprintf(f, doNotEditComment+"\n\n")
		fmt.Fprintf(f, "import (\n")
		fmt.Fprintf(f, "\tipld \"github.com/ipld/go-ipld-prime\"\n") // referenced for links
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "var _ ipld.Node = nil // suppress errors when this dependency is not referenced\n")

		// Emit the type table.
		EmitTypeTable(pkgName, ts, adjCfg, f)

		// Emit the type defns matching the schema types.
		fmt.Fprintf(f, "\n// --- type definitions follow ---\n\n")
		applyToEachType(func(tg TypeGenerator, w io.Writer) {
			tg.EmitNativeType(w)
			fmt.Fprintf(f, "\n")
		}, f)

	})

	// Emit a file with all the Node/NodeBuilder/NodeAssembler boilerplate.
	//  Also includes typedefs for representation-level data.
	//  Also includes the MaybeT boilerplate.
	withFile(filepath.Join(pth, "ipldsch_satisfaction.go"), func(f io.Writer) {
		// Emit headers, import statements, etc.
		fmt.Fprintf(f, "package %s\n\n", pkgName)
		fmt.Fprintf(f, doNotEditComment+"\n\n")
		fmt.Fprintf(f, "import (\n")
		fmt.Fprintf(f, "\tipld \"github.com/ipld/go-ipld-prime\"\n")        // referenced everywhere.
		fmt.Fprintf(f, "\t\"github.com/ipld/go-ipld-prime/node/mixins\"\n") // referenced by node implementation guts.
		fmt.Fprintf(f, "\t\"github.com/ipld/go-ipld-prime/schema\"\n")      // referenced by maybes (and surprisingly little else).
		fmt.Fprintf(f, ")\n\n")

		// For each type, we'll emit... everything except the native type, really.
		applyToEachType(func(tg TypeGenerator, w io.Writer) {
			tg.EmitNativeAccessors(w)
			tg.EmitNativeBuilder(w)
			tg.EmitNativeMaybe(w)
			EmitNode(tg, w)
			tg.EmitTypedNodeMethodType(w)
			tg.EmitTypedNodeMethodRepresentation(w)

			nrg := tg.GetRepresentationNodeGen()
			EmitNode(nrg, w)

			fmt.Fprintf(f, "\n")
		}, f)
	})
}

// GenerateSplayed is like Generate, but emits a differnet pattern of files.
// GenerateSplayed emits many more individual files than Generate.
//
// This function should be considered deprecated and may be removed in the future.
func GenerateSplayed(pth string, pkgName string, ts schema.TypeSystem, adjCfg *AdjunctCfg) {
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
				case schema.UnionRepresentation_Kinded:
					EmitEntireType(NewUnionReprKindedGenerator(pkgName, t2, adjCfg), f)
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
		fmt.Fprintf(f, "package %s\n\n", pkgName)
		fmt.Fprintf(f, doNotEditComment+"\n\n")
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

type sortableTypeNames []schema.TypeName

func (a sortableTypeNames) Len() int           { return len(a) }
func (a sortableTypeNames) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortableTypeNames) Less(i, j int) bool { return a[i] < a[j] }
