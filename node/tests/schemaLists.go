package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestListsContainingMaybe(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnList("List__String",
		"String", false))
	ts.Accumulate(schema.SpawnList("List__nullableString",
		"String", true))
	engine.Init(t, ts)

	t.Run("non-nullable", func(t *testing.T) {
		np := engine.PrototypeByName("List__String")
		nrp := engine.PrototypeByName("List__String.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildList(np, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("2")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, n.Length(), qt.Equals, int64(2))

				qt.Check(t, must.String(must.Node(n.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(n.LookupByIndex(1))), qt.Equals, "2")

				qt.Check(t, must.String(must.Node(n.LookupBySegment(datamodel.PathSegmentOfInt(0)))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(n.LookupByNode(basicnode.NewInt(0)))), qt.Equals, "1")

				_, err := n.LookupByIndex(3)
				qt.Check(t, err, qt.ErrorAs, &datamodel.ErrNotExists{})
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))

				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(1))), qt.Equals, "2")

				qt.Check(t, must.String(must.Node(n.LookupBySegment(datamodel.PathSegmentOfInt(0)))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(n.LookupByNode(basicnode.NewInt(0)))), qt.Equals, "1")

				_, err := n.LookupByIndex(3)
				qt.Check(t, err, qt.ErrorAs, &datamodel.ErrNotExists{})
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("2")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})
	t.Run("nullable", func(t *testing.T) {
		np := engine.PrototypeByName("List__nullableString")
		nrp := engine.PrototypeByName("List__nullableString.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildList(np, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignNull()
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, n.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(n.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.Node(n.LookupByIndex(1)), qt.Equals, datamodel.Null)
				_, err := n.LookupByIndex(3)
				qt.Check(t, err, qt.ErrorAs, &datamodel.ErrNotExists{})
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(n.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.Node(n.LookupByIndex(1)), qt.Equals, datamodel.Null)
				_, err := n.LookupByIndex(3)
				qt.Check(t, err, qt.ErrorAs, &datamodel.ErrNotExists{})
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignNull()
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})
}

// TestListsContainingLists is probing *two* things:
//   - that lists can nest, obviously
//   - that representation semantics are held correctly when we recurse, both in builders and in reading
// To cover that latter situation, this depends on structs (so we can use rename directives on the representation to make it distinctive).
func SchemaTestListsContainingLists(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Frub",
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false), // plain field.
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"field": "encoded",
		}),
	))
	ts.Accumulate(schema.SpawnList("List__Frub",
		"Frub", false))
	ts.Accumulate(schema.SpawnList("List__List__Frub",
		"List__Frub", true))
	engine.Init(t, ts)

	np := engine.PrototypeByName("List__List__Frub")
	nrp := engine.PrototypeByName("List__List__Frub.Repr")
	var n schema.TypedNode
	t.Run("typed-create", func(t *testing.T) {
		n = fluent.MustBuildList(np, 3, func(la fluent.ListAssembler) {
			la.AssembleValue().CreateList(3, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("11") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("12") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("13") })
			})
			la.AssembleValue().CreateList(1, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("21") })
			})
			la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("31") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("field").AssignString("32") })
			})
		}).(schema.TypedNode)
		t.Run("typed-read", func(t *testing.T) {
			qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_List)
			qt.Assert(t, n.Length(), qt.Equals, int64(3))
			qt.Assert(t, must.Node(n.LookupByIndex(0)).Length(), qt.Equals, int64(3))
			qt.Assert(t, must.Node(n.LookupByIndex(1)).Length(), qt.Equals, int64(1))
			qt.Assert(t, must.Node(n.LookupByIndex(2)).Length(), qt.Equals, int64(2))

			qt.Check(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(0)).LookupByIndex(0)).LookupByString("field"))), qt.Equals, "11")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(0)).LookupByIndex(2)).LookupByString("field"))), qt.Equals, "13")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(1)).LookupByIndex(0)).LookupByString("field"))), qt.Equals, "21")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(2)).LookupByIndex(1)).LookupByString("field"))), qt.Equals, "32")
		})
		t.Run("repr-read", func(t *testing.T) {
			nr := n.Representation()
			qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
			qt.Assert(t, nr.Length(), qt.Equals, int64(3))
			qt.Assert(t, must.Node(nr.LookupByIndex(0)).Length(), qt.Equals, int64(3))
			qt.Assert(t, must.Node(nr.LookupByIndex(1)).Length(), qt.Equals, int64(1))
			qt.Assert(t, must.Node(nr.LookupByIndex(2)).Length(), qt.Equals, int64(2))

			qt.Check(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(0)).LookupByIndex(0)).LookupByString("encoded"))), qt.Equals, "11")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(0)).LookupByIndex(2)).LookupByString("encoded"))), qt.Equals, "13")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(1)).LookupByIndex(0)).LookupByString("encoded"))), qt.Equals, "21")
			qt.Check(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(2)).LookupByIndex(1)).LookupByString("encoded"))), qt.Equals, "32")
		})
	})
	t.Run("repr-create", func(t *testing.T) {
		nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
			// This is the same as the type-level create earlier, except note the field names are now all different.
			la.AssembleValue().CreateList(3, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("11") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("12") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("13") })
			})
			la.AssembleValue().CreateList(1, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("21") })
			})
			la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("31") })
				la.AssembleValue().CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry("encoded").AssignString("32") })
			})
		})
		qt.Check(t, n, NodeContentEquals, nr)
	})
}
