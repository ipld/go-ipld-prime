package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestStructNesting(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("SmolStruct",
		[]schema.StructField{
			schema.SpawnStructField("s", ts.TypeByName("String"), false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"s": "q",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("GulpoStruct",
		[]schema.StructField{
			schema.SpawnStructField("x", ts.TypeByName("SmolStruct"), false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"x": "r",
		}),
	))

	prefix := "struct-nesting"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		ns := getStyleByName("GulpoStruct")
		nsr := getStyleByName("GulpoStruct.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(ns, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("x").CreateMap(1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("s").AssignString("woo")
				})
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 1)
				n2 := must.Node(n.LookupString("x"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, must.String(must.Node(n2.LookupString("s"))), ShouldEqual, "woo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, nr.Length(), ShouldEqual, 1)
				n2 := must.Node(nr.LookupString("r"))
				Require(t, n2.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, must.String(must.Node(n2.LookupString("q"))), ShouldEqual, "woo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nsr, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("r").CreateMap(1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("q").AssignString("woo")
				})
			})
			Wish(t, n, ShouldEqual, nr)
		})
	})
}
