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
	qt.Check(t, Sprint(n), qt.CmpEquals(), wish.Dedent(`
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
		qt.Check(t, Sprint(n), qt.CmpEquals(), wish.Dedent(`
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
		qt.Check(t, Sprint(n), qt.CmpEquals(), wish.Dedent(`
			map<WowMap>{
				struct<FooBar>{foo: string<String>{"x"}, bar: string<String>{"y"}}: string<String>{"a"}
				struct<FooBar>{foo: string<String>{"z"}, bar: string<String>{"z"}}: string<String>{"b"}
			}`,
		))
	})
	t.Run("map-with-nested-struct-keys", func(t *testing.T) {
		type Baz struct {
			Baz string
		}
		type FooBar struct {
			Foo string
			Bar Baz
			Baz Baz
		}
		type WowMap struct {
			Keys   []FooBar
			Values map[FooBar]Baz
		}
		ts := schema.MustTypeSystem(
			schema.SpawnString("String"),
			schema.SpawnStruct("FooBar", []schema.StructField{
				schema.SpawnStructField("foo", "String", false, false),
				schema.SpawnStructField("bar", "Baz", false, false),
				schema.SpawnStructField("baz", "Baz", false, false),
			}, schema.SpawnStructRepresentationStringjoin(":")),
			schema.SpawnStruct("Baz", []schema.StructField{
				schema.SpawnStructField("baz", "String", false, false),
			}, schema.SpawnStructRepresentationStringjoin(":")),
			schema.SpawnMap("WowMap", "FooBar", "Baz", false),
		)
		n := bindnode.Wrap(&WowMap{
			Keys: []FooBar{{"x", Baz{"y"}, Baz{"y"}}, {"z", Baz{"z"}, Baz{"z"}}},
			Values: map[FooBar]Baz{
				{"x", Baz{"y"}, Baz{"y"}}: Baz{"a"},
				{"z", Baz{"z"}, Baz{"z"}}: Baz{"b"},
			},
		}, ts.TypeByName("WowMap"))
		t.Run("complex-keys-in-effect", func(t *testing.T) {
			cfg := Config{
				UseMapComplexStyleAlways: true,
			}
			qt.Check(t, cfg.Sprint(n), qt.CmpEquals(), wish.Dedent(`
				map<WowMap>{
					struct<FooBar>{
							foo: string<String>{"x"}
							bar: struct<Baz>{
								baz: string<String>{"y"}
							}
							baz: struct<Baz>{
								baz: string<String>{"y"}
							}
					}: struct<Baz>{
						baz: string<String>{"a"}
					}
					struct<FooBar>{
							foo: string<String>{"z"}
							bar: struct<Baz>{
								baz: string<String>{"z"}
							}
							baz: struct<Baz>{
								baz: string<String>{"z"}
							}
					}: struct<Baz>{
						baz: string<String>{"b"}
					}
				}`,
			))
		})
		t.Run("complex-keys-in-disabled", func(t *testing.T) {
			cfg := Config{
				UseMapComplexStyleOnType: map[schema.TypeName]bool{
					"WowMap": false,
				},
			}
			qt.Check(t, cfg.Sprint(n), qt.CmpEquals(), wish.Dedent(`
				map<WowMap>{
					struct<FooBar>{foo: string<String>{"x"}, bar: struct<Baz>{baz: string<String>{"y"}}, baz: struct<Baz>{baz: string<String>{"y"}}}: struct<Baz>{
						baz: string<String>{"a"}
					}
					struct<FooBar>{foo: string<String>{"z"}, bar: struct<Baz>{baz: string<String>{"z"}}, baz: struct<Baz>{baz: string<String>{"z"}}}: struct<Baz>{
						baz: string<String>{"b"}
					}
				}`,
			))
		})
	})

}
