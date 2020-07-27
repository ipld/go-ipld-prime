package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionKeyed(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnString("Strung"))
	ts.Accumulate(schema.SpawnUnion("StrStr",
		[]schema.TypeName{
			"String",
			"Strung",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"a": "String",
			"b": "Strung",
		}),
	))

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("StrStr")
		nrp := getPrototypeByName("StrStr.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
				na.AssembleEntry("Strung").AssignString("whee")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 1)
				Wish(t, must.String(must.Node(n.LookupByString("Strung"))), ShouldEqual, "whee")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, nr.Length(), ShouldEqual, 1)
				Wish(t, must.String(must.Node(nr.LookupByString("b"))), ShouldEqual, "whee")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 2, func(na fluent.MapAssembler) {
				na.AssembleEntry("b").AssignString("whee")
			})
			Wish(t, n, ShouldEqual, nr)
		})
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"StrStr": "embedAll"}

		prefix := "union-keyed-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"StrStr": "interface"}

		prefix := "union-keyed-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}

// Test keyed unions again, but this time with more complex types as children.
//
// The previous tests used scalar types as the children; this exercises most things,
// but also has a couple (extremely non-obvious) simplifications:
// namely, because the default representation for strings are "natural" representations,
// the ReprAssemblers are actually aliases of the type-level Assemblers!
// Aaaand that makes a few things "work" by coincidence that wouldn't otherwise fly.
func TestUnionKeyedComplexChildren(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("SmolStruct",
		[]schema.StructField{
			schema.SpawnStructField("s", "String", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"s": "q",
		}),
	))
	ts.Accumulate(schema.SpawnUnion("WheeUnion",
		[]schema.TypeName{
			"String",
			"SmolStruct",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"a": "String",
			"b": "SmolStruct",
		}),
	))

	test := func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("WheeUnion")
		nrp := getPrototypeByName("WheeUnion.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
				na.AssembleEntry("SmolStruct").CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry("s").AssignString("whee")
				})
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 1)
				n2 := must.Node(n.LookupByString("SmolStruct"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, must.String(must.Node(n2.LookupByString("s"))), ShouldEqual, "whee")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, nr.Length(), ShouldEqual, 1)
				n2 := must.Node(nr.LookupByString("b"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, must.String(must.Node(n2.LookupByString("q"))), ShouldEqual, "whee")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 2, func(na fluent.MapAssembler) {
				na.AssembleEntry("b").CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry("q").AssignString("whee")
				})
			})
			Wish(t, n, ShouldEqual, nr)
		})
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "embedAll"}

		prefix := "union-keyed-complex-child-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "interface"}

		prefix := "union-keyed-complex-child-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}
