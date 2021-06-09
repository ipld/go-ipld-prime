package bindnode_test

import (
	"os"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	refmtjson "github.com/polydawn/refmt/json"
)

func ExampleWrap_withSchema() {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnStruct("Person",
		[]schema.StructField{
			schema.SpawnStructField("Name", "String", false, false),
			schema.SpawnStructField("Age", "Int", true, false),
			schema.SpawnStructField("Friends", "List_String", false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	ts.Accumulate(schema.SpawnList("List_String", "String", false))

	schemaType := ts.TypeByName("Person")

	type Person struct {
		Name    string
		Age     *int64 // optional
		Friends []string
	}
	person := &Person{
		Name:    "Michael",
		Friends: []string{"Sarah", "Alex"},
	}
	node := bindnode.Wrap(person, schemaType)

	nodeRepr := node.Representation()
	dagjson.Marshal(nodeRepr, refmtjson.NewEncoder(os.Stdout, refmtjson.EncodeOptions{}), true)

	// Output:
	// {"Name":"Michael","Friends":["Sarah","Alex"]}
}

func ExamplePrototype_onlySchema() {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnStruct("Person",
		[]schema.StructField{
			schema.SpawnStructField("Name", "String", false, false),
			schema.SpawnStructField("Age", "Int", true, false),
			schema.SpawnStructField("Friends", "List_String", false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	ts.Accumulate(schema.SpawnList("List_String", "String", false))

	schemaType := ts.TypeByName("Person")
	proto := bindnode.Prototype(nil, schemaType)

	node, err := qp.BuildMap(proto, -1, func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, "Name", qp.String("Michael"))
		qp.MapEntry(ma, "Friends", qp.List(-1, func(la ipld.ListAssembler) {
			qp.ListEntry(la, qp.String("Sarah"))
			qp.ListEntry(la, qp.String("Alex"))
		}))
	})
	if err != nil {
		panic(err)
	}

	nodeRepr := node.(schema.TypedNode).Representation()
	dagjson.Marshal(nodeRepr, refmtjson.NewEncoder(os.Stdout, refmtjson.EncodeOptions{}), true)

	// Output:
	// {"Name":"Michael","Friends":["Sarah","Alex"]}
}
