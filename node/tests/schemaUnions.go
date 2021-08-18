package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestUnionKeyed(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
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
	engine.Init(t, ts)

	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `{"a":"whee"}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"String", "whee"},
				//{"Strung", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"a", "whee"},
				//{"b", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"Strung":"whee"}`,
			reprJson: `{"b":"whee"}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				//{"String", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"Strung", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				//{"a", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"b", "whee"},
			},
		},
	}

	np := engine.PrototypeByName("StrStr")
	nrp := engine.PrototypeByName("StrStr.Repr")
	for _, tcase := range specs {
		tcase.Test(t, np, nrp)
	}
}

// Test keyed unions again, but this time with more complex types as children.
//
// The previous tests used scalar types as the children; this exercises most things,
// but also has a couple (extremely non-obvious) simplifications:
// namely, because the default representation for strings are "natural" representations,
// the ReprAssemblers are actually aliases of the type-level Assemblers!
// Aaaand that makes a few things "work" by coincidence that wouldn't otherwise fly.
func SchemaTestUnionKeyedComplexChildren(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
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
	engine.Init(t, ts)

	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `{"a":"whee"}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"String", "whee"},
				//{"SmolStruct", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"a", "whee"},
				//{"b", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"SmolStruct":{"s":"whee"}}`,
			reprJson: `{"b":{"q":"whee"}}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				//{"String", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"SmolStruct", datamodel.Kind_Map},
				{"SmolStruct/s", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				//{"a", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"b", datamodel.Kind_Map},
				{"b/q", "whee"},
			},
		},
	}

	np := engine.PrototypeByName("WheeUnion")
	nrp := engine.PrototypeByName("WheeUnion.Repr")
	for _, tcase := range specs {
		tcase.Test(t, np, nrp)
	}
}

// TestUnionKeyedReset puts a union inside a list, so that we can use the list's reuse of assembler as a test of the assembler's reset feature.
// The value inside the union is also more complex than a scalar value so that we test resetting gets passed down, too.
func SchemaTestUnionKeyedReset(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
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
	engine.Init(t, ts)

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

	np := engine.PrototypeByName("OuterList")
	nrp := engine.PrototypeByName("OuterList.Repr")
	for _, tcase := range specs {
		tcase.Test(t, np, nrp)
	}
}
