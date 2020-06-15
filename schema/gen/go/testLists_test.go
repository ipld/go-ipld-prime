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

// TestListsContainingLists is probing *two* things:
//   - that lists can nest, obviously
//   - that representation semantics are held correctly when we recurse, both in builders and in reading
// To cover that latter situation, this depends on structs (so we can use rename directives on the representation to make it distinctive).
func TestListsContainingLists(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Frub",
		[]schema.StructField{
			schema.SpawnStructField("field", ts.TypeByName("String"), false, false), // plain field.
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"field": "encoded",
		}),
	))
	ts.Accumulate(schema.SpawnList("List__Frub",
		ts.TypeByName("Frub"), false))
	ts.Accumulate(schema.SpawnList("List__List__Frub",
		ts.TypeByName("List__Frub"), true))

	prefix := "lists-of-lists"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		ns := getStyleByName("List__List__Frub")
		nsr := getStyleByName("List__List__Frub.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildList(ns, 3, func(la fluent.ListAssembler) {
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
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
				Require(t, n.Length(), ShouldEqual, 3)
				Require(t, must.Node(n.LookupIndex(0)).Length(), ShouldEqual, 3)
				Require(t, must.Node(n.LookupIndex(1)).Length(), ShouldEqual, 1)
				Require(t, must.Node(n.LookupIndex(2)).Length(), ShouldEqual, 2)

				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupIndex(0)).LookupIndex(0)).LookupString("field"))), ShouldEqual, "11")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupIndex(0)).LookupIndex(2)).LookupString("field"))), ShouldEqual, "13")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupIndex(1)).LookupIndex(0)).LookupString("field"))), ShouldEqual, "21")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupIndex(2)).LookupIndex(1)).LookupString("field"))), ShouldEqual, "32")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_List)
				Require(t, nr.Length(), ShouldEqual, 3)
				Require(t, must.Node(nr.LookupIndex(0)).Length(), ShouldEqual, 3)
				Require(t, must.Node(nr.LookupIndex(1)).Length(), ShouldEqual, 1)
				Require(t, must.Node(nr.LookupIndex(2)).Length(), ShouldEqual, 2)

				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupIndex(0)).LookupIndex(0)).LookupString("encoded"))), ShouldEqual, "11")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupIndex(0)).LookupIndex(2)).LookupString("encoded"))), ShouldEqual, "13")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupIndex(1)).LookupIndex(0)).LookupString("encoded"))), ShouldEqual, "21")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupIndex(2)).LookupIndex(1)).LookupString("encoded"))), ShouldEqual, "32")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nsr, 2, func(la fluent.ListAssembler) {
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
			Wish(t, n, ShouldEqual, nr)
		})
	})

}
