package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionKinded(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("SmolStruct",
		[]schema.StructField{
			schema.SpawnStructField("s", "String", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"s": "q",
		}),
	))
	ts.Accumulate(schema.SpawnUnion("WheeUnion",
		[]schema.TypeName{
			"String",
			"SmolStruct",
		},
		schema.SpawnUnionRepresentationKinded(map[ipld.ReprKind]schema.TypeName{
			ipld.ReprKind_String: "String",
			ipld.ReprKind_Map:    "SmolStruct",
		}),
	))

	// These are the same *type-level* as in TestUnionKeyedComplexChildren,
	//  but (of course) have very different representations.
	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `"whee"`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"String", "whee"},
				//{"SmolStruct", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_String},
				{"", "whee"},
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"SmolStruct":{"s":"whee"}}`,
			reprJson: `{"q":"whee"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				//{"String", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"SmolStruct", ipld.ReprKind_Map},
				{"SmolStruct/s", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"q", "whee"},
			},
		},
	}

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("WheeUnion")
		nrp := getPrototypeByName("WheeUnion.Repr")
		for _, tcase := range specs {
			tcase.Test(t, np, nrp)
		}
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "embedAll"}

		prefix := "union-kinded-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "interface"}

		prefix := "union-kinded-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}
