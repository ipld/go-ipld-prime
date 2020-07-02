package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionPoke(t *testing.T) {
	t.Skip("not implemented yet")

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnString("Strung"))
	ts.Accumulate(schema.SpawnUnion("StrStr",
		[]schema.Type{
			ts.TypeByName("String"),
			ts.TypeByName("Strung"),
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.Type{
			"a": ts.TypeByName("String"),
			"b": ts.TypeByName("Strung"),
		},
		)))

	prefix := "union-keyed"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
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
				Wish(t, must.String(must.Node(n.LookupByString("b"))), ShouldEqual, "whee")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 2, func(na fluent.MapAssembler) {
				na.AssembleEntry("b").AssignString("whee")
			})
			Wish(t, n, ShouldEqual, nr)
		})
	})
}
