package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestUnionStringprefix(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("SmolStruct",
		[]schema.StructField{
			schema.SpawnStructField("a", "String", false, false),
			schema.SpawnStructField("b", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	ts.Accumulate(schema.SpawnUnion("WheeUnion",
		[]schema.TypeName{
			"String",
			"SmolStruct",
		},
		schema.SpawnUnionRepresentationStringprefix(
			":",
			map[string]schema.TypeName{
				"simple":  "String",
				"complex": "SmolStruct",
			},
		),
	))
	engine.Init(t, ts)

	// These are the same *type-level* as in TestUnionKeyedComplexChildren,
	//  but (of course) have very different representations.
	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `"simple:whee"`,
			typePoints: []testcasePoint{
				{"", ipld.Kind_Map},
				{"String", "whee"},
				//{"SmolStruct", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", ipld.Kind_String},
				{"", "simple:whee"},
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"SmolStruct":{"a":"whee","b":"woo"}}`,
			reprJson: `"complex:whee:woo"`,
			typePoints: []testcasePoint{
				{"", ipld.Kind_Map},
				//{"String", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"SmolStruct", ipld.Kind_Map},
				{"SmolStruct/a", "whee"},
				{"SmolStruct/b", "woo"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.Kind_String},
				{"", "complex:whee:woo"},
			},
		},
	}

	np := engine.PrototypeByName("WheeUnion")
	nrp := engine.PrototypeByName("WheeUnion.Repr")
	for _, tcase := range specs {
		tcase.Test(t, np, nrp)
	}
}
