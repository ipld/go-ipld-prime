package ipld_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
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
