package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestStructReprListPairs(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("OneListPair",
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false),
		},
		schema.SpawnStructRepresentationListPairs(),
	))
	ts.Accumulate(schema.SpawnStruct("FourListPairs",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, true),
			schema.SpawnStructField("bar", "String", true, true),
			schema.SpawnStructField("baz", "String", true, false),
			schema.SpawnStructField("qux", "String", false, false),
		},
		schema.SpawnStructRepresentationListPairs(),
	))
	engine.Init(t, ts)

	t.Run("onelistpair works", func(t *testing.T) {
		np := engine.PrototypeByName("OneListPair")
		nrp := engine.PrototypeByName("OneListPair.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("field").AssignString("valoo")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(1))
				qt.Check(t, must.String(must.Node(n.LookupByString("field"))), qt.Equals, "valoo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "field")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(1))), qt.Equals, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("field")
				la.AssembleValue().AssignString("valoo")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourlistpairs works", func(t *testing.T) {
		np := engine.PrototypeByName("FourListPairs")
		nrp := engine.PrototypeByName("FourListPairs.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 4, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignString("0")
				ma.AssembleEntry("bar").AssignString("1")
				ma.AssembleEntry("baz").AssignString("2")
				ma.AssembleEntry("qux").AssignString("3")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(4))
				qt.Check(t, must.String(must.Node(n.LookupByString("foo"))), qt.Equals, "0")
				qt.Check(t, must.String(must.Node(n.LookupByString("bar"))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(n.LookupByString("baz"))), qt.Equals, "2")
				qt.Check(t, must.String(must.Node(n.LookupByString("qux"))), qt.Equals, "3")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(8))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "foo")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(1))), qt.Equals, "0")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(2))), qt.Equals, "bar")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(3))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(4))), qt.Equals, "baz")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(5))), qt.Equals, "2")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(6))), qt.Equals, "qux")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(7))), qt.Equals, "3")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 8, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("foo")
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignString("bar")
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("baz")
				la.AssembleValue().AssignString("2")
				la.AssembleValue().AssignString("qux")
				la.AssembleValue().AssignString("3")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
		t.Run("repr-create out-of-order", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 8, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("bar")
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("foo")
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignString("qux")
				la.AssembleValue().AssignString("3")
				la.AssembleValue().AssignString("baz")
				la.AssembleValue().AssignString("2")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourlistpairs with absents", func(t *testing.T) {
		np := engine.PrototypeByName("FourListPairs")
		nrp := engine.PrototypeByName("FourListPairs.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignNull()
				ma.AssembleEntry("qux").AssignString("3")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(4))
				qt.Check(t, must.Node(n.LookupByString("foo")), qt.Equals, datamodel.Null)
				qt.Check(t, must.Node(n.LookupByString("bar")), qt.Equals, datamodel.Absent)
				qt.Check(t, must.Node(n.LookupByString("baz")), qt.Equals, datamodel.Absent)
				qt.Check(t, must.String(must.Node(n.LookupByString("qux"))), qt.Equals, "3")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(4))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "foo")
				qt.Check(t, must.Node(nr.LookupByIndex(1)), qt.Equals, datamodel.Null)
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(2))), qt.Equals, "qux")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(3))), qt.Equals, "3")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("foo")
				la.AssembleValue().AssignNull()
				la.AssembleValue().AssignString("qux")
				la.AssembleValue().AssignString("3")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})
}
