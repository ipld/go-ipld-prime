package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestRequiredFields(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{}
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

	prefix := "struct-required-fields"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		t.Run("building-type-without-required-fields-errors", func(t *testing.T) {
			np := getPrototypeByName("StructOne")

			nb := np.NewBuilder()
			ma, _ := nb.BeginMap(0)
			err := ma.Finish()

			Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
			Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b`)
		})
		t.Run("building-representation-without-required-fields-errors", func(t *testing.T) {
			nrp := getPrototypeByName("StructOne.Repr")

			nb := nrp.NewBuilder()
			ma, _ := nb.BeginMap(0)
			err := ma.Finish()

			Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
			Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b`)
		})
		t.Run("building-representation-with-renames-without-required-fields-errors", func(t *testing.T) {
			nrp := getPrototypeByName("StructTwo.Repr")

			nb := nrp.NewBuilder()
			ma, _ := nb.BeginMap(0)
			err := ma.Finish()

			Wish(t, err, ShouldBeSameTypeAs, ipld.ErrMissingRequiredField{})
			Wish(t, err.Error(), ShouldEqual, `missing required fields: a,b (serial:"z")`)
		})
	})
}

func TestStructNesting(t *testing.T) {
	t.Parallel()

	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
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

	prefix := "struct-nesting"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype) {
		np := getPrototypeByName("GulpoStruct")
		nrp := getPrototypeByName("GulpoStruct.Repr")
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
				Wish(t, must.String(must.Node(n2.LookupByString("s"))), ShouldEqual, "woo")
			})
			t.Run("repr-read", func(t *testing.T) {
				nr := n.Representation()
				Require(t, nr.Kind(), ShouldEqual, ipld.Kind_Map)
				Wish(t, nr.Length(), ShouldEqual, int64(1))
				n2 := must.Node(nr.LookupByString("r"))
				Require(t, n2.Kind(), ShouldEqual, ipld.Kind_Map)
				Wish(t, must.String(must.Node(n2.LookupByString("q"))), ShouldEqual, "woo")
			})
		})
		t.Run("repr-create", func(t *testing.T) {
			nr := fluent.MustBuildMap(nrp, 1, func(ma fluent.MapAssembler) {
				ma.AssembleEntry("r").CreateMap(1, func(ma fluent.MapAssembler) {
					ma.AssembleEntry("q").AssignString("woo")
				})
			})
			Wish(t, n, ShouldEqual, nr)
		})
	})
}
