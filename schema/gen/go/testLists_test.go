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
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnList("List__String",
		"String", false))
	ts.Accumulate(schema.SpawnList("List__nullableString",
		"String", true))

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		t.Run("non-nullable", func(t *testing.T) {
			np := getPrototypeByName("List__String")
			nrp := getPrototypeByName("List__String.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildList(np, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignString("2")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.Kind(), ShouldEqual, ipld.Kind_List)
					Wish(t, n.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(n.LookupByIndex(0))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(n.LookupByIndex(1))), ShouldEqual, "2")
					_, err := n.LookupByIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
					Wish(t, nr.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(nr.LookupByIndex(0))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(nr.LookupByIndex(1))), ShouldEqual, "2")
					_, err := n.LookupByIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignString("2")
				})
				Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
			})
		})
		t.Run("nullable", func(t *testing.T) {
			np := getPrototypeByName("List__nullableString")
			nrp := getPrototypeByName("List__nullableString.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildList(np, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignNull()
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.Kind(), ShouldEqual, ipld.Kind_List)
					Wish(t, n.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(n.LookupByIndex(0))), ShouldEqual, "1")
					Wish(t, must.Node(n.LookupByIndex(1)), ShouldEqual, ipld.Null)
					_, err := n.LookupByIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
					Wish(t, nr.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(n.LookupByIndex(0))), ShouldEqual, "1")
					Wish(t, must.Node(n.LookupByIndex(1)), ShouldEqual, ipld.Null)
					_, err := n.LookupByIndex(3)
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignNull()
				})
				Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
			})
		})
	}

	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "lists-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "lists-mptr"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}

// TestListsContainingLists is probing *two* things:
//   - that lists can nest, obviously
//   - that representation semantics are held correctly when we recurse, both in builders and in reading
// To cover that latter situation, this depends on structs (so we can use rename directives on the representation to make it distinctive).
func TestListsContainingLists(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
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

	prefix := "lists-of-lists"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("List__List__Frub")
		nrp := getPrototypeByName("List__List__Frub.Repr")
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
				Require(t, n.Kind(), ShouldEqual, ipld.Kind_List)
				Require(t, n.Length(), ShouldEqual, int64(3))
				Require(t, must.Node(n.LookupByIndex(0)).Length(), ShouldEqual, int64(3))
				Require(t, must.Node(n.LookupByIndex(1)).Length(), ShouldEqual, int64(1))
				Require(t, must.Node(n.LookupByIndex(2)).Length(), ShouldEqual, int64(2))

				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(0)).LookupByIndex(0)).LookupByString("field"))), ShouldEqual, "11")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(0)).LookupByIndex(2)).LookupByString("field"))), ShouldEqual, "13")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(1)).LookupByIndex(0)).LookupByString("field"))), ShouldEqual, "21")
				Wish(t, must.String(must.Node(must.Node(must.Node(n.LookupByIndex(2)).LookupByIndex(1)).LookupByString("field"))), ShouldEqual, "32")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, ipld.Kind_List)
				Require(t, nr.Length(), ShouldEqual, int64(3))
				Require(t, must.Node(nr.LookupByIndex(0)).Length(), ShouldEqual, int64(3))
				Require(t, must.Node(nr.LookupByIndex(1)).Length(), ShouldEqual, int64(1))
				Require(t, must.Node(nr.LookupByIndex(2)).Length(), ShouldEqual, int64(2))

				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(0)).LookupByIndex(0)).LookupByString("encoded"))), ShouldEqual, "11")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(0)).LookupByIndex(2)).LookupByString("encoded"))), ShouldEqual, "13")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(1)).LookupByIndex(0)).LookupByString("encoded"))), ShouldEqual, "21")
				Wish(t, must.String(must.Node(must.Node(must.Node(nr.LookupByIndex(2)).LookupByIndex(1)).LookupByString("encoded"))), ShouldEqual, "32")
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
			Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
		})
	})

}
