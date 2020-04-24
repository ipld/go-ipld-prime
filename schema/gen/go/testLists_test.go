package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestListsContainingMaybe(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnList("List__String",
		ts.TypeByName("String"), false))
	ts.Accumulate(schema.SpawnList("List__nullableString",
		ts.TypeByName("String"), true))

	test := func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		t.Run("non-nullable", func(t *testing.T) {
			ns := getStyleByName("List__String")
			nsr := getStyleByName("List__String.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildList(ns, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignString("2")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, n.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(n.LookupIndex(0))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(n.LookupIndex(1))), ShouldEqual, "2")
					_, err := n.LookupIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, nr.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(nr.LookupIndex(0))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(nr.LookupIndex(1))), ShouldEqual, "2")
					_, err := n.LookupIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nsr, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignString("2")
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
		t.Run("nullable", func(t *testing.T) {
			ns := getStyleByName("List__nullableString")
			nsr := getStyleByName("List__nullableString.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildList(ns, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignNull()
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, n.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(n.LookupIndex(0))), ShouldEqual, "1")
					Wish(t, must.Node(n.LookupIndex(1)), ShouldEqual, ipld.Null)
					_, err := n.LookupIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
					Wish(t, nr.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(n.LookupIndex(0))), ShouldEqual, "1")
					Wish(t, must.Node(n.LookupIndex(1)), ShouldEqual, ipld.Null)
					_, err := n.LookupIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nsr, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignNull()
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
	}

	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "lists-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			test(t, getStyleByName)
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "lists-mptr"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			test(t, getStyleByName)
		})
	})
}
