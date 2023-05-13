package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestStructReprListPairs(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnList("List__String", "String", false))
	ts.Accumulate(schema.SpawnStruct("OneListPair",
		[]schema.StructField{
			schema.SpawnStructField("field", "String", false, false),
		},
		schema.SpawnStructRepresentationListPairs(),
	))
	ts.Accumulate(schema.SpawnStruct("FourListPairs",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, true),
			schema.SpawnStructField("bar", "String", true, true),
			schema.SpawnStructField("baz", "String", true, false),
			schema.SpawnStructField("qux", "List__String", false, false),
		},
		schema.SpawnStructRepresentationListPairs(),
	))
	ts.Accumulate(schema.SpawnStruct("NestedListPairs",
		[]schema.StructField{
			schema.SpawnStructField("str", "String", false, false),
			schema.SpawnStructField("lp", "OneListPair", false, false),
		},
		schema.SpawnStructRepresentationListPairs(),
	))
	engine.Init(t, ts)

	t.Run("onelistpair works", func(t *testing.T) {
		np := engine.PrototypeByName("OneListPair")
		nrp := engine.PrototypeByName("OneListPair.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("field").AssignString("valoo")
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(1))
				qt.Check(t, must.String(must.Node(n.LookupByString("field"))), qt.Equals, "valoo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(1))
				kv := must.Node(nr.LookupByIndex(0))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "field")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 1, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("field")
					la.AssembleValue().AssignString("valoo")
				})
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourlistpairs works", func(t *testing.T) {
		np := engine.PrototypeByName("FourListPairs")
		nrp := engine.PrototypeByName("FourListPairs.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 4, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignString("0")
				ma.AssembleEntry("bar").AssignString("1")
				ma.AssembleEntry("baz").AssignString("2")
				ma.AssembleEntry("qux").CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("3")
					la.AssembleValue().AssignString("4")
				})
			}).(schema.TypedNode)

			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(4))
				qt.Check(t, must.String(must.Node(n.LookupByString("foo"))), qt.Equals, "0")
				qt.Check(t, must.String(must.Node(n.LookupByString("bar"))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(n.LookupByString("baz"))), qt.Equals, "2")
				qux := must.Node(n.LookupByString("qux"))
				qt.Assert(t, qux.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, qux.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(0))), qt.Equals, "3")
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(1))), qt.Equals, "4")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(4))
				kv := must.Node(nr.LookupByIndex(0))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "foo")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "0")
				kv = must.Node(nr.LookupByIndex(1))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "bar")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "1")
				kv = must.Node(nr.LookupByIndex(2))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "baz")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "2")
				kv = must.Node(nr.LookupByIndex(3))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "qux")
				qux := must.Node(kv.LookupByIndex(1))
				qt.Assert(t, qux.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, qux.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(0))), qt.Equals, "3")
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(1))), qt.Equals, "4")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("foo")
					la.AssembleValue().AssignString("0")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("bar")
					la.AssembleValue().AssignString("1")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("baz")
					la.AssembleValue().AssignString("2")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("qux")
					la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
						la.AssembleValue().AssignString("3")
						la.AssembleValue().AssignString("4")
					})
				})
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
		t.Run("repr-create out-of-order", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 4, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("bar")
					la.AssembleValue().AssignString("1")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("foo")
					la.AssembleValue().AssignString("0")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("qux")
					la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
						la.AssembleValue().AssignString("3")
						la.AssembleValue().AssignString("4")
					})
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("baz")
					la.AssembleValue().AssignString("2")
				})
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})

	t.Run("fourlistpairs with absents", func(t *testing.T) {
		np := engine.PrototypeByName("FourListPairs")
		nrp := engine.PrototypeByName("FourListPairs.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 2, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("foo").AssignNull()
				ma.AssembleEntry("qux").CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("1")
					la.AssembleValue().AssignString("2")
				})
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(4))
				qt.Check(t, must.Node(n.LookupByString("foo")), qt.Equals, datamodel.Null)
				qt.Check(t, must.Node(n.LookupByString("bar")), qt.Equals, datamodel.Absent)
				qt.Check(t, must.Node(n.LookupByString("baz")), qt.Equals, datamodel.Absent)
				qux := must.Node(n.LookupByString("qux"))
				qt.Assert(t, qux.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, qux.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(1))), qt.Equals, "2")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))
				kv := must.Node(nr.LookupByIndex(0))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "foo")
				qt.Check(t, must.Node(kv.LookupByIndex(1)), qt.Equals, datamodel.Null)
				kv = must.Node(nr.LookupByIndex(1))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "qux")
				qux := must.Node(kv.LookupByIndex(1))
				qt.Assert(t, qux.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, qux.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(0))), qt.Equals, "1")
				qt.Check(t, must.String(must.Node(qux.LookupByIndex(1))), qt.Equals, "2")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("foo")
					la.AssembleValue().AssignNull()
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("qux")
					la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
						la.AssembleValue().AssignString("1")
						la.AssembleValue().AssignString("2")
					})
				})
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
		t.Run("repr-create with AssignNode", func(t *testing.T) {
			nr := fluent.MustBuildList(basicnode.Prototype.Any, 2, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("foo")
					la.AssembleValue().AssignNull()
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("qux")
					la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
						la.AssembleValue().AssignString("1")
						la.AssembleValue().AssignString("2")
					})
				})
			})
			builder := nrp.NewBuilder()
			err := builder.AssignNode(nr)
			qt.Assert(t, err, qt.IsNil)
			anr := builder.Build()
			qt.Check(t, n, NodeContentEquals, anr)
		})
	})

	t.Run("nestedlistpairs works", func(t *testing.T) {
		np := engine.PrototypeByName("NestedListPairs")
		nrp := engine.PrototypeByName("NestedListPairs.Repr")
		var n schema.TypedNode
		t.Run("typed-create", func(t *testing.T) {
			n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("str").AssignString("boop")
				ma.AssembleEntry("lp").CreateMap(1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("field").AssignString("valoo")
				})
			}).(schema.TypedNode)
			t.Run("typed-read", func(t *testing.T) {
				qt.Assert(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, n.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(n.LookupByString("str"))), qt.Equals, "boop")
				lp := must.Node(n.LookupByString("lp"))
				qt.Check(t, lp.Kind(), qt.Equals, datamodel.Kind_Map)
				qt.Check(t, must.String(must.Node(lp.LookupByString("field"))), qt.Equals, "valoo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, nr.Length(), qt.Equals, int64(2))
				kv := must.Node(nr.LookupByIndex(0))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "str")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "boop")
				kv = must.Node(nr.LookupByIndex(1))
				qt.Assert(t, nr.Kind(), qt.Equals, datamodel.Kind_List)
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "lp")
				lp := must.Node(kv.LookupByIndex(1))
				qt.Check(t, lp.Kind(), qt.Equals, datamodel.Kind_List)
				kv = must.Node(lp.LookupByIndex(0))
				qt.Check(t, kv.Length(), qt.Equals, int64(2))
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(0))), qt.Equals, "field")
				qt.Check(t, must.String(must.Node(kv.LookupByIndex(1))), qt.Equals, "valoo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildList(nrp, 1, func(la fluent.ListAssembler) {
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("str")
					la.AssembleValue().AssignString("boop")
				})
				la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
					la.AssembleValue().AssignString("lp")
					la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
						la.AssembleValue().CreateList(2, func(la fluent.ListAssembler) {
							la.AssembleValue().AssignString("field")
							la.AssembleValue().AssignString("valoo")
						})
					})
				})
			})
			qt.Check(t, n, NodeContentEquals, nr)
		})
	})
}
