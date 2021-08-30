package printer

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestSimpleData(t *testing.T) {
	n, _ := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "some key", qp.String("some value"))
		qp.MapEntry(ma, "another key", qp.String("another value"))
		qp.MapEntry(ma, "nested map", qp.Map(2, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, "deeper entries", qp.String("deeper values"))
			qp.MapEntry(ma, "more deeper entries", qp.String("more deeper values"))
		}))
		qp.MapEntry(ma, "nested list", qp.List(2, func(la datamodel.ListAssembler) {
			qp.ListEntry(la, qp.Int(1))
			qp.ListEntry(la, qp.Int(2))
		}))
	})
	qt.Check(t, Sprint(n), qt.Equals, wish.Dedent(`
		map{
			string{"some key"}: string{"some value"}
			string{"another key"}: string{"another value"}
			string{"nested map"}: map{
				string{"deeper entries"}: string{"deeper values"}
				string{"more deeper entries"}: string{"more deeper values"}
			}
			string{"nested list"}: list{
				0: int{1}
				1: int{2}
			}
		}`,
	))
}

func TestTypedData(t *testing.T) {
	t.Run("structs", func(t *testing.T) {
		type FooBar struct {
			Foo string
			Bar string
		}
		ts := schema.MustTypeSystem(
			schema.SpawnString("String"),
			schema.SpawnStruct("FooBar", []schema.StructField{
				schema.SpawnStructField("foo", "String", false, false),
				schema.SpawnStructField("bar", "String", false, false),
			}, nil),
		)
		n := bindnode.Wrap(&FooBar{"x", "y"}, ts.TypeByName("FooBar"))
		qt.Check(t, Sprint(n), qt.Equals, wish.Dedent(`
			struct<FooBar>{
				foo: string<String>{"x"}
				bar: string<String>{"y"}
			}`,
		))
	})
	t.Run("map-with-struct-keys", func(t *testing.T) {
		type FooBar struct {
			Foo string
			Bar string
		}
		type WowMap struct {
			Keys   []FooBar
			Values map[FooBar]string
		}
		ts := schema.MustTypeSystem(
			schema.SpawnString("String"),
			schema.SpawnStruct("FooBar", []schema.StructField{
				schema.SpawnStructField("foo", "String", false, false),
				schema.SpawnStructField("bar", "String", false, false),
			}, schema.SpawnStructRepresentationStringjoin(":")),
			schema.SpawnMap("WowMap", "FooBar", "String", false),
		)
		n := bindnode.Wrap(&WowMap{
			Keys: []FooBar{{"x", "y"}, {"z", "z"}},
			Values: map[FooBar]string{
				{"x", "y"}: "a",
				{"z", "z"}: "b",
			},
		}, ts.TypeByName("WowMap"))
		qt.Check(t, Sprint(n), qt.Equals, wish.Dedent(`
			map<WowMap>{
				struct<FooBar>{foo: string<String>{"x"}, bar: string<String>{"y"}}: string<String>{"a"}
				struct<FooBar>{foo: string<String>{"z"}, bar: string<String>{"z"}}: string<String>{"b"}
			}`,
		))
	})

}
