//go:build ignore
// +build ignore

package main

import (
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
	gengo.Generate(".", pkgName, ts, adjCfg)
}
