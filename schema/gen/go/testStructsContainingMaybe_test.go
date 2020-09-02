package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

// TestStructsContainingMaybe checks all the variations of "nullable" and "optional" on struct fields.
// It does this twice: once for the child maybes being implemented with pointers,
// and once with maybes implemented as embeds.
// The child values are scalars.
//
// Both type-level generic build and access as well as representation build and access are exercised;
// the representation used is map (the native representation for structs).
func TestStructsContainingMaybe(t *testing.T) {
	// Type declarations.
	//  The tests here will all be targetted against this "Stroct" type.
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
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

	// There's a lot of cases to cover so a shorthand code labels each case for clarity:
	//  - 'v' -- value in that entry
	//  - 'n' -- null in that entry
	//  - 'z' -- absent entry
	// There's also a semantic description of the main detail being probed suffixed to the shortcode.
	specs := []testcase{
		{
			name:     "vvvvv-AllFieldsSet",
			typeJson: `{"f1":"a","f2":"b","f3":"c","f4":"d","f5":"e"}`,
			reprJson: `{"r1":"a","r2":"b","r3":"c","r4":"d","f5":"e"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", "c"},
				{"f4", "d"},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
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
			reprJson: `{"r1":"a","r2":"b","r3":null,"r4":null,"f5":"e"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", ipld.Null},
				{"f4", ipld.Null},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"r1", "a"},
				{"r2", "b"},
				{"r3", ipld.Null},
				{"r4", ipld.Null},
				{"f5", "e"},
			},
		},
		{
			name:     "vzvzv-AbsentOptionals",
			typeJson: `{"f1":"a","f3":"c","f5":"e"}`,
			reprJson: `{"r1":"a","r3":"c","f5":"e"}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"f1", "a"},
				{"f2", ipld.Absent},
				{"f3", "c"},
				{"f4", ipld.Absent},
				{"f5", "e"},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"r1", "a"},
				//{"r2", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"r3", "c"},
				//{"r4", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				{"f5", "e"},
			},
			typeItr: []entry{
				{"f1", "a"},
				{"f2", ipld.Absent},
				{"f3", "c"},
				{"f4", ipld.Absent},
				{"f5", "e"},
			},
		},
		{
			name:     "vvnzz-AbsentTrailingOptionals",
			typeJson: `{"f1":"a","f2":"b","f3":null}`,
			reprJson: `{"r1":"a","r2":"b","r3":null}`,
			typePoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"f1", "a"},
				{"f2", "b"},
				{"f3", ipld.Null},
				{"f4", ipld.Absent},
				{"f5", ipld.Absent},
			},
			reprPoints: []testcasePoint{
				{"", ipld.ReprKind_Map},
				{"r1", "a"},
				{"r2", "b"},
				{"r3", ipld.Null},
				//{"r4", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
				//{"f5", ipld.ErrNotExists{}}, // TODO: need better error typing from traversal package.
			},
			typeItr: []entry{
				{"f1", "a"},
				{"f2", "b"},
				{"f3", ipld.Null},
			},
		},
	}

	// And finally, launch tests! ...while specializing the adjunct config a bit.
	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "stroct"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			for _, tcase := range specs {
				tcase.Test(t, getPrototypeByName("Stroct"), getPrototypeByName("Stroct.Repr"))
			}
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "stroct2"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			for _, tcase := range specs {
				tcase.Test(t, getPrototypeByName("Stroct"), getPrototypeByName("Stroct.Repr"))
			}
		})
	})
}
