package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestUnionKinded(t *testing.T, engine Engine) {
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
		schema.SpawnUnionRepresentationKinded(map[ipld.Kind]schema.TypeName{
			ipld.Kind_String: "String",
			ipld.Kind_Map:    "SmolStruct",
		}),
	))
	engine.Init(t, ts)

	// These are the same *type-level* as in TestUnionKeyedComplexChildren,
	//  but (of course) have very different representations.
	specs := []testcase{
		{
			name:     "InhabitantA",
			typeJson: `{"String":"whee"}`,
			reprJson: `"whee"`,
			typePoints: []testcasePoint{
				{"", ipld.Kind_Map},
				{"String", "whee"},
				//{"SmolStruct", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			reprPoints: []testcasePoint{
				{"", ipld.Kind_String},
				{"", "whee"},
			},
		},
		{
			name:     "InhabitantB",
			typeJson: `{"SmolStruct":{"s":"whee"}}`,
			reprJson: `{"q":"whee"}`,
			typePoints: []testcasePoint{
				{"", ipld.Kind_Map},
				//{"String", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"SmolStruct", ipld.Kind_Map},
				{"SmolStruct/s", "whee"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.Kind_Map},
				{"q", "whee"},
			},
		},
	}

	np := engine.PrototypeByName("WheeUnion")
	nrp := engine.PrototypeByName("WheeUnion.Repr")
	for _, tcase := range specs {
		tcase.Test(t, np, nrp)
	}
}
