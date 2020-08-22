package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionKinded(t *testing.T) {
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
		schema.SpawnUnionRepresentationKinded(map[ipld.ReprKind]schema.TypeName{
			ipld.ReprKind_String: "String",
			ipld.ReprKind_Map:    "SmolStruct",
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
				Wish(t, must.String(must.Node(nr.LookupByString("q"))), ShouldEqual, "whee")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 1, func(na fluent.MapAssembler) {
				na.AssembleEntry("q").AssignString("whee")
			})
			Wish(t, n, ShouldEqual, nr)
		})
	}

	t.Run("union-using-embed", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "embedAll"}

		prefix := "union-kinded-using-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
	t.Run("union-using-interface", func(t *testing.T) {
		adjCfg.CfgUnionMemlayout = map[schema.TypeName]string{"WheeUnion": "interface"}

		prefix := "union-kinded-using-interface"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
			test(t, getPrototypeByName)
		})
	})
}
