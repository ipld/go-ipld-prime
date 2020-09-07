package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionKeyed(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnString("Strung"))
	ts.Accumulate(schema.SpawnUnion("StrStr",
		[]schema.TypeName{
			"String",
			"Strung",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"a": "String",
			"b": "Strung",
		}),
	))

	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `{"a":"whee"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"String", "whee"},
				//{"Strung", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"a", "whee"},
				//{"b", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"Strung":"whee"}`,
			reprJson: `{"b":"whee"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				//{"String", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"Strung", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				//{"a", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"b", "whee"},
			},
		},
	}

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("StrStr")
		nrp := getPrototypeByName("StrStr.Repr")
		for _, tcase := range specs {
			tcase.Test(t, np, nrp)
		}
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"StrStr": "embedAll"}

		prefix := "union-keyed-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"StrStr": "interface"}

		prefix := "union-keyed-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}

// Test keyed unions again, but this time with more complex types as children.
//
// The previous tests used scalar types as the children; this exercises most things,
// but also has a couple (extremely non-obvious) simplifications:
// namely, because the default representation for strings are "natural" representations,
// the ReprAssemblers are actually aliases of the type-level Assemblers!
// Aaaand that makes a few things "work" by coincidence that wouldn't otherwise fly.
func TestUnionKeyedComplexChildren(t *testing.T) {
	t.Parallel()

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
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"a": "String",
			"b": "SmolStruct",
		}),
	))

	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `{"a":"whee"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"String", "whee"},
				//{"SmolStruct", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"a", "whee"},
				//{"b", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"SmolStruct":{"s":"whee"}}`,
			reprJson: `{"b":{"q":"whee"}}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				//{"String", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"SmolStruct", ipld.ReprKind_Map},
				{"SmolStruct/s", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				//{"a", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"b", ipld.ReprKind_Map},
				{"b/q", "whee"},
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

		prefix := "union-keyed-complex-child-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "interface"}

		prefix := "union-keyed-complex-child-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}

// TestUnionKeyedReset puts a union inside a list, so that we can use the list's reuse of assembler as a test of the assembler's reset feature.
// The value inside the union is also more complex than a scalar value so that we test resetting gets passed down, too.
func TestUnionKeyedReset(t *testing.T) {
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
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"a": "String",
			"b": "SmolStruct",
		}),
	))
	ts.Accumulate(schema.SpawnList("OuterList",
		"WheeUnion", false,
	))

	specs := []testcase{
		{
			typeJson: `[{"SmolStruct":{"s":"one"}}, {"SmolStruct":{"s":"two"}}, {"String":"three"}]`,
			reprJson: `[{"b":{"q":"one"}}, {"b":{"q":"two"}}, {"a":"three"}]`,
			typePoints: []testcasePoint{
				{"0/SmolStruct/s", "one"},
				{"1/SmolStruct/s", "two"},
				{"2/String", "three"},
			},
			reprPoints: []testcasePoint{
				{"0/b/q", "one"},
				{"1/b/q", "two"},
				{"2/a", "three"},
			},
		},
	}

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("OuterList")
		nrp := getPrototypeByName("OuterList.Repr")
		for _, tcase := range specs {
			tcase.Test(t, np, nrp)
		}
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "embedAll"}

		prefix := "union-keyed-reset-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "interface"}

		prefix := "union-keyed-reset-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}
