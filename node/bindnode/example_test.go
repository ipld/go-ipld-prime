package bindnode_test

import (
	"os"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/polydawn/refmt/json"
)

func ExamplePrototypeOnlySchema() {
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
	proto := bindnode.PrototypeOnlySchema(schemaType)

	n, err := qp.BuildMap(proto, -1, func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, "Name", qp.String("Michael"))
		qp.MapEntry(ma, "Friends", qp.List(-1, func(la ipld.ListAssembler) {
			qp.ListEntry(la, qp.String("Sarah"))
			qp.ListEntry(la, qp.String("Alex"))
		}))
	})
	if err != nil {
		panic(err)
	}
	nr := n.(schema.TypedNode).Representation()
	dagjson.Marshal(nr, json.NewEncoder(os.Stdout, json.EncodeOptions{}), true)

	// Output:
	// {"Name":"Michael","Friends":["Sarah","Alex"]}
}
