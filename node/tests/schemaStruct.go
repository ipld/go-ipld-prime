package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestRequiredFields(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("StructOne",
		[]schema.StructField{
			schema.SpawnStructField("a", "String", false, false),
			schema.SpawnStructField("b", "String", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			// no renames.  we expect a simpler error message in this case.
		}),
	))
	ts.Accumulate(schema.SpawnStruct("StructTwo",
		[]schema.StructField{
			schema.SpawnStructField("a", "String", false, false),
			schema.SpawnStructField("b", "String", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"b": "z",
		}),
	))
	engine.Init(t, ts)

	t.Run("building-type-without-required-fields-errors", func(t *testing.T) {
		np := engine.PrototypeByName("StructOne")

		nb := np.NewBuilder()
		ma, _ := nb.BeginMap(0)
		err := ma.Finish()

		Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
		Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b`)
	})
	t.Run("building-representation-without-required-fields-errors", func(t *testing.T) {
		nrp := engine.PrototypeByName("StructOne.Repr")

		nb := nrp.NewBuilder()
		ma, _ := nb.BeginMap(0)
		err := ma.Finish()

		Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
		Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b`)
	})
	t.Run("building-representation-with-renames-without-required-fields-errors", func(t *testing.T) {
		nrp := engine.PrototypeByName("StructTwo.Repr")

		nb := nrp.NewBuilder()
		ma, _ := nb.BeginMap(0)
		err := ma.Finish()

		Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
		Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b (serial:"z")`)
	})
}

func SchemaTestStructNesting(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("SmolStruct",
		[]schema.StructField{
			schema.SpawnStructField("s", "String", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"s": "q",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("GulpoStruct",
		[]schema.StructField{
			schema.SpawnStructField("x", "SmolStruct", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"x": "r",
		}),
	))
	engine.Init(t, ts)

	np := engine.PrototypeByName("GulpoStruct")
	nrp := engine.PrototypeByName("GulpoStruct.Repr")
	var n schema.TypedNode
	t.Run("typed-create", func(t *testing.T) {
		n = fluent.MustBuildMap(np, 1, func(ma fluent.MapAssembler) {
			ma.AssembleEntry("x").CreateMap(1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("s").AssignString("woo")
			})
		}).(schema.TypedNode)
		t.Run("typed-read", func(t *testing.T) {
			Require(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
			Wish(t, n.Length(), ShouldEqual, int64(1))

			n2 := must.Node(n.LookupByString("x"))
			Require(t, n2.Kind(), ShouldEqual, ipld.Kind_Map)

			n2Seg := must.Node(n.LookupBySegment(ipld.PathSegmentOfString("x")))
			Wish(t, ipld.DeepEqual(n2, n2Seg), ShouldEqual, true)

			n2Node := must.Node(n.LookupByNode(basicnode.NewString("x")))
			Wish(t, ipld.DeepEqual(n2, n2Node), ShouldEqual, true)

			Wish(t, must.String(must.Node(n2.LookupByString("s"))), ShouldEqual, "woo")
		})
		t.Run("repr-read", func(t *testing.T) {
			nr := n.Representation()
			Require(t, nr.Kind(), ShouldEqual, ipld.Kind_Map)
			Wish(t, nr.Length(), ShouldEqual, int64(1))

			n2 := must.Node(nr.LookupByString("r"))
			Require(t, n2.Kind(), ShouldEqual, ipld.Kind_Map)

			n2Seg := must.Node(nr.LookupBySegment(ipld.PathSegmentOfString("r")))
			Wish(t, ipld.DeepEqual(n2, n2Seg), ShouldEqual, true)

			n2Node := must.Node(nr.LookupByNode(basicnode.NewString("r")))
			Wish(t, ipld.DeepEqual(n2, n2Node), ShouldEqual, true)

			Wish(t, must.String(must.Node(n2.LookupByString("q"))), ShouldEqual, "woo")
		})
	})
	t.Run("repr-create", func(t *testing.T) {
		nr := fluent.MustBuildMap(nrp, 1, func(ma fluent.MapAssembler) {
			ma.AssembleEntry("r").CreateMap(1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("q").AssignString("woo")
			})
		})
		Wish(t, ipld.DeepEqual(n, nr), ShouldEqual, true)
	})
}
