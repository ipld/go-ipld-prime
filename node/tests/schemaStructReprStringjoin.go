package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

// TestStructReprStringjoin exercises... well, what it says on the tin.
//
// These should pass even if the natural map representation doesn't.
// No maybes are exercised.
func SchemaTestStructReprStringjoin(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("StringyStruct",
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	ts.Accumulate(schema.SpawnStruct("ManystringStruct",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, false),
			schema.SpawnStructField("bar", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	ts.Accumulate(schema.SpawnStruct("Recurzorator",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, false),
			schema.SpawnStructField("zap", "ManystringStruct", false, false),
			schema.SpawnStructField("bar", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin("-"),
	))
	engine.Init(t, ts)

	t.Run("single field works", func(t *testing.T) {
		np := engine.PrototypeByName("StringyStruct")
		nrp := engine.PrototypeByName("StringyStruct.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("field").AssignString("valoo")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.Kind(), ShouldEqual, datamodel.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(1))
				Wish(t, must.String(must.Node(n.LookupByString("field"))), ShouldEqual, "valoo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, datamodel.Kind_String)
				Wish(t, must.String(nr), ShouldEqual, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuild(nrp, func(na fluent.NodeAssembler) {
				na.AssignString("valoo")
			})
			Wish(t, datamodel.DeepEqual(n, nr), ShouldEqual, true)
		})
	})

	t.Run("several fields work", func(t *testing.T) {
		np := engine.PrototypeByName("ManystringStruct")
		nrp := engine.PrototypeByName("ManystringStruct.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignString("v1")
				ma.AssembleEntry("bar").AssignString("v2")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.Kind(), ShouldEqual, datamodel.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(2))
				Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "v1")
				Wish(t, must.String(must.Node(n.LookupByString("bar"))), ShouldEqual, "v2")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, datamodel.Kind_String)
				Wish(t, must.String(nr), ShouldEqual, "v1:v2")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuild(nrp, func(na fluent.NodeAssembler) {
				na.AssignString("v1:v2")
			})
			Wish(t, datamodel.DeepEqual(n, nr), ShouldEqual, true)
		})
	})

	t.Run("nested stringjoin structs work", func(t *testing.T) {
		np := engine.PrototypeByName("Recurzorator")
		nrp := engine.PrototypeByName("Recurzorator.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 3, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignString("v1")
				ma.AssembleEntry("zap").CreateMap(2, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("foo").AssignString("v2")
					ma.AssembleEntry("bar").AssignString("v3")
				})
				ma.AssembleEntry("bar").AssignString("v4")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				Require(t, n.Kind(), ShouldEqual, datamodel.Kind_Map)
				Wish(t, n.Length(), ShouldEqual, int64(3))
				Wish(t, must.String(must.Node(n.LookupByString("foo"))), ShouldEqual, "v1")
				Wish(t, must.String(must.Node(n.LookupByString("bar"))), ShouldEqual, "v4")
				n2 := must.Node(n.LookupByString("zap"))
				Wish(t, n2.Length(), ShouldEqual, int64(2))
				Wish(t, must.String(must.Node(n2.LookupByString("foo"))), ShouldEqual, "v2")
				Wish(t, must.String(must.Node(n2.LookupByString("bar"))), ShouldEqual, "v3")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, datamodel.Kind_String)
				Wish(t, must.String(nr), ShouldEqual, "v1-v2:v3-v4")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuild(nrp, func(na fluent.NodeAssembler) {
				na.AssignString("v1-v2:v3-v4")
			})
			Wish(t, datamodel.DeepEqual(n, nr), ShouldEqual, true)
		})
	})
}
