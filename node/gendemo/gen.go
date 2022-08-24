//go:build ignore

package main

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

func main() {
	pkgName := "gendemo"
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &gengo.AdjunctCfg{}
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Msg3",
		[]schema.StructField{
			schema.SpawnStructField("whee", "Int", false, false),
			schema.SpawnStructField("woot", "Int", false, false),
			schema.SpawnStructField("waga", "Int", false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	ts.Accumulate(schema.SpawnMap("Map__String__Msg3",
		"String", "Msg3", false))
	ts.Accumulate(schema.SpawnBool("Bar"))
	ts.Accumulate(schema.SpawnString("Baz"))
	ts.Accumulate(schema.SpawnInt("Foo"))
	ts.Accumulate(schema.SpawnUnion("UnionKinded",
		[]schema.TypeName{
			"Foo",
			"Bar",
			"Baz",
		}, schema.SpawnUnionRepresentationKinded(
			map[datamodel.Kind]schema.TypeName{
				datamodel.Kind_Int:    "Foo",
				datamodel.Kind_Bool:   "Bar",
				datamodel.Kind_String: "Baz",
			}),
	))
	gengo.Generate(".", pkgName, ts, adjCfg)
}
