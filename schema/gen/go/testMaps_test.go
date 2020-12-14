package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestMapsContainingMaybe(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnMap("Map__String__String",
		"String", "String", false))
	ts.Accumulate(schema.SpawnMap("Map__String__nullableString",
		"String", "String", true))

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		t.Run("non-nullable", func(t *testing.T) {
			np := getPrototypeByName("Map__String__String")
			nrp := getPrototypeByName("Map__String__String.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("one").AssignString("1")
					ma.AssembleEntry("two").AssignString("2")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(n.LookupByString("one"))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(n.LookupByString("two"))), ShouldEqual, "2")
					_, err := n.LookupByString("miss")
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, nr.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(nr.LookupByString("one"))), ShouldEqual, "1")
					Wish(t, must.String(must.Node(nr.LookupByString("two"))), ShouldEqual, "2")
					_, err := nr.LookupByString("miss")
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildMap(nrp, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("one").AssignString("1")
					ma.AssembleEntry("two").AssignString("2")
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
		t.Run("nullable", func(t *testing.T) {
			np := getPrototypeByName("Map__String__nullableString")
			nrp := getPrototypeByName("Map__String__nullableString.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("one").AssignString("1")
					ma.AssembleEntry("none").AssignNull()
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(n.LookupByString("one"))), ShouldEqual, "1")
					Wish(t, must.Node(n.LookupByString("none")), ShouldEqual, ipld.Null)
					_, err := n.LookupByString("miss")
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, nr.Length(), ShouldEqual, int64(2))
					Wish(t, must.String(must.Node(nr.LookupByString("one"))), ShouldEqual, "1")
					Wish(t, must.Node(nr.LookupByString("none")), ShouldEqual, ipld.Null)
					_, err := nr.LookupByString("miss")
					Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuildMap(nrp, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("one").AssignString("1")
					ma.AssembleEntry("none").AssignNull()
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
	}

	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "maps-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "maps-mptr"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}

// TestMapsContainingMaps is probing *two* things:
//   - that maps can nest, obviously
//   - that representation semantics are held correctly when we recurse, both in builders and in reading
// To cover that latter situation, this depends on structs (so we can use rename directives on the representation to make it distinctive).
func TestMapsContainingMaps(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Frub", // "type Frub struct { field String (rename "encoded") }"
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false), // plain field.
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"field": "encoded",
		}),
	))
	ts.Accumulate(schema.SpawnMap("Map__String__Frub", // "{String:Frub}"
		"String", "Frub", false))
	ts.Accumulate(schema.SpawnMap("Map__String__nullableMap__String__Frub", // "{String:nullable {String:Frub}}"
		"String", "Map__String__Frub", true))

	prefix := "maps-recursive"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("Map__String__nullableMap__String__Frub")
		nrp := getPrototypeByName("Map__String__nullableMap__String__Frub.Repr")
		creation := func(t *testing.T, np ipld.NodePrototype, fieldName string) schema.TypedNode {
			return fluent.MustBuildMap(np, 3, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("one").CreateMap(2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("zot").CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry(fieldName).AssignString("11") })
					ma.AssembleEntry("zop").CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry(fieldName).AssignString("12") })
				})
				ma.AssembleEntry("two").CreateMap(1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("zim").CreateMap(1, func(ma fluent.MapAssembler) { ma.AssembleEntry(fieldName).AssignString("21") })
				})
				ma.AssembleEntry("none").AssignNull()
			}).(schema.TypedNode)
		}
		reading := func(t *testing.T, n ipld.Node, fieldName string) {
			withNode(n, func(n ipld.Node) {
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(3))
				withNode(must.Node(n.LookupByString("one")), func(n ipld.Node) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, int64(2))
					withNode(must.Node(n.LookupByString("zot")), func(n ipld.Node) {
						Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
						Wish(t, n.Length(), ShouldEqual, int64(1))
						Wish(t, must.String(must.Node(n.LookupByString(fieldName))), ShouldEqual, "11")
					})
					withNode(must.Node(n.LookupByString("zop")), func(n ipld.Node) {
						Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
						Wish(t, n.Length(), ShouldEqual, int64(1))
						Wish(t, must.String(must.Node(n.LookupByString(fieldName))), ShouldEqual, "12")
					})
				})
				withNode(must.Node(n.LookupByString("two")), func(n ipld.Node) {
					Wish(t, n.Length(), ShouldEqual, int64(1))
					withNode(must.Node(n.LookupByString("zim")), func(n ipld.Node) {
						Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
						Wish(t, n.Length(), ShouldEqual, int64(1))
						Wish(t, must.String(must.Node(n.LookupByString(fieldName))), ShouldEqual, "21")
					})
				})
				withNode(must.Node(n.LookupByString("none")), func(n ipld.Node) {
					Wish(t, n, ShouldEqual, ipld.Null)
				})
				_, err := n.LookupByString("miss")
				Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
			})
		}
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = creation(t, np, "field")
			t.Run("typed-read", func(t *testing.T) {
				reading(t, n, "field")
			})
			t.Run("repr-read", func(t *testing.T) {
				reading(t, n.Representation(), "encoded")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := creation(t, nrp, "encoded")
			Wish(t, n, ShouldEqual, nr)
		})
	})
}

func TestMapsWithComplexKeys(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("StringyStruct",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, false),
			schema.SpawnStructField("bar", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	ts.Accumulate(schema.SpawnMap("Map__StringyStruct__String",
		"StringyStruct", "String", false))

	prefix := "maps-cmplx-keys"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("Map__StringyStruct__String")
		nrp := getPrototypeByName("Map__StringyStruct__String.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 3, func(ma fluent.MapAssembler) {
				ma.AssembleKey().CreateMap(2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("a")
					ma.AssembleEntry("bar").AssignString("b")
				})
				ma.AssembleValue().AssignString("1")
				ma.AssembleKey().CreateMap(2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("c")
					ma.AssembleEntry("bar").AssignString("d")
				})
				ma.AssembleValue().AssignString("2")
				ma.AssembleKey().CreateMap(2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("e")
					ma.AssembleEntry("bar").AssignString("f")
				})
				ma.AssembleValue().AssignString("3")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(3))
				n2 := must.Node(n.LookupByString("c:d"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_String)
				Wish(t, must.String(n2), ShouldEqual, "2")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, nr.Length(), ShouldEqual, int64(3))
				n2 := must.Node(nr.LookupByString("c:d"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_String)
				Wish(t, must.String(n2), ShouldEqual, "2")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 3, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("a:b").AssignString("1")
				ma.AssembleEntry("c:d").AssignString("2")
				ma.AssembleEntry("e:f").AssignString("3")
			})
			Wish(t, n, ShouldEqual, nr)
		})
	})
}
