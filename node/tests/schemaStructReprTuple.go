package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
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
				Require(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(1))
				Wish(t, must.String(must.Node(n.LookupByString("field"))), ShouldEqual, "valoo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
				Wish(t, nr.Length(), ShouldEqual, int64(1))
				Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 1, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("valoo")
			})
			Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
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
				Require(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(4))
				Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "0")
				Wish(t, must.String(must.Node(n.LookupByString("bar"))), ShouldEqual, "1")
				Wish(t, must.String(must.Node(n.LookupByString("baz"))), ShouldEqual, "2")
				Wish(t, must.String(must.Node(n.LookupByString("qux"))), ShouldEqual, "3")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
				Wish(t, nr.Length(), ShouldEqual, int64(4))
				Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "0")
				Wish(t, must.String(must.Node(nr.LookupByIndex(1))), ShouldEqual, "1")
				Wish(t, must.String(must.Node(nr.LookupByIndex(2))), ShouldEqual, "2")
				Wish(t, must.String(must.Node(nr.LookupByIndex(3))), ShouldEqual, "3")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignString("1")
				la.AssembleValue().AssignString("2")
				la.AssembleValue().AssignString("3")
			})
			Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
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
				Require(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(4))
				Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "0")
				Wish(t, must.Node(n.LookupByString("bar")), ShouldEqual, ipld.Null)
				Wish(t, must.Node(n.LookupByString("baz")), ShouldEqual, ipld.Absent)
				Wish(t, must.Node(n.LookupByString("qux")), ShouldEqual, ipld.Absent)
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
				Wish(t, nr.Length(), ShouldEqual, int64(2))
				Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "0")
				Wish(t, must.Node(nr.LookupByIndex(1)), ShouldEqual, ipld.Null)
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().AssignString("0")
				la.AssembleValue().AssignNull()
			})
			Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
		})
	})
}
