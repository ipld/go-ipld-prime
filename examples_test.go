package ipld_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

// Example_createDataAndMarshal shows how you can feed data into a NodeBuilder,
// and also how to then hand that to an Encoder.
//
// Often you'll encoding implicitly through a LinkSystem.Store call instead,
// but you can do it directly, too.
func Example_createDataAndMarshal() {
	np := basicnode.Prototype.Any // Pick a prototype: this is how we decide what implementation will store the in-memory data.
	nb := np.NewBuilder()         // Create a builder.
	ma, _ := nb.BeginMap(2)       // Begin assembling a map.
	ma.AssembleKey().AssignString("hey")
	ma.AssembleValue().AssignString("it works!")
	ma.AssembleKey().AssignString("yes")
	ma.AssembleValue().AssignBool(true)
	ma.Finish()     // Call 'Finish' on the map assembly to let it know no more data is coming.
	n := nb.Build() // Call 'Build' to get the resulting Node.  (It's immutable!)

	dagjson.Encode(n, os.Stdout)

	// Output:
	// {"hey":"it works!","yes":true}
}

// Example_unmarshalData shows how you can use a Decoder
// and a NodeBuilder (or NodePrototype) together to do unmarshalling.
//
// Often you'll do this implicitly through a LinkSystem.Load call instead,
// but you can do it directly, too.
func Example_unmarshalData() {
	serial := strings.NewReader(`{"hey":"it works!","yes": true}`)

	np := basicnode.Prototype.Any // Pick a stle for the in-memory data.
	nb := np.NewBuilder()         // Create a builder.
	dagjson.Decode(nb, serial)    // Hand the builder to decoding -- decoding will fill it in!
	n := nb.Build()               // Call 'Build' to get the resulting Node.  (It's immutable!)

	fmt.Printf("the data decoded was a %s kind\n", n.Kind())
	fmt.Printf("the length of the node is %d\n", n.Length())

	// Output:
	// the data decoded was a map kind
	// the length of the node is 2
}

func ExampleLoadSchema() {
	ts, err := ipld.LoadSchema("sample.ipldsch", strings.NewReader(`
		type Root struct {
			foo Int
			bar nullable String
		}
		`))
	if err != nil {
		panic(err)
	}
	typeRoot := ts.TypeByName("Root").(*schema.TypeStruct)
	for _, field := range typeRoot.Fields() {
		fmt.Printf("field name=%q nullable=%t type=%v\n",
			field.Name(), field.IsNullable(), field.Type().Name())
	}
	// Output:
	// field name="foo" nullable=false type=Int
	// field name="bar" nullable=true type=String
}

// Example_goValueWithSchema shows how to combine a Go value with an IPLD
// schema, which can then be used as an IPLD node.
//
// For more examples and documentation, see the node/bindnode package.
func Example_goValueWithSchema() {
	type Person struct {
		Name    string
		Age     int
		Friends []string
	}

	ts, err := ipld.LoadSchemaBytes([]byte(`
		type Person struct {
			name    String
			age     Int
			friends [String]
		} representation tuple
	`))
	if err != nil {
		panic(err)
	}
	schemaType := ts.TypeByName("Person")
	person := &Person{Name: "Alice", Age: 34, Friends: []string{"Bob"}}
	node := bindnode.Wrap(person, schemaType)

	fmt.Printf("%#v\n", person)
	dagjson.Encode(node.Representation(), os.Stdout)

	// Output:
	// &ipld_test.Person{Name:"Alice", Age:34, Friends:[]string{"Bob"}}
	// ["Alice",34,["Bob"]]
}
