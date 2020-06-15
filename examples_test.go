package ipld_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func ExampleCreateDataAndMarshal() {

	ns := basicnode.Style.Any // Pick a style for the in-memory data.
	nb := ns.NewBuilder()     // Create a builder.
	ma, _ := nb.BeginMap(2)   // Begin assembling a map.
	ma.AssembleKey().AssignString("hey")
	ma.AssembleValue().AssignString("it works!")
	ma.AssembleKey().AssignString("yes")
	ma.AssembleValue().AssignBool(true)
	ma.Finish()     // Call 'Finish' on the map assembly to let it know no more data is coming.
	n := nb.Build() // Call 'Build' to get the resulting Node.  (It's immutable!)

	dagjson.Encoder(n, os.Stdout)

	// Output:
	// {
	//	"hey": "it works!",
	//	"yes": true
	// }
}

func ExampleUnmarshalData() {
	serial := bytes.NewBufferString(`{"hey":"it works!","yes": true}`)

	ns := basicnode.Style.Any   // Pick a stle for the in-memory data.
	nb := ns.NewBuilder()       // Create a builder.
	dagjson.Decoder(nb, serial) // Hand the builder to decoding -- decoding will fill it in!
	n := nb.Build()             // Call 'Build' to get the resulting Node.  (It's immutable!)

	fmt.Printf("the data decoded was a %s kind\n", n.ReprKind())
	fmt.Printf("the length of the node is %d\n", n.Length())

	// Output:
	// the data decoded was a map kind
	// the length of the node is 2
}
