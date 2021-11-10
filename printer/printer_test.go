package printer

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/warpfork/go-wish"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

var testLink = func() datamodel.Link {
	someCid, _ := cid.Cast([]byte{1, 85, 0, 5, 0, 1, 2, 3, 4})
	return cidlink.Link{Cid: someCid}
}()

func TestSimpleData(t *testing.T) {
	t.Run("nested-maps", func(t *testing.T) {
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
	})

	t.Run("map-with-link-and-bytes", func(t *testing.T) {
		n, _ := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, "some key", qp.Link(testLink))
			qp.MapEntry(ma, "another key", qp.String("another value"))
			qp.MapEntry(ma, "nested map", qp.Map(2, func(ma datamodel.MapAssembler) {
				qp.MapEntry(ma, "deeper entries", qp.String("deeper values"))
				qp.MapEntry(ma, "more deeper entries", qp.Link(testLink))
				qp.MapEntry(ma, "yet another deeper entries", qp.Bytes([]byte("fish")))
			}))
			qp.MapEntry(ma, "nested list", qp.List(2, func(la datamodel.ListAssembler) {
				qp.ListEntry(la, qp.Bytes([]byte("ghoti")))
				qp.ListEntry(la, qp.Int(1))
				qp.ListEntry(la, qp.Link(testLink))
			}))
		})
		qt.Check(t, Sprint(n), qt.CmpEquals(), wish.Dedent(`
		map{
			string{"some key"}: link{bafkqabiaaebagba}
			string{"another key"}: string{"another value"}
			string{"nested map"}: map{
				string{"deeper entries"}: string{"deeper values"}
				string{"more deeper entries"}: link{bafkqabiaaebagba}
				string{"yet another deeper entries"}: bytes{66697368}
			}
			string{"nested list"}: list{
				0: bytes{67686f7469}
				1: int{1}
				2: link{bafkqabiaaebagba}
			}
		}`,
		))
	})
}

func TestTypedData(t *testing.T) {
	t.Run("structs", func(t *testing.T) {
		type FooBar struct {
			Foo  string
			Bar  string
			Baz  []byte
			Jazz datamodel.Link
		}
		ts := schema.MustTypeSystem(
			schema.SpawnString("String"),
			schema.SpawnBytes("Bytes"),
			schema.SpawnLink("Link"),
			schema.SpawnStruct("FooBar", []schema.StructField{
				schema.SpawnStructField("foo", "String", false, false),
				schema.SpawnStructField("bar", "String", false, false),
				schema.SpawnStructField("baz", "Bytes", false, false),
				schema.SpawnStructField("jazz", "Link", false, false),
			}, nil),
		)
		n := bindnode.Wrap(&FooBar{"x", "y", []byte("zed"), testLink}, ts.TypeByName("FooBar"))
		qt.Check(t, Sprint(n), qt.CmpEquals(), wish.Dedent(`
			struct<FooBar>{
				foo: string<String>{"x"}
				bar: string<String>{"y"}
				baz: bytes<Bytes>{7a6564}
				jazz: link<Link>{bafkqabiaaebagba}
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
				{"x", Baz{"y"}, Baz{"y"}}: {"a"},
				{"z", Baz{"z"}, Baz{"z"}}: {"b"},
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
