package bindnode_test

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
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
	dagjson.Encode(nodeRepr, os.Stdout)

	// Output:
	// {"Friends":["Sarah","Alex"],"Name":"Michael"}
}

func ExampleWrap_noSchema() {
	type Person struct {
		Name    string
		Age     int64 // TODO: optional to match other examples
		Friends []string
	}
	person := &Person{
		Name:    "Michael",
		Friends: []string{"Sarah", "Alex"},
	}
	node := bindnode.Wrap(person, nil)

	nodeRepr := node.Representation()
	dagjson.Encode(nodeRepr, os.Stdout)

	// Output:
	// {"Age":0,"Friends":["Sarah","Alex"],"Name":"Michael"}
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
	dagjson.Encode(nodeRepr, os.Stdout)

	// Output:
	// {"Friends":["Sarah","Alex"],"Name":"Michael"}
}

func ExamplePrototype_union() {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnUnion("StringOrInt",
		[]schema.TypeName{
			"String",
			"Int",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"hasString": "String",
			"hasInt":    "Int",
		}),
	))

	schemaType := ts.TypeByName("StringOrInt")

	type CustomIntType int64
	type StringOrInt struct {
		String *string
		Int    *CustomIntType // We can use custom types, too.
	}

	proto := bindnode.Prototype((*StringOrInt)(nil), schemaType)

	node, err := qp.BuildMap(proto.Representation(), -1, func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, "hasInt", qp.Int(123))
	})
	if err != nil {
		panic(err)
	}

	fmt.Print("Type level DAG-JSON: ")
	dagjson.Encode(node, os.Stdout)
	fmt.Println()

	fmt.Print("Representation level DAG-JSON: ")
	nodeRepr := node.(schema.TypedNode).Representation()
	dagjson.Encode(nodeRepr, os.Stdout)
	fmt.Println()

	// Inspect what the underlying Go value contains.
	union := bindnode.Unwrap(node).(*StringOrInt)
	switch {
	case union.String != nil:
		fmt.Printf("Go StringOrInt.String: %v\n", *union.String)
	case union.Int != nil:
		fmt.Printf("Go StringOrInt.Int: %v\n", *union.Int)
	}

	// Output:
	// Type level DAG-JSON: {"Int":123}
	// Representation level DAG-JSON: {"hasInt":123}
	// Go StringOrInt.Int: 123
}
