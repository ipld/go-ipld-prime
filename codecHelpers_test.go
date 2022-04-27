package ipld_test

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/schema"
)

func Example_marshal() {
	type Foobar struct {
		Foo string
		Bar string
	}
	encoded, err := ipld.Marshal(json.Encode, &Foobar{"wow", "whee"}, nil)
	fmt.Printf("error: %v\n", err)
	fmt.Printf("data: %s\n", string(encoded))

	// Output:
	// error: <nil>
	// data: {
	// 	"Foo": "wow",
	// 	"Bar": "whee"
	// }
}

// TODO: Example_Unmarshal, which uses nil and infers a typesystem.  However, to match Example_Unmarshal_withSchema, that appears to need more features in bindnode.

func Example_unmarshal_withSchema() {
	typesys := schema.MustTypeSystem(
		schema.SpawnStruct("Foobar",
			[]schema.StructField{
				schema.SpawnStructField("foo", "String", false, false),
				schema.SpawnStructField("bar", "String", false, false),
			},
			schema.SpawnStructRepresentationMap(nil),
		),
		schema.SpawnString("String"),
	)

	type Foobar struct {
		Foo string
		Bar string
	}
	serial := []byte(`{"foo":"wow","bar":"whee"}`)
	foobar := Foobar{}
	n, err := ipld.Unmarshal(serial, json.Decode, &foobar, typesys.TypeByName("Foobar"))
	fmt.Printf("error: %v\n", err)
	fmt.Printf("go struct: %v\n", foobar)
	fmt.Printf("node kind and length: %s, %d\n", n.Kind(), n.Length())
	fmt.Printf("node lookup 'foo': %q\n", must.String(must.Node(n.LookupByString("foo"))))

	// Output:
	// error: <nil>
	// go struct: {wow whee}
	// node kind and length: map, 2
	// node lookup 'foo': "wow"
}
