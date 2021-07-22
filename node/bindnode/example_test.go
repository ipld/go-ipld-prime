package bindnode_test

import (
	"bytes"
	"fmt"
	"os"

	ipld "github.com/ipld/go-ipld-prime"
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

func ExampleWrap_withSchemaUsingStringjoinStruct() {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("FooBarBaz",
		[]schema.StructField{
			schema.SpawnStructField("foo", "String", false, false),
			schema.SpawnStructField("bar", "String", false, false),
			schema.SpawnStructField("baz", "String", false, false),
		},
		schema.SpawnStructRepresentationStringjoin(":"),
	))
	schemaType := ts.TypeByName("FooBarBaz")

	type FooBarBaz struct {
		Foo string
		Bar string
		Baz string
	}
	fbb := &FooBarBaz{
		Foo: "x",
		Bar: "y",
		Baz: "z",
	}
	node := bindnode.Wrap(fbb, schemaType)

	// Take the representation of the node, and serialize it.
	nodeRepr := node.Representation()
	var buf bytes.Buffer
	dagjson.Encode(nodeRepr, &buf)

	// Output how this was serialized, for the example.
	fmt.Fprintf(os.Stdout, "json: %s\n", buf.Bytes())

	// Now unmarshal that again and print that too, to show it working both ways.
	np := bindnode.Prototype(&FooBarBaz{}, schemaType)
	nb := np.Representation().NewBuilder()
	err := dagjson.Decode(nb, &buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("golang: %#v\n", bindnode.Unwrap(nb.Build()))

	// Output:
	// json: "x:y:z"
	// golang: &bindnode_test.FooBarBaz{Foo:"x", Bar:"y", Baz:"z"}
}

func ExampleWrap_withSchemaUsingStringprefixUnion() {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnString("Foo"))
	ts.Accumulate(schema.SpawnString("Bar"))
	ts.Accumulate(schema.SpawnUnion("FooOrBar",
		[]schema.TypeName{
			"Foo",
			"Bar",
		},
		schema.SpawnUnionRepresentationStringprefix(
			":", // n.b. this API will change soon; a schema-schema iteration has removed the distinct joiner string.
			map[string]schema.TypeName{
				"foo": "Foo",
				"bar": "Bar",
			},
		),
	))
	schemaType := ts.TypeByName("FooOrBar")

	// The golang structures for unions don't look as simple as one might like.  Golang has no native sum types, so we do something interesting here.
	type FooOrBar struct {
		Index int
		Value interface{}
	}
	fob := &FooOrBar{
		Index: 1,
		Value: "oi",
	}
	node := bindnode.Wrap(fob, schemaType)

	// Take the representation of the node, and serialize it.
	nodeRepr := node.Representation()
	var buf bytes.Buffer
	dagjson.Encode(nodeRepr, &buf)

	// Output how this was serialized, for the example.
	fmt.Fprintf(os.Stdout, "json: %s\n", buf.Bytes())

	// Now unmarshal that again and print that too, to show it working both ways.
	np := bindnode.Prototype((*FooOrBar)(nil), schemaType)
	nb := np.Representation().NewBuilder()
	err := dagjson.Decode(nb, &buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("golang: %#v\n", bindnode.Unwrap(nb.Build()))

	// Output:
	// json: "bar:oi"
	// golang: &bindnode_test.FooOrBar{Index:1, Value:"oi"}
}
