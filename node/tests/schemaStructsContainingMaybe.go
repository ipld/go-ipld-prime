package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

// TestStructsContainingMaybe checks all the variations of "nullable" and "optional" on struct fields.
// It does this twice: once for the child maybes being implemented with pointers,
// and once with maybes implemented as embeds.
// The child values are scalars.
//
// Both type-level generic build and access as well as representation build and access are exercised;
// the representation used is map (the native representation for structs).
func SchemaTestStructsContainingMaybe(t *testing.T, engine Engine) {
	// Type declarations.
	//  The tests here will all be targetted against this "Stroct" type.
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Stroct",
		[]schema.StructField{
			// Every field in this struct (including their order) is exercising an interesting case...
			schema.SpawnStructField("f1", "String", false, false), // plain field.
			schema.SpawnStructField("f2", "String", true, false),  // optional; later we have more than one optional field, nonsequentially.
			schema.SpawnStructField("f3", "String", false, true),  // nullable; but required.
			schema.SpawnStructField("f4", "String", true, true),   // optional and nullable; trailing optional.
			schema.SpawnStructField("f5", "String", true, false),  // optional; and the second one in a row, trailing.
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"f1": "r1",
			"f2": "r2",
			"f3": "r3",
			"f4": "r4",
		}),
	))
	engine.Init(t, ts)

	// There's a lot of cases to cover so a shorthand code labels each case for clarity:
	//  - 'v' -- value in that entry
	//  - 'n' -- null in that entry
	//  - 'z' -- absent entry
	// There's also a semantic description of the main detail being probed suffixed to the shortcode.
	specs := []testcase{
		{
			name:     "vvvvv-AllFieldsSet",
			typeJson: `{"f1":"a","f2":"b","f3":"c","f4":"d","f5":"e"}`,
			reprJson: `{"f5":"e","r1":"a","r2":"b","r3":"c","r4":"d"}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", "c"},
				{"f4", "d"},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"r1", "a"},
				{"r2", "b"},
				{"r3", "c"},
				{"r4", "d"},
				{"f5", "e"},
			},
		},
		{
			name:     "vvnnv-Nulls",
			typeJson: `{"f1":"a","f2":"b","f3":null,"f4":null,"f5":"e"}`,
			reprJson: `{"f5":"e","r1":"a","r2":"b","r3":null,"r4":null}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", datamodel.Null},
				{"f4", datamodel.Null},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"r1", "a"},
				{"r2", "b"},
				{"r3", datamodel.Null},
				{"r4", datamodel.Null},
				{"f5", "e"},
			},
		},
		{
			name:     "vzvzv-AbsentOptionals",
			typeJson: `{"f1":"a","f3":"c","f5":"e"}`,
			reprJson: `{"f5":"e","r1":"a","r3":"c"}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"f1", "a"},
				{"f2", datamodel.Absent},
				{"f3", "c"},
				{"f4", datamodel.Absent},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"r1", "a"},
				//{"r2", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"r3", "c"},
				//{"r4", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"f5", "e"},
			},
			typeItr: []entry{
				{"f1", "a"},
				{"f2", datamodel.Absent},
				{"f3", "c"},
				{"f4", datamodel.Absent},
				{"f5", "e"},
			},
		},
		{
			name:     "vvnzz-AbsentTrailingOptionals",
			typeJson: `{"f1":"a","f2":"b","f3":null}`,
			reprJson: `{"r1":"a","r2":"b","r3":null}`,
			typePoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", datamodel.Null},
				{"f4", datamodel.Absent},
				{"f5", datamodel.Absent},
			},
			reprPoints: []testcasePoint{
				{"", datamodel.Kind_Map},
				{"r1", "a"},
				{"r2", "b"},
				{"r3", datamodel.Null},
				//{"r4", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				//{"f5", datamodel.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			typeItr: []entry{
				{"f1", "a"},
				{"f2", "b"},
				{"f3", datamodel.Null},
				{"f4", datamodel.Absent},
				{"f5", datamodel.Absent},
			},
		},
	}

	for _, tcase := range specs {
		tcase.Test(t, engine.PrototypeByName("Stroct"), engine.PrototypeByName("Stroct.Repr"))
	}
}
