package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestStructReprTuple(t *testing.T) {
	t.Parallel()

	prefix := "structtuple"
	pkgName := "main"

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
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

	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		t.Run("onetuple works", func(t *testing.T) {
			np := getPrototypeByName("OneTuple")
			nrp := getPrototypeByName("OneTuple.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("field").AssignString("valoo")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, 1)
					Wish(t, must.String(must.Node(n.LookupByString("field"))), ShouldEqual, "valoo")
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, nr.Length(), ShouldEqual, 1)
					Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "valoo")
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nrp, 1, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("valoo")
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})

		t.Run("fourtuple works", func(t *testing.T) {
			np := getPrototypeByName("FourTuple")
			nrp := getPrototypeByName("FourTuple.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(np, 4, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("0")
					ma.AssembleEntry("bar").AssignString("1")
					ma.AssembleEntry("baz").AssignString("2")
					ma.AssembleEntry("qux").AssignString("3")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, 4)
					Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "0")
					Wish(t, must.String(must.Node(n.LookupByString("bar"))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(n.LookupByString("baz"))), ShouldEqual, "2")
					Wish(t, must.String(must.Node(n.LookupByString("qux"))), ShouldEqual, "3")
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, nr.Length(), ShouldEqual, 4)
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
				Wish(t, n, ShouldEqual, nr)
			})
		})

		t.Run("fourtuple with absents", func(t *testing.T) {
			np := getPrototypeByName("FourTuple")
			nrp := getPrototypeByName("FourTuple.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("0")
					ma.AssembleEntry("bar").AssignNull()
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, 4)
					Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "0")
					Wish(t, must.Node(n.LookupByString("bar")), ShouldEqual, ipld.Null)
					Wish(t, must.Node(n.LookupByString("baz")), ShouldEqual, ipld.Absent)
					Wish(t, must.Node(n.LookupByString("qux")), ShouldEqual, ipld.Absent)
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, nr.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "0")
					Wish(t, must.Node(nr.LookupByIndex(1)), ShouldEqual, ipld.Null)
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("0")
					la.AssembleValue().AssignNull()
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
	})
}
