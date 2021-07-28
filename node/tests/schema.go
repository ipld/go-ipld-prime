package tests

import "testing"

// use a table instead of a map, to get a consistent order

var allSchemaTests = []struct {
	name string
	fn   func(*testing.T, Engine)
}{
	{"Links", SchemaTestLinks},
	{"ListsContainingMaybe", SchemaTestListsContainingMaybe},
	{"ListsContainingLists", SchemaTestListsContainingLists},
	{"MapsContainingMaybe", SchemaTestMapsContainingMaybe},
	{"MapsContainingMaps", SchemaTestMapsContainingMaps},
	{"MapsWithComplexKeys", SchemaTestMapsWithComplexKeys},
	{"Scalars", SchemaTestScalars},
	{"RequiredFields", SchemaTestRequiredFields},
	{"StructNesting", SchemaTestStructNesting},
	{"StructReprStringjoin", SchemaTestStructReprStringjoin},
	{"StructReprTuple", SchemaTestStructReprTuple},
	{"StructsContainingMaybe", SchemaTestStructsContainingMaybe},
	{"UnionKeyed", SchemaTestUnionKeyed},
	{"UnionKeyedComplexChildren", SchemaTestUnionKeyedComplexChildren},
	{"UnionKeyedReset", SchemaTestUnionKeyedReset},
	{"UnionKinded", SchemaTestUnionKinded},
	{"UnionStringprefix", SchemaTestUnionStringprefix},
}

type EngineSubtest struct {
	Name   string // subtest name
	Engine Engine
}

func SchemaTestAll(t *testing.T, forTest func(name string) []EngineSubtest) {
	for _, test := range allSchemaTests {
		test := test // do not reuse the range variable
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			subtests := forTest(test.name)
			if len(subtests) == 0 {
				t.Skip("no engine provided to SchemaTestAll")
			}
			if len(subtests) == 1 {
				sub := subtests[0]
				if sub.Name != "" {
					t.Fatal("a single engine shouldn't be named")
				}
				test.fn(t, sub.Engine)
				return
			}
			for _, sub := range subtests {
				if sub.Name == "" {
					t.Fatal("multiple engines should be named")
				}
				t.Run(sub.Name, func(t *testing.T) {
					test.fn(t, sub.Engine)
				})
			}
		})
	}
}
