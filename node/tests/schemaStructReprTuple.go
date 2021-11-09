package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestStructReprTuple(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("OneTuple",
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false),
		},
		schema.SpawnStructRepresentationTuple(),
	))
	ts.Accumulate(schema.SpawnStruct("FourTuple",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, false),
			schema.SpawnStructField("bar", "String", false, true),
			schema.SpawnStructField("baz", "String", true, true),
			schema.SpawnStructField("qux", "String", true, false),
		},
		schema.SpawnStructRepresentationTuple(),
	))
	engine.Init(t, ts)

	t.Run("onetuple works", func(t *testing.T) {
		np := engine.PrototypeByName("OneTuple")
		nrp := engine.PrototypeByName("OneTuple.Repr")
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
				qt.Check(t, nr.Length(), qt.Equals, int64(1))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 1, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("valoo")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourtuple works", func(t *testing.T) {
		np := engine.PrototypeByName("FourTuple")
		nrp := engine.PrototypeByName("FourTuple.Repr")
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
				qt.Check(t, nr.Length(), qt.Equals, int64(4))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "0")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(1))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(2))), qt.Equals, "2")
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(3))), qt.Equals, "3")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("2")
				la.AssembleValue().AssignString("3")
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourtuple with absents", func(t *testing.T) {
		np := engine.PrototypeByName("FourTuple")
		nrp := engine.PrototypeByName("FourTuple.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignString("0")
				ma.AssembleEntry("bar").AssignNull()
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(4))
				qt.Check(t, must.String(must.Node(n.LookupByString("foo"))), qt.Equals, "0")
				qt.Check(t, must.Node(n.LookupByString("bar")), qt.Equals, datamodel.Null)
				qt.Check(t, must.Node(n.LookupByString("baz")), qt.Equals, datamodel.Absent)
				qt.Check(t, must.Node(n.LookupByString("qux")), qt.Equals, datamodel.Absent)
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(nr.LookupByIndex(0))), qt.Equals, "0")
				qt.Check(t, must.Node(nr.LookupByIndex(1)), qt.Equals, datamodel.Null)
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignNull()
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})
}
