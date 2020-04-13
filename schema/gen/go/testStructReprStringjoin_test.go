package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

// TestStructReprStringjoin exercises... well, what it says on the tin.
//
// These should pass even if the natural map representation doesn't.
// No maybes are exercised.
func TestStructReprStringjoin(t *testing.T) {
	prefix := "structstrjoin"
	pkgName := "main"

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("StringyStruct",
		[]schema.StructField{
			schema.SpawnStructField("field", ts.TypeByName("String"), false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	ts.Accumulate(schema.SpawnStruct("ManystringStruct",
		[]schema.StructField{
			schema.SpawnStructField("foo", ts.TypeByName("String"), false, false),
			schema.SpawnStructField("bar", ts.TypeByName("String"), false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	// TODO coming soon:
	// ts.Accumulate(schema.SpawnStruct("Recurzorator",
	// 	[]schema.StructField{
	// 		schema.SpawnStructField("foo", ts.TypeByName("String"), false, false),
	// 		schema.SpawnStructField("zap", ts.TypeByName("ManystringStruct"), false, false),
	// 		schema.SpawnStructField("bar", ts.TypeByName("String"), false, false),
	// 	},
	// 	schema.SpawnStructRepresentationStringjoin("-"),
	// ))

	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		t.Run("single field works", func(t *testing.T) {
			ns := getStyleByName("StringyStruct")
			nsr := getStyleByName("StringyStruct.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(ns, 1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("field").AssignString("valoo")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, 1)
					Wish(t, must.String(must.Node(n.LookupString("field"))), ShouldEqual, "valoo")
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_String)
					Wish(t, must.String(nr), ShouldEqual, "valoo")
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuild(nsr, func(na fluent.NodeAssembler) {
					na.AssignString("valoo")
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})

		t.Run("several fields work", func(t *testing.T) {
			ns := getStyleByName("ManystringStruct")
			nsr := getStyleByName("ManystringStruct.Repr")
			var n schema.TypedNode
			t.Run("typed-create", func(t *testing.T) {
				n = fluent.MustBuildMap(ns, 2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("v1")
					ma.AssembleEntry("bar").AssignString("v2")
				}).(schema.TypedNode)
				t.Run("typed-read", func(t *testing.T) {
					Require(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
					Wish(t, n.Length(), ShouldEqual, 2)
					Wish(t, must.String(must.Node(n.LookupString("foo"))), ShouldEqual, "v1")
					Wish(t, must.String(must.Node(n.LookupString("bar"))), ShouldEqual, "v2")
				})
				t.Run("repr-read", func(t *testing.T) {
					nr := n.Representation()
					Require(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_String)
					Wish(t, must.String(nr), ShouldEqual, "v1:v2")
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				nr := fluent.MustBuild(nsr, func(na fluent.NodeAssembler) {
					na.AssignString("v1:v2")
				})
				Wish(t, n, ShouldEqual, nr)
			})
		})
	})
}
